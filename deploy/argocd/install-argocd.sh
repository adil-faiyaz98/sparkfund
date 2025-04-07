#!/bin/bash

# This script installs ArgoCD and configures it for the SparkFund platform

# Create ArgoCD namespace
kubectl create namespace argocd

# Install ArgoCD
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for ArgoCD to be ready
echo "Waiting for ArgoCD to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/argocd-server -n argocd

# Apply ArgoCD Project
kubectl apply -f projects/sparkfund-project.yaml

# Apply infrastructure applications
kubectl apply -f applications/infrastructure/cert-manager.yaml
kubectl apply -f applications/infrastructure/ingress-nginx.yaml
kubectl apply -f applications/infrastructure/prometheus-stack.yaml

# Wait for infrastructure to be ready
echo "Waiting for infrastructure to be ready..."
sleep 60

# Apply SparkFund applications
kubectl apply -f applications/dev/sparkfund-dev.yaml
kubectl apply -f applications/staging/sparkfund-staging.yaml
kubectl apply -f applications/prod/sparkfund-prod.yaml

# Get ArgoCD admin password
echo "ArgoCD admin password:"
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
echo

# Port forward ArgoCD server
echo "Port forwarding ArgoCD server to http://localhost:8080"
echo "Press Ctrl+C to stop"
kubectl port-forward svc/argocd-server -n argocd 8080:443
