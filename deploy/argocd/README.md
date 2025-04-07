# SparkFund ArgoCD Configuration

This directory contains ArgoCD configuration for deploying the SparkFund platform using GitOps principles.

## Overview

ArgoCD is a declarative, GitOps continuous delivery tool for Kubernetes. It follows the GitOps pattern of using Git repositories as the source of truth for defining the desired application state. ArgoCD automates the deployment of the desired application states in the specified target environments.

## Directory Structure

- `applications/`: ArgoCD Application resources
  - `dev/`: Development environment applications
  - `staging/`: Staging environment applications
  - `prod/`: Production environment applications
- `projects/`: ArgoCD Project resources
- `applicationsets/`: ArgoCD ApplicationSet resources

## Prerequisites

- Kubernetes cluster with ArgoCD installed
- kubectl configured to communicate with the cluster
- Git repository with Helm charts

## Installation

### Installing ArgoCD

```bash
# Create ArgoCD namespace
kubectl create namespace argocd

# Install ArgoCD
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Access ArgoCD UI
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

### Deploying Applications

```bash
# Apply ArgoCD Project
kubectl apply -f projects/sparkfund-project.yaml

# Apply ArgoCD ApplicationSet
kubectl apply -f applicationsets/sparkfund-applicationset.yaml
```

## Usage

### Accessing ArgoCD UI

The ArgoCD UI can be accessed at https://localhost:8080 when port-forwarding is active.

Default credentials:
- Username: admin
- Password: (retrieve with `kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d`)

### Syncing Applications

Applications can be synced manually through the ArgoCD UI or automatically based on the sync policy defined in the Application resource.

### Adding New Applications

To add a new application:

1. Create a new Application resource in the appropriate environment directory
2. Apply the Application resource to the cluster
3. Verify the application is synced in the ArgoCD UI

## GitOps Workflow

1. Developers make changes to the application code and push to the Git repository
2. CI pipeline builds, tests, and pushes the application image to the container registry
3. Developers update the Helm chart values with the new image tag and push to the Git repository
4. ArgoCD detects the changes in the Git repository and automatically syncs the application to the target environment

## Best Practices

- Use separate branches for different environments (e.g., dev, staging, prod)
- Use semantic versioning for application images
- Use Helm chart values for environment-specific configuration
- Use ArgoCD Projects to enforce RBAC and resource constraints
- Use ApplicationSets for managing multiple similar applications
- Use sync waves for controlling the order of resource creation
- Use health checks to verify application health after deployment
