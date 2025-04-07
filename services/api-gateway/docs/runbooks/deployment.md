# API Gateway Deployment Runbook

This runbook describes the process for deploying the API Gateway service.

## Prerequisites

- Access to the Kubernetes cluster
- `kubectl` configured to access the cluster
- Access to the container registry
- CI/CD pipeline access

## Deployment Process

### 1. Standard Deployment

The API Gateway is deployed automatically through the CI/CD pipeline when changes are merged to the main branch. The pipeline performs the following steps:

1. Build and test the code
2. Run security scans
3. Build and push the Docker image
4. Deploy to staging environment
5. Run integration tests
6. Deploy to production environment

### 2. Manual Deployment

In case a manual deployment is needed:

```bash
# 1. Build the Docker image
docker build -t ghcr.io/sparkfund/api-gateway:latest ./services/api-gateway

# 2. Push the Docker image
docker push ghcr.io/sparkfund/api-gateway:latest

# 3. Update the Kubernetes deployment
kubectl set image deployment/api-gateway api-gateway=ghcr.io/sparkfund/api-gateway:latest -n sparkfund

# 4. Monitor the rollout
kubectl rollout status deployment/api-gateway -n sparkfund
```

### 3. Canary Deployment

For a canary deployment:

```bash
# 1. Deploy the canary version
kubectl apply -f services/api-gateway/k8s/canary-deployment.yaml

# 2. Monitor the canary deployment
kubectl get canary api-gateway -n sparkfund -o yaml

# 3. Check metrics during the canary deployment
kubectl get canary api-gateway -n sparkfund -o jsonpath='{.status.canaryWeight}'
kubectl get canary api-gateway -n sparkfund -o jsonpath='{.status.failedChecks}'
kubectl get canary api-gateway -n sparkfund -o jsonpath='{.status.iterations}'
```

## Post-Deployment Verification

After deployment, verify that the API Gateway is functioning correctly:

1. Check that the pods are running:
   ```bash
   kubectl get pods -n sparkfund -l app=api-gateway
   ```

2. Check the logs for any errors:
   ```bash
   kubectl logs -n sparkfund -l app=api-gateway
   ```

3. Verify that the service is responding:
   ```bash
   curl -I https://api.sparkfund.com/health
   ```

4. Check the metrics:
   ```bash
   kubectl port-forward svc/prometheus-server 9090:9090 -n monitoring
   # Open http://localhost:9090 in your browser
   # Query: sum(rate(http_requests_total{service="api-gateway"}[5m])) by (status_code)
   ```

5. Verify that all routes are working:
   ```bash
   curl -I https://api.sparkfund.com/api/v1/users
   curl -I https://api.sparkfund.com/api/v1/kyc
   curl -I https://api.sparkfund.com/api/v1/investments
   ```

## Rollback Procedure

If issues are detected after deployment, follow these steps to rollback:

### 1. Automatic Rollback

The CI/CD pipeline will automatically rollback if the deployment fails or if the health checks fail.

### 2. Manual Rollback

To manually rollback:

```bash
# 1. Rollback the deployment
kubectl rollout undo deployment/api-gateway -n sparkfund

# 2. Monitor the rollback
kubectl rollout status deployment/api-gateway -n sparkfund

# 3. Verify that the previous version is running
kubectl describe deployment api-gateway -n sparkfund | grep Image
```

### 3. Canary Rollback

To rollback a canary deployment:

```bash
# 1. Delete the canary
kubectl delete canary api-gateway -n sparkfund

# 2. Ensure the primary deployment is using the previous version
kubectl set image deployment/api-gateway api-gateway=ghcr.io/sparkfund/api-gateway:previous-version -n sparkfund
```

## Deployment Schedule

- **Production Deployments**: Monday-Thursday, 10:00-16:00 ET
- **Emergency Deployments**: Any time, with approval from the on-call engineer and product owner
- **Blackout Periods**: End of month (28th-3rd), major holidays, and announced blackout periods

## Contacts

- **DevOps Team**: devops@sparkfund.com
- **On-Call Engineer**: oncall@sparkfund.com
- **Slack Channel**: #api-gateway-deployments
