# HELM

## 1. Install the Helm Chart

```bash
# For dev environment
helm install money-pulse ./helm/money-pulse -f ./helm/money-pulse/values.yaml --namespace money-pulse --create-namespace

# For production environment
helm install money-pulse ./helm/money-pulse -f ./helm/money-pulse/values-prod.yaml --namespace money-pulse --create-namespace
```

## 2. Upgrade the Helm Chart

```bash
helm upgrade money-pulse ./helm/money-pulse -f ./helm/money-pulse/values.yaml --namespace money-pulse
```

## 3. Uninstall the Helm Chart

```bash
helm uninstall money-pulse --namespace money-pulse
```

## 4. Package the Helm Chart

```bash
helm uninstall money-pulse --namespace money-pulse
```
