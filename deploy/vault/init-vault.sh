#!/bin/bash

# This script initializes and configures Vault for the SparkFund platform

# Set variables
VAULT_ADDR="https://vault.sparkfund.com"
VAULT_TOKEN_FILE="vault-token.json"
VAULT_KEYS_FILE="vault-keys.json"
VAULT_CONFIG_DIR="$(dirname "$0")/config"

# Check if Vault is already initialized
echo "Checking if Vault is already initialized..."
INIT_STATUS=$(curl -s -k $VAULT_ADDR/v1/sys/init | jq -r '.initialized')

if [ "$INIT_STATUS" == "true" ]; then
  echo "Vault is already initialized."
else
  echo "Initializing Vault..."
  # Initialize Vault with 5 key shares and 3 key threshold
  curl -s -k -X PUT -d '{"secret_shares": 5, "secret_threshold": 3}' $VAULT_ADDR/v1/sys/init > $VAULT_KEYS_FILE
  
  # Extract root token and unseal keys
  ROOT_TOKEN=$(cat $VAULT_KEYS_FILE | jq -r '.root_token')
  UNSEAL_KEY_1=$(cat $VAULT_KEYS_FILE | jq -r '.keys[0]')
  UNSEAL_KEY_2=$(cat $VAULT_KEYS_FILE | jq -r '.keys[1]')
  UNSEAL_KEY_3=$(cat $VAULT_KEYS_FILE | jq -r '.keys[2]')
  
  echo "Vault initialized successfully."
  echo "Root Token: $ROOT_TOKEN"
  echo "Unseal Keys: saved to $VAULT_KEYS_FILE"
  
  # Save root token to file
  echo "{\"token\": \"$ROOT_TOKEN\"}" > $VAULT_TOKEN_FILE
  
  # Unseal Vault
  echo "Unsealing Vault..."
  curl -s -k -X PUT -d "{\"key\": \"$UNSEAL_KEY_1\"}" $VAULT_ADDR/v1/sys/unseal
  curl -s -k -X PUT -d "{\"key\": \"$UNSEAL_KEY_2\"}" $VAULT_ADDR/v1/sys/unseal
  curl -s -k -X PUT -d "{\"key\": \"$UNSEAL_KEY_3\"}" $VAULT_ADDR/v1/sys/unseal
  
  echo "Vault unsealed successfully."
fi

# Set Vault token for subsequent commands
ROOT_TOKEN=$(cat $VAULT_TOKEN_FILE | jq -r '.token')
export VAULT_TOKEN=$ROOT_TOKEN

# Enable Kubernetes authentication
echo "Enabling Kubernetes authentication..."
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{"type": "kubernetes"}' $VAULT_ADDR/v1/sys/auth/kubernetes

# Configure Kubernetes authentication
echo "Configuring Kubernetes authentication..."
KUBE_CA_CERT=$(kubectl config view --raw --minify --flatten --output='jsonpath={.clusters[].cluster.certificate-authority-data}' | base64 --decode)
KUBE_HOST=$(kubectl config view --raw --minify --flatten --output='jsonpath={.clusters[].cluster.server}')
TOKEN_REVIEWER_JWT=$(kubectl create token vault-auth -n vault --duration=87600h)

curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{
  \"kubernetes_host\": \"$KUBE_HOST\",
  \"kubernetes_ca_cert\": \"$KUBE_CA_CERT\",
  \"token_reviewer_jwt\": \"$TOKEN_REVIEWER_JWT\"
}" $VAULT_ADDR/v1/auth/kubernetes/config

# Enable secrets engines
echo "Enabling secrets engines..."
# KV version 2 for static secrets
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{"type": "kv", "options": {"version": "2"}}' $VAULT_ADDR/v1/sys/mounts/kv

# Database secrets engine for dynamic database credentials
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{"type": "database"}' $VAULT_ADDR/v1/sys/mounts/database

# Transit engine for encryption as a service
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{"type": "transit"}' $VAULT_ADDR/v1/sys/mounts/transit

# PKI engine for certificate management
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{"type": "pki"}' $VAULT_ADDR/v1/sys/mounts/pki
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{"max_lease_ttl": "87600h"}' $VAULT_ADDR/v1/sys/mounts/pki/tune

# Configure database secrets engine
echo "Configuring database secrets engine..."
# PostgreSQL connection
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "plugin_name": "postgresql-database-plugin",
  "allowed_roles": "kyc-service,investment-service,user-service",
  "connection_url": "postgresql://{{username}}:{{password}}@postgres:5432/sparkfund?sslmode=disable",
  "username": "postgres",
  "password": "postgres"
}' $VAULT_ADDR/v1/database/config/postgres

# Create database roles
echo "Creating database roles..."
# KYC Service role
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "db_name": "postgres",
  "creation_statements": "CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD \"{{password}}\" VALID UNTIL \"{{expiration}}\"; GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO \"{{name}}\";",
  "default_ttl": "1h",
  "max_ttl": "24h"
}' $VAULT_ADDR/v1/database/roles/kyc-service

# Investment Service role
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "db_name": "postgres",
  "creation_statements": "CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD \"{{password}}\" VALID UNTIL \"{{expiration}}\"; GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO \"{{name}}\";",
  "default_ttl": "1h",
  "max_ttl": "24h"
}' $VAULT_ADDR/v1/database/roles/investment-service

# User Service role
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "db_name": "postgres",
  "creation_statements": "CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD \"{{password}}\" VALID UNTIL \"{{expiration}}\"; GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO \"{{name}}\";",
  "default_ttl": "1h",
  "max_ttl": "24h"
}' $VAULT_ADDR/v1/database/roles/user-service

# Configure transit secrets engine
echo "Configuring transit secrets engine..."
# Create encryption keys
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "type": "aes256-gcm96"
}' $VAULT_ADDR/v1/transit/keys/kyc-data

curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "type": "aes256-gcm96"
}' $VAULT_ADDR/v1/transit/keys/investment-data

curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "type": "aes256-gcm96"
}' $VAULT_ADDR/v1/transit/keys/user-data

# Configure PKI secrets engine
echo "Configuring PKI secrets engine..."
# Generate root CA
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "common_name": "SparkFund Root CA",
  "ttl": "87600h"
}' $VAULT_ADDR/v1/pki/root/generate/internal

# Configure PKI URLs
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{
  \"issuing_certificates\": \"$VAULT_ADDR/v1/pki/ca\",
  \"crl_distribution_points\": \"$VAULT_ADDR/v1/pki/crl\"
}" $VAULT_ADDR/v1/pki/config/urls

# Create PKI role
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "allowed_domains": "sparkfund.com,sparkfund.svc",
  "allow_subdomains": true,
  "max_ttl": "720h"
}' $VAULT_ADDR/v1/pki/roles/sparkfund-dot-com

# Create Kubernetes auth roles
echo "Creating Kubernetes auth roles..."
# KYC Service role
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "bound_service_account_names": "kyc-service",
  "bound_service_account_namespaces": "sparkfund-dev,sparkfund-staging,sparkfund-prod",
  "policies": "kyc-service",
  "ttl": "1h"
}' $VAULT_ADDR/v1/auth/kubernetes/role/kyc-service

# Investment Service role
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "bound_service_account_names": "investment-service",
  "bound_service_account_namespaces": "sparkfund-dev,sparkfund-staging,sparkfund-prod",
  "policies": "investment-service",
  "ttl": "1h"
}' $VAULT_ADDR/v1/auth/kubernetes/role/investment-service

# User Service role
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "bound_service_account_names": "user-service",
  "bound_service_account_namespaces": "sparkfund-dev,sparkfund-staging,sparkfund-prod",
  "policies": "user-service",
  "ttl": "1h"
}' $VAULT_ADDR/v1/auth/kubernetes/role/user-service

# API Gateway role
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "bound_service_account_names": "api-gateway",
  "bound_service_account_namespaces": "sparkfund-dev,sparkfund-staging,sparkfund-prod",
  "policies": "api-gateway",
  "ttl": "1h"
}' $VAULT_ADDR/v1/auth/kubernetes/role/api-gateway

# Create policies
echo "Creating policies..."
# KYC Service policy
curl -s -k -X PUT -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "policy": "path \"database/creds/kyc-service\" {\n  capabilities = [\"read\"]\n}\n\npath \"transit/encrypt/kyc-data\" {\n  capabilities = [\"update\"]\n}\n\npath \"transit/decrypt/kyc-data\" {\n  capabilities = [\"update\"]\n}\n\npath \"kv/data/kyc-service/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\"]\n}\n\npath \"pki/issue/sparkfund-dot-com\" {\n  capabilities = [\"create\", \"update\"]\n}"
}' $VAULT_ADDR/v1/sys/policies/acl/kyc-service

# Investment Service policy
curl -s -k -X PUT -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "policy": "path \"database/creds/investment-service\" {\n  capabilities = [\"read\"]\n}\n\npath \"transit/encrypt/investment-data\" {\n  capabilities = [\"update\"]\n}\n\npath \"transit/decrypt/investment-data\" {\n  capabilities = [\"update\"]\n}\n\npath \"kv/data/investment-service/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\"]\n}\n\npath \"pki/issue/sparkfund-dot-com\" {\n  capabilities = [\"create\", \"update\"]\n}"
}' $VAULT_ADDR/v1/sys/policies/acl/investment-service

# User Service policy
curl -s -k -X PUT -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "policy": "path \"database/creds/user-service\" {\n  capabilities = [\"read\"]\n}\n\npath \"transit/encrypt/user-data\" {\n  capabilities = [\"update\"]\n}\n\npath \"transit/decrypt/user-data\" {\n  capabilities = [\"update\"]\n}\n\npath \"kv/data/user-service/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\"]\n}\n\npath \"pki/issue/sparkfund-dot-com\" {\n  capabilities = [\"create\", \"update\"]\n}"
}' $VAULT_ADDR/v1/sys/policies/acl/user-service

# API Gateway policy
curl -s -k -X PUT -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "policy": "path \"kv/data/api-gateway/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\"]\n}\n\npath \"pki/issue/sparkfund-dot-com\" {\n  capabilities = [\"create\", \"update\"]\n}"
}' $VAULT_ADDR/v1/sys/policies/acl/api-gateway

# Store static secrets
echo "Storing static secrets..."
# KYC Service secrets
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "data": {
    "api-key": "kyc-service-api-key-12345",
    "jwt-secret": "kyc-service-jwt-secret-12345"
  }
}' $VAULT_ADDR/v1/kv/data/kyc-service/config

# Investment Service secrets
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "data": {
    "api-key": "investment-service-api-key-12345",
    "jwt-secret": "investment-service-jwt-secret-12345"
  }
}' $VAULT_ADDR/v1/kv/data/investment-service/config

# User Service secrets
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "data": {
    "api-key": "user-service-api-key-12345",
    "jwt-secret": "user-service-jwt-secret-12345"
  }
}' $VAULT_ADDR/v1/kv/data/user-service/config

# API Gateway secrets
curl -s -k -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{
  "data": {
    "api-key": "api-gateway-api-key-12345",
    "jwt-secret": "api-gateway-jwt-secret-12345"
  }
}' $VAULT_ADDR/v1/kv/data/api-gateway/config

echo "Vault configuration completed successfully."
