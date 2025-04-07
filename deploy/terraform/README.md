# SparkFund Terraform Infrastructure

This directory contains Terraform configurations for provisioning the infrastructure required by the SparkFund platform.

## Directory Structure

- `modules/`: Reusable Terraform modules
  - `eks/`: Amazon EKS cluster module
  - `vpc/`: VPC and networking module
  - `rds/`: RDS database module
  - `elasticache/`: ElastiCache (Redis) module
  - `s3/`: S3 bucket module
  - `iam/`: IAM roles and policies module
- `environments/`: Environment-specific configurations
  - `dev/`: Development environment
  - `staging/`: Staging environment
  - `prod/`: Production environment

## Prerequisites

- Terraform 1.0.0+
- AWS CLI configured with appropriate credentials
- kubectl installed
- helm installed

## Usage

### Initialize Terraform

```bash
cd environments/dev
terraform init
```

### Plan the Infrastructure

```bash
terraform plan -out=tfplan
```

### Apply the Infrastructure

```bash
terraform apply tfplan
```

### Destroy the Infrastructure

```bash
terraform destroy
```

## Modules

### EKS Module

The EKS module provisions an Amazon EKS cluster with the following features:

- Managed node groups
- Cluster autoscaler
- AWS Load Balancer Controller
- Prometheus and Grafana for monitoring
- Fluent Bit for logging
- Cert Manager for TLS certificates

### VPC Module

The VPC module provisions a VPC with the following features:

- Public and private subnets
- NAT gateways
- Internet gateway
- VPC endpoints for AWS services

### RDS Module

The RDS module provisions an Amazon RDS instance with the following features:

- PostgreSQL database
- Multi-AZ deployment
- Automated backups
- Encryption at rest

### ElastiCache Module

The ElastiCache module provisions an Amazon ElastiCache cluster with the following features:

- Redis engine
- Multi-AZ deployment
- Encryption at rest
- Automatic failover

### S3 Module

The S3 module provisions an Amazon S3 bucket with the following features:

- Versioning
- Encryption at rest
- Lifecycle policies
- Access control

### IAM Module

The IAM module provisions IAM roles and policies with the following features:

- Least privilege principle
- Service accounts for Kubernetes
- IRSA (IAM Roles for Service Accounts)
- Boundary policies
