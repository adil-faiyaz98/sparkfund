Here is the **pure text format** version of your Helm commands for the `README.md` file:

---

# HELM

## 1. Install the Helm Chart

For development environment:

```
helm install money-pulse ./helm/money-pulse -f ./helm/money-pulse/values.yaml --namespace money-pulse --create-namespace
```

For production environment:

```
helm install money-pulse ./helm/money-pulse -f ./helm/money-pulse/values-prod.yaml --namespace money-pulse --create-namespace
```

## 2. Upgrade the Helm Chart

```
helm upgrade money-pulse ./helm/money-pulse -f ./helm/money-pulse/values.yaml --namespace money-pulse
```

## 3. Uninstall the Helm Chart

```
helm uninstall money-pulse --namespace money-pulse
```

## 4. Package the Helm Chart

```
helm package ./helm/money-pulse
```

## 5. Add a Chart Repository

```
helm repo add bitnami https://charts.bitnami.com/bitnami
```

## 6. Update Repositories

```
helm repo update
```

## 7. List Repositories

```
helm repo list
```

## 8. Remove a Repository

```
helm repo remove bitnami
```

## 9. Search for PostgreSQL Charts

```
helm search repo postgresql
```

## 10. Search for Charts with Detailed Output

```
helm search repo nginx -l
```

## 11. List All Releases in a Specific Namespace

```
helm list -n <app>
```

## 12. List All Releases in All Namespaces

```
helm list --all-namespaces
```

## 13. List Releases Including Deleted Ones with Pending Purge

```
helm list --all -n <app>
```

## 14. View Release History

```
helm history money-pulse -n money-pulse
```

## 15. Rollback to a Specific Revision

```
helm rollback money-pulse 1 -n money-pulse
```

## 16. Rollback with a Custom Description

```
helm rollback money-pulse 2 -n money-pulse --description "Rolling back due to database issues"
```

## 17. Get All Information About a Release

```
helm get all money-pulse -n money-pulse
```

## 18. Get Only the Values

```
helm get values money-pulse -n money-pulse
```

## 19. Get Only the Manifest

```
helm get manifest money-pulse -n money-pulse
```

## 20. Get the Notes

```
helm get notes money-pulse -n money-pulse
```

## 21. Get the Hooks

```
helm get hooks money-pulse -n money-pulse
```

## 22. Render Chart Templates Locally Without Installing

```
helm template money-pulse ./helm/money-pulse -f ./helm/money-pulse/values-dev.yaml
```

## 23. Validate Chart Structure and Formatting

```
helm lint ./helm/money-pulse
```

## 24. Debug Template Rendering Issues

```
helm template --debug money-pulse ./helm/money-pulse
```

## 25. Update Chart Dependencies

```
helm dependency update ./helm/money-pulse
```

## 26. List Chart Dependencies

```
helm dependency list ./helm/money-pulse
```

## 27. Build Chart Dependencies

```
helm dependency build ./helm/money-pulse
```

## 28. Show Chart Information

```
helm show chart ./helm/money-pulse
```

## 29. Show Chart Values

```
helm show values ./helm/money-pulse
```

## 30. Show Chart README

```
helm show readme ./helm/money-pulse
```

## 31. Show All Chart Information

```
helm show all ./helm/money-pulse
```

## 32. Run Tests for a Release

```
helm test money-pulse -n money-pulse
```

## 33. Check Release Status

```
helm status money-pulse -n money-pulse
```

## 34. Download a Chart for Customization

```
helm pull bitnami/postgresql --untar --untardir ./custom-charts
```

## 35. Dry Run to Check What Would Be Installed

```
helm install --dry-run --debug money-pulse ./helm/money-pulse
```

## 36. Install with Timeout Override

```
helm install money-pulse ./helm/money-pulse --timeout 10m
```

## 37. Install with Specific Kubernetes Version

```
helm install money-pulse ./helm/money-pulse --kube-version 1.25.0
```

## 38. Install with Custom Release Name

```
helm install finance-app ./helm/money-pulse --namespace money-pulse
```

## 39. Set Helm Configuration with Environment Variables

```
export HELM_NAMESPACE=money-pulse
export HELM_KUBECONTEXT=my-cluster-context
```

## 40. Then run Helm commands without specifying the namespace:

```
helm list # Will use HELM_NAMESPACE
```
