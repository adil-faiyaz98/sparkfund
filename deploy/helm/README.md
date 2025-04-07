# SparkFund Helm Charts

This directory contains Helm charts for deploying the SparkFund platform to Kubernetes.

## Chart Structure

- `common/`: Common templates used by all service charts
- `api-gateway/`: API Gateway service chart
- `kyc-service/`: KYC Service chart
- `investment-service/`: Investment Service chart
- `user-service/`: User Service chart
- `ai-service/`: AI Service chart
- `sparkfund/`: Parent chart that includes all services

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- PV provisioner support in the underlying infrastructure
- Ingress controller (e.g., nginx-ingress)
- Cert-manager for TLS certificates

## Installation

### Installing the Chart

To install the chart with the release name `sparkfund`:

```bash
# Add the SparkFund Helm repository
helm repo add sparkfund https://charts.sparkfund.com

# Update the repository
helm repo update

# Install the chart
helm install sparkfund sparkfund/sparkfund
```

### Installing from Local Directory

```bash
# Install the chart from the local directory
helm install sparkfund ./sparkfund
```

### Installing with Custom Values

```bash
# Install the chart with custom values
helm install sparkfund ./sparkfund -f values-custom.yaml
```

### Installing for Different Environments

```bash
# Install for development environment
helm install sparkfund ./sparkfund -f sparkfund/values-dev.yaml

# Install for staging environment
helm install sparkfund ./sparkfund -f sparkfund/values-staging.yaml

# Install for production environment
helm install sparkfund ./sparkfund -f sparkfund/values-prod.yaml
```

## Upgrading the Chart

```bash
# Upgrade the chart
helm upgrade sparkfund ./sparkfund
```

## Uninstalling the Chart

```bash
# Uninstall the chart
helm uninstall sparkfund
```

## Configuration

The following table lists the configurable parameters of the SparkFund chart and their default values.

### Global Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `api-gateway.enabled` | Enable API Gateway | `true` |
| `kyc-service.enabled` | Enable KYC Service | `true` |
| `investment-service.enabled` | Enable Investment Service | `true` |
| `user-service.enabled` | Enable User Service | `true` |
| `ai-service.enabled` | Enable AI Service | `true` |

### API Gateway Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `api-gateway.replicaCount` | Number of replicas | `3` |
| `api-gateway.image.repository` | Image repository | `ghcr.io/adil-faiyaz98/sparkfund/api-gateway` |
| `api-gateway.image.tag` | Image tag | `latest` |
| `api-gateway.image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `api-gateway.service.type` | Service type | `ClusterIP` |
| `api-gateway.service.port` | Service port | `8080` |
| `api-gateway.ingress.enabled` | Enable ingress | `true` |
| `api-gateway.ingress.hosts` | Ingress hosts | `[{host: api.sparkfund.com, paths: [{path: /, pathType: Prefix}]}]` |
| `api-gateway.ingress.tls` | Ingress TLS configuration | `[{secretName: api-gateway-tls, hosts: [api.sparkfund.com]}]` |
| `api-gateway.resources` | Resource requests and limits | `{}` |
| `api-gateway.nodeSelector` | Node selector | `{}` |
| `api-gateway.tolerations` | Tolerations | `[]` |
| `api-gateway.affinity` | Affinity | `{}` |

### KYC Service Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `kyc-service.replicaCount` | Number of replicas | `3` |
| `kyc-service.image.repository` | Image repository | `ghcr.io/adil-faiyaz98/sparkfund/kyc-service` |
| `kyc-service.image.tag` | Image tag | `latest` |
| `kyc-service.image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `kyc-service.service.type` | Service type | `ClusterIP` |
| `kyc-service.service.port` | Service port | `8080` |
| `kyc-service.ingress.enabled` | Enable ingress | `true` |
| `kyc-service.ingress.hosts` | Ingress hosts | `[{host: kyc.sparkfund.com, paths: [{path: /, pathType: Prefix}]}]` |
| `kyc-service.ingress.tls` | Ingress TLS configuration | `[{secretName: kyc-service-tls, hosts: [kyc.sparkfund.com]}]` |
| `kyc-service.resources` | Resource requests and limits | `{}` |
| `kyc-service.nodeSelector` | Node selector | `{}` |
| `kyc-service.tolerations` | Tolerations | `[]` |
| `kyc-service.affinity` | Affinity | `{}` |

### Investment Service Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `investment-service.replicaCount` | Number of replicas | `3` |
| `investment-service.image.repository` | Image repository | `ghcr.io/adil-faiyaz98/sparkfund/investment-service` |
| `investment-service.image.tag` | Image tag | `latest` |
| `investment-service.image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `investment-service.service.type` | Service type | `ClusterIP` |
| `investment-service.service.port` | Service port | `8080` |
| `investment-service.ingress.enabled` | Enable ingress | `true` |
| `investment-service.ingress.hosts` | Ingress hosts | `[{host: investment.sparkfund.com, paths: [{path: /, pathType: Prefix}]}]` |
| `investment-service.ingress.tls` | Ingress TLS configuration | `[{secretName: investment-service-tls, hosts: [investment.sparkfund.com]}]` |
| `investment-service.resources` | Resource requests and limits | `{}` |
| `investment-service.nodeSelector` | Node selector | `{}` |
| `investment-service.tolerations` | Tolerations | `[]` |
| `investment-service.affinity` | Affinity | `{}` |

### User Service Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `user-service.replicaCount` | Number of replicas | `3` |
| `user-service.image.repository` | Image repository | `ghcr.io/adil-faiyaz98/sparkfund/user-service` |
| `user-service.image.tag` | Image tag | `latest` |
| `user-service.image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `user-service.service.type` | Service type | `ClusterIP` |
| `user-service.service.port` | Service port | `8080` |
| `user-service.ingress.enabled` | Enable ingress | `true` |
| `user-service.ingress.hosts` | Ingress hosts | `[{host: user.sparkfund.com, paths: [{path: /, pathType: Prefix}]}]` |
| `user-service.ingress.tls` | Ingress TLS configuration | `[{secretName: user-service-tls, hosts: [user.sparkfund.com]}]` |
| `user-service.resources` | Resource requests and limits | `{}` |
| `user-service.nodeSelector` | Node selector | `{}` |
| `user-service.tolerations` | Tolerations | `[]` |
| `user-service.affinity` | Affinity | `{}` |

### AI Service Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `ai-service.replicaCount` | Number of replicas | `2` |
| `ai-service.image.repository` | Image repository | `ghcr.io/adil-faiyaz98/sparkfund/ai-service` |
| `ai-service.image.tag` | Image tag | `latest` |
| `ai-service.image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `ai-service.service.type` | Service type | `ClusterIP` |
| `ai-service.service.port` | Service port | `8000` |
| `ai-service.ingress.enabled` | Enable ingress | `true` |
| `ai-service.ingress.hosts` | Ingress hosts | `[{host: ai.sparkfund.com, paths: [{path: /, pathType: Prefix}]}]` |
| `ai-service.ingress.tls` | Ingress TLS configuration | `[{secretName: ai-service-tls, hosts: [ai.sparkfund.com]}]` |
| `ai-service.resources` | Resource requests and limits | `{}` |
| `ai-service.nodeSelector` | Node selector | `{}` |
| `ai-service.tolerations` | Tolerations | `[]` |
| `ai-service.affinity` | Affinity | `{}` |
