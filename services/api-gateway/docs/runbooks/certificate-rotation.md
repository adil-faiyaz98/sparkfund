# Certificate Rotation Runbook

This runbook describes the process for rotating TLS certificates for the API Gateway service.

## Overview

The API Gateway uses TLS certificates for:
1. Securing external traffic (HTTPS)
2. Mutual TLS (mTLS) for service-to-service communication

Certificates should be rotated:
- Before expiration (typically every 90 days)
- After a security incident
- When changing certificate authorities
- When changing domain names

## Prerequisites

- Access to the Kubernetes cluster
- `kubectl` configured to access the cluster
- Access to cert-manager resources
- Access to the certificate authority (if using a private CA)

## Certificate Types

### 1. External TLS Certificate

This certificate secures traffic between clients and the API Gateway.

- **Secret Name**: `api-gateway-tls`
- **Certificate Name**: `api-gateway-cert`
- **Domains**: `api.sparkfund.com`
- **Issuer**: `sparkfund-issuer` (Let's Encrypt in production)

### 2. Internal mTLS Certificate

This certificate secures traffic between the API Gateway and backend services.

- **Secret Name**: `api-gateway-mtls`
- **Certificate Name**: `api-gateway-mtls-cert`
- **Domains**: `api-gateway.sparkfund.svc.cluster.local`
- **Issuer**: `sparkfund-internal-issuer` (Internal CA)

## Automatic Rotation with cert-manager

The API Gateway uses cert-manager for automatic certificate rotation. cert-manager will automatically renew certificates before they expire.

### Checking Certificate Status

```bash
# Check external certificate
kubectl get certificate api-gateway-cert -n sparkfund -o wide

# Check internal certificate
kubectl get certificate api-gateway-mtls-cert -n sparkfund -o wide

# Check certificate details
kubectl describe certificate api-gateway-cert -n sparkfund
```

### Triggering Manual Rotation

If you need to manually rotate a certificate:

```bash
# Annotate the certificate for renewal
kubectl annotate certificate api-gateway-cert -n sparkfund cert-manager.io/renew="true"

# Monitor the renewal process
kubectl get certificaterequest -n sparkfund -l cert-manager.io/certificate-name=api-gateway-cert
```

## Manual Certificate Rotation

If cert-manager is unavailable or you need to use a certificate from an external source:

### 1. Prepare the New Certificate

Ensure you have:
- Certificate file (tls.crt)
- Private key file (tls.key)
- CA certificate (ca.crt) for mTLS

### 2. Create a New Secret

```bash
# For external TLS
kubectl create secret tls api-gateway-tls-new -n sparkfund \
  --cert=path/to/tls.crt \
  --key=path/to/tls.key \
  --dry-run=client -o yaml | kubectl apply -f -

# For internal mTLS
kubectl create secret generic api-gateway-mtls-new -n sparkfund \
  --from-file=tls.crt=path/to/tls.crt \
  --from-file=tls.key=path/to/tls.key \
  --from-file=ca.crt=path/to/ca.crt \
  --dry-run=client -o yaml | kubectl apply -f -
```

### 3. Update References

Update the API Gateway deployment to use the new secret:

```bash
# Edit the deployment
kubectl edit deployment api-gateway -n sparkfund

# Find the volume mounts and update the secret name
# volumes:
# - name: tls-cert
#   secret:
#     secretName: api-gateway-tls-new
```

### 4. Restart the API Gateway

```bash
kubectl rollout restart deployment api-gateway -n sparkfund
```

### 5. Verify the New Certificate

```bash
# Check the certificate in use
kubectl exec -it $(kubectl get pod -n sparkfund -l app=api-gateway -o jsonpath='{.items[0].metadata.name}') -n sparkfund -- \
  openssl s_client -connect localhost:8443 -showcerts

# Verify the expiration date
kubectl exec -it $(kubectl get pod -n sparkfund -l app=api-gateway -o jsonpath='{.items[0].metadata.name}') -n sparkfund -- \
  openssl x509 -in /etc/ssl/certs/tls.crt -noout -dates
```

### 6. Clean Up

Once the new certificate is working correctly:

```bash
# Rename the new secret to the standard name
kubectl get secret api-gateway-tls-new -n sparkfund -o yaml | \
  sed 's/name: api-gateway-tls-new/name: api-gateway-tls/' | \
  kubectl apply -f -

# Delete the old secret
kubectl delete secret api-gateway-tls-old -n sparkfund
```

## Rotating CA Certificates

If you need to rotate the CA certificate:

### 1. Update the CA ConfigMap

```bash
# Update the CA ConfigMap
kubectl create configmap api-gateway-ca -n sparkfund \
  --from-file=ca.crt=path/to/new-ca.crt \
  --dry-run=client -o yaml | kubectl apply -f -
```

### 2. Update the Issuer

```bash
# Update the issuer with the new CA
kubectl edit issuer sparkfund-internal-issuer -n sparkfund
```

### 3. Rotate Service Certificates

Follow the steps above to rotate all service certificates that were signed by the old CA.

## Emergency Certificate Rotation

In case of a security incident (e.g., private key compromise):

### 1. Revoke the Compromised Certificate

If using Let's Encrypt:

```bash
# Using certbot
certbot revoke --cert-path /path/to/compromised-cert.pem --key-path /path/to/compromised-key.pem
```

### 2. Generate a New Certificate Immediately

```bash
# Force renewal
kubectl annotate certificate api-gateway-cert -n sparkfund cert-manager.io/renew="true"

# Monitor the renewal process
kubectl get certificaterequest -n sparkfund -l cert-manager.io/certificate-name=api-gateway-cert
```

### 3. Update All Clients

Ensure all clients have the new CA certificate in their trust store.

### 4. Rotate Related Secrets

If the certificate was used for authentication, rotate any related credentials.

## Troubleshooting

### Certificate Not Renewing

1. Check cert-manager logs:
   ```bash
   kubectl logs -n cert-manager -l app=cert-manager -c cert-manager
   ```

2. Check certificate events:
   ```bash
   kubectl describe certificate api-gateway-cert -n sparkfund
   ```

3. Check certificate request:
   ```bash
   kubectl get certificaterequest -n sparkfund -l cert-manager.io/certificate-name=api-gateway-cert
   ```

### Certificate Validation Failures

1. Check the certificate chain:
   ```bash
   kubectl get secret api-gateway-tls -n sparkfund -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -noout -text
   ```

2. Verify the certificate matches the domain:
   ```bash
   kubectl get secret api-gateway-tls -n sparkfund -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -noout -subject
   ```

3. Check for expiration:
   ```bash
   kubectl get secret api-gateway-tls -n sparkfund -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -noout -dates
   ```

## Contacts

- **Security Team**: security@sparkfund.com
- **DevOps Team**: devops@sparkfund.com
- **On-Call Engineer**: oncall@sparkfund.com
- **Slack Channel**: #api-gateway-security
