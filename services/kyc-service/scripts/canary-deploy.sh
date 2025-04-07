#!/bin/bash

# Canary deployment script for KYC service
# This script manages canary deployments with gradual traffic shifting

set -e

# Default values
NAMESPACE="sparkfund"
SERVICE="kyc-service"
CANARY_WEIGHT=10
STABLE_WEIGHT=90
CANARY_IMAGE=""
STEP_PERCENTAGE=10
STEP_INTERVAL=60  # seconds
METRICS_CHECK=true
ERROR_THRESHOLD=1.0  # percentage

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --namespace)
      NAMESPACE="$2"
      shift
      shift
      ;;
    --service)
      SERVICE="$2"
      shift
      shift
      ;;
    --canary-image)
      CANARY_IMAGE="$2"
      shift
      shift
      ;;
    --initial-weight)
      CANARY_WEIGHT="$2"
      STABLE_WEIGHT=$((100 - CANARY_WEIGHT))
      shift
      shift
      ;;
    --step-percentage)
      STEP_PERCENTAGE="$2"
      shift
      shift
      ;;
    --step-interval)
      STEP_INTERVAL="$2"
      shift
      shift
      ;;
    --skip-metrics)
      METRICS_CHECK=false
      shift
      ;;
    --error-threshold)
      ERROR_THRESHOLD="$2"
      shift
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Validate required parameters
if [ -z "$CANARY_IMAGE" ]; then
  echo "Error: --canary-image is required"
  exit 1
fi

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Check for required commands
for cmd in kubectl jq curl; do
  if ! command_exists "$cmd"; then
    echo "Error: $cmd is required but not installed"
    exit 1
  fi
done

# Check if Istio is installed
if ! kubectl get crd virtualservices.networking.istio.io >/dev/null 2>&1; then
  echo "Error: Istio CRDs not found. Is Istio installed?"
  exit 1
fi

# Deploy canary version
echo "Deploying canary version with image: $CANARY_IMAGE"
kubectl -n $NAMESPACE set image deployment/$SERVICE-canary $SERVICE=$CANARY_IMAGE

# Wait for canary deployment to be ready
echo "Waiting for canary deployment to be ready..."
kubectl -n $NAMESPACE rollout status deployment/$SERVICE-canary --timeout=300s

# Create or update VirtualService for initial traffic split
echo "Setting initial traffic split: $STABLE_WEIGHT% stable, $CANARY_WEIGHT% canary"
cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: $SERVICE
  namespace: $NAMESPACE
spec:
  hosts:
  - $SERVICE
  - $SERVICE.${NAMESPACE}.svc.cluster.local
  gateways:
  - mesh
  http:
  - route:
    - destination:
        host: $SERVICE
        subset: stable
      weight: $STABLE_WEIGHT
    - destination:
        host: $SERVICE
        subset: canary
      weight: $CANARY_WEIGHT
EOF

# Function to check error rate
check_error_rate() {
  if [ "$METRICS_CHECK" = false ]; then
    return 0
  fi
  
  # Get error rate from Prometheus
  ERROR_RATE=$(curl -s "http://prometheus:9090/api/v1/query" \
    --data-urlencode "query=sum(rate(http_requests_total{app=\"$SERVICE\",version=\"canary\",status=~\"5..\"}[1m])) / sum(rate(http_requests_total{app=\"$SERVICE\",version=\"canary\"}[1m])) * 100" | \
    jq -r '.data.result[0].value[1]')
  
  if [ -z "$ERROR_RATE" ] || [ "$ERROR_RATE" = "null" ]; then
    echo "Warning: Could not get error rate from Prometheus"
    return 0
  fi
  
  echo "Current canary error rate: $ERROR_RATE%"
  
  # Check if error rate exceeds threshold
  if (( $(echo "$ERROR_RATE > $ERROR_THRESHOLD" | bc -l) )); then
    echo "Error rate exceeds threshold ($ERROR_THRESHOLD%). Rolling back."
    return 1
  fi
  
  return 0
}

# Gradually increase canary traffic
while [ $CANARY_WEIGHT -lt 100 ]; do
  # Check error rate before increasing traffic
  if ! check_error_rate; then
    echo "Rolling back canary deployment..."
    kubectl -n $NAMESPACE set image deployment/$SERVICE-canary $SERVICE=$(kubectl -n $NAMESPACE get deployment $SERVICE -o jsonpath='{.spec.template.spec.containers[0].image}')
    
    # Reset traffic to 100% stable
    cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: $SERVICE
  namespace: $NAMESPACE
spec:
  hosts:
  - $SERVICE
  - $SERVICE.${NAMESPACE}.svc.cluster.local
  gateways:
  - mesh
  http:
  - route:
    - destination:
        host: $SERVICE
        subset: stable
      weight: 100
    - destination:
        host: $SERVICE
        subset: canary
      weight: 0
EOF
    exit 1
  fi
  
  # Wait before increasing traffic
  echo "Waiting $STEP_INTERVAL seconds before increasing canary traffic..."
  sleep $STEP_INTERVAL
  
  # Increase canary traffic
  CANARY_WEIGHT=$((CANARY_WEIGHT + STEP_PERCENTAGE))
  if [ $CANARY_WEIGHT -gt 100 ]; then
    CANARY_WEIGHT=100
  fi
  STABLE_WEIGHT=$((100 - CANARY_WEIGHT))
  
  echo "Updating traffic split: $STABLE_WEIGHT% stable, $CANARY_WEIGHT% canary"
  cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: $SERVICE
  namespace: $NAMESPACE
spec:
  hosts:
  - $SERVICE
  - $SERVICE.${NAMESPACE}.svc.cluster.local
  gateways:
  - mesh
  http:
  - route:
    - destination:
        host: $SERVICE
        subset: stable
      weight: $STABLE_WEIGHT
    - destination:
        host: $SERVICE
        subset: canary
      weight: $CANARY_WEIGHT
EOF
done

# Promote canary to stable
echo "Canary deployment successful. Promoting canary to stable..."
kubectl -n $NAMESPACE set image deployment/$SERVICE $SERVICE=$CANARY_IMAGE

# Wait for stable deployment to be ready
echo "Waiting for stable deployment to be ready..."
kubectl -n $NAMESPACE rollout status deployment/$SERVICE --timeout=300s

# Reset traffic to 100% stable
cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: $SERVICE
  namespace: $NAMESPACE
spec:
  hosts:
  - $SERVICE
  - $SERVICE.${NAMESPACE}.svc.cluster.local
  gateways:
  - mesh
  http:
  - route:
    - destination:
        host: $SERVICE
        subset: stable
      weight: 100
    - destination:
        host: $SERVICE
        subset: canary
      weight: 0
EOF

echo "Canary deployment completed successfully!"
