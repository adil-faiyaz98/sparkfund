# Money Pulse Kubernetes Deployment

This directory contains the Kubernetes manifests for deploying the Money Pulse application.

## Directory Structure

```sh
k8s/
    base/
        namespace.yaml
        secrets/
        configmaps/
        storage/
    infrastructure/
    services/
    scaling/
    overlays/
        dev/
        prod/
```

## Deployment Commands

### Development Environment

### Apply development configuration

```bash
kubectl apply -k overlays/dev
```

### View all resources in the money-pulse namespace

```bash
kubectl get all -n money-pulse
```

### Port forward to access services locally

```bash
kubectl port-forward svc/accounts-service -n money-pulse 8080:80
```

## Accessing Services

### Development

```bash
Access the API at http://api.money-pulse.local/accounts (add to your hosts file)
```

### Production

```bash
Access the API at https://api.money-pulse.com/accounts
```

These Kubernetes manifests provide a complete deployment solution for the application with:

1. **Base resources** - Common configurations like namespace, ConfigMaps, and PVCs
2. **Infrastructure components** - Database and ingress controller
3. **Service deployments** - All microservices with proper configuration
4. **Scaling policies** - Horizontal Pod Autoscalers for dynamic scaling
5. **Environment overlays** - Kustomize configurations for dev and prod environments
