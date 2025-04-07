# API Gateway General Troubleshooting Guide

This guide provides general troubleshooting steps for common issues with the API Gateway.

## Diagnostic Tools

Before diving into specific issues, here are some useful diagnostic tools:

### 1. Check Pod Status

```bash
# List all API Gateway pods
kubectl get pods -n sparkfund -l app=api-gateway

# Get detailed information about a specific pod
kubectl describe pod <pod-name> -n sparkfund
```

### 2. Check Logs

```bash
# Get logs from all API Gateway pods
kubectl logs -n sparkfund -l app=api-gateway

# Get logs from a specific pod
kubectl logs <pod-name> -n sparkfund

# Get logs with timestamps
kubectl logs <pod-name> -n sparkfund --timestamps

# Get logs from the previous container instance
kubectl logs <pod-name> -n sparkfund --previous
```

### 3. Check Metrics

```bash
# Port forward Prometheus
kubectl port-forward svc/prometheus-server 9090:9090 -n monitoring

# Open http://localhost:9090 in your browser
# Useful queries:
# - sum(rate(http_requests_total{service="api-gateway"}[5m])) by (status_code)
# - histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{service="api-gateway"}[5m])) by (le))
# - sum(rate(http_requests_total{service="api-gateway", status_code=~"5.."}[5m]))
```

### 4. Check Configuration

```bash
# Get ConfigMap
kubectl get configmap api-gateway-config -n sparkfund -o yaml

# Get Secrets (encoded)
kubectl get secret api-gateway-tls -n sparkfund -o yaml
```

### 5. Check Network Policies

```bash
# Get Network Policies
kubectl get networkpolicy -n sparkfund

# Describe Network Policy
kubectl describe networkpolicy api-gateway -n sparkfund
```

## Common Issues and Solutions

### 1. API Gateway Not Responding

**Symptoms:**
- HTTP 502/503/504 errors
- Connection timeouts
- "Service Unavailable" errors

**Possible Causes:**
- Pods are not running
- Resource constraints (CPU/Memory)
- Network issues
- Backend services are down

**Troubleshooting Steps:**

1. Check if pods are running:
   ```bash
   kubectl get pods -n sparkfund -l app=api-gateway
   ```

2. Check for resource constraints:
   ```bash
   kubectl top pods -n sparkfund -l app=api-gateway
   ```

3. Check pod events:
   ```bash
   kubectl describe pod <pod-name> -n sparkfund
   ```

4. Check logs for errors:
   ```bash
   kubectl logs -n sparkfund -l app=api-gateway
   ```

5. Check if backend services are reachable:
   ```bash
   # Get a shell in the pod
   kubectl exec -it <pod-name> -n sparkfund -- sh
   
   # Test connectivity to backend services
   wget -O- http://user-service:8084/health
   wget -O- http://kyc-service:8081/health
   wget -O- http://investment-service:8082/health
   ```

**Solutions:**

1. Restart the pods:
   ```bash
   kubectl rollout restart deployment api-gateway -n sparkfund
   ```

2. Scale up the deployment:
   ```bash
   kubectl scale deployment api-gateway -n sparkfund --replicas=5
   ```

3. Check and adjust resource limits:
   ```bash
   kubectl edit deployment api-gateway -n sparkfund
   ```

4. Verify network policies allow traffic:
   ```bash
   kubectl describe networkpolicy api-gateway -n sparkfund
   ```

### 2. High Latency

**Symptoms:**
- Slow response times
- Timeouts
- Increased error rates during high load

**Possible Causes:**
- Insufficient resources
- Backend service latency
- Network issues
- Database bottlenecks
- Inefficient routing

**Troubleshooting Steps:**

1. Check response time metrics:
   ```
   histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{service="api-gateway"}[5m])) by (le))
   ```

2. Check CPU and memory usage:
   ```bash
   kubectl top pods -n sparkfund -l app=api-gateway
   ```

3. Check backend service latency:
   ```
   histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{service="api-gateway"}[5m])) by (le, backend_service))
   ```

4. Check network metrics:
   ```
   sum(rate(http_request_size_bytes_sum{service="api-gateway"}[5m])) / sum(rate(http_request_size_bytes_count{service="api-gateway"}[5m]))
   ```

**Solutions:**

1. Scale up the deployment:
   ```bash
   kubectl scale deployment api-gateway -n sparkfund --replicas=5
   ```

2. Adjust resource limits:
   ```bash
   kubectl edit deployment api-gateway -n sparkfund
   ```

3. Enable or adjust caching:
   ```bash
   kubectl edit configmap api-gateway-config -n sparkfund
   ```

4. Optimize backend service calls:
   - Implement circuit breakers
   - Add caching
   - Batch requests

### 3. Certificate Issues

**Symptoms:**
- TLS handshake failures
- Certificate validation errors
- "Your connection is not private" browser warnings

**Possible Causes:**
- Expired certificates
- Misconfigured certificates
- Missing intermediate certificates
- Hostname mismatch

**Troubleshooting Steps:**

1. Check certificate expiration:
   ```bash
   kubectl get secret api-gateway-tls -n sparkfund -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -noout -dates
   ```

2. Verify certificate chain:
   ```bash
   kubectl get secret api-gateway-tls -n sparkfund -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -noout -text
   ```

3. Check certificate configuration:
   ```bash
   kubectl describe certificate api-gateway-cert -n sparkfund
   ```

**Solutions:**

1. Renew the certificate:
   ```bash
   kubectl annotate certificate api-gateway-cert -n sparkfund cert-manager.io/renew="true"
   ```

2. Update the certificate configuration:
   ```bash
   kubectl edit certificate api-gateway-cert -n sparkfund
   ```

3. Manually replace the certificate:
   ```bash
   kubectl create secret tls api-gateway-tls -n sparkfund --cert=path/to/tls.crt --key=path/to/tls.key --dry-run=client -o yaml | kubectl apply -f -
   ```

### 4. Routing Issues

**Symptoms:**
- 404 errors for valid endpoints
- Requests sent to wrong backend services
- Inconsistent routing behavior

**Possible Causes:**
- Misconfigured routes
- Backend service discovery issues
- Path prefix issues
- Header-based routing issues

**Troubleshooting Steps:**

1. Check route configuration:
   ```bash
   kubectl get configmap api-gateway-config -n sparkfund -o yaml
   ```

2. Check service discovery:
   ```bash
   kubectl get endpoints -n sparkfund
   ```

3. Test routing manually:
   ```bash
   # Get a shell in the pod
   kubectl exec -it <pod-name> -n sparkfund -- sh
   
   # Test routing
   curl -v http://localhost:8080/api/v1/users
   ```

4. Check logs for routing decisions:
   ```bash
   kubectl logs <pod-name> -n sparkfund | grep "route"
   ```

**Solutions:**

1. Update route configuration:
   ```bash
   kubectl edit configmap api-gateway-config -n sparkfund
   ```

2. Restart the API Gateway:
   ```bash
   kubectl rollout restart deployment api-gateway -n sparkfund
   ```

3. Verify service endpoints:
   ```bash
   kubectl get endpoints user-service -n sparkfund
   kubectl get endpoints kyc-service -n sparkfund
   kubectl get endpoints investment-service -n sparkfund
   ```

## Escalation Procedure

If you cannot resolve the issue using this guide:

1. **Level 1**: Contact the on-call engineer via Slack (#api-gateway-support) or PagerDuty
2. **Level 2**: Escalate to the API Gateway team lead
3. **Level 3**: Escalate to the Platform Engineering team

## Useful Commands

```bash
# Get API Gateway version
kubectl exec -it <pod-name> -n sparkfund -- /api-gateway version

# Get API Gateway configuration
kubectl exec -it <pod-name> -n sparkfund -- cat /config/config.yaml

# Check connectivity to backend services
kubectl exec -it <pod-name> -n sparkfund -- wget -O- http://user-service:8084/health

# Check TLS configuration
kubectl exec -it <pod-name> -n sparkfund -- openssl s_client -connect api.sparkfund.com:443 -servername api.sparkfund.com

# Check DNS resolution
kubectl exec -it <pod-name> -n sparkfund -- nslookup user-service.sparkfund.svc.cluster.local
```
