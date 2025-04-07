# SparkFund Terraform Environments

This directory contains Terraform configurations for different environments.

## Environments

- `dev/`: Development environment
- `staging/`: Staging environment
- `prod/`: Production environment

## Usage

### Initialize Terraform

```bash
cd dev
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

## Environment Variables

Each environment has its own set of variables defined in `terraform.tfvars`. These variables are used to configure the infrastructure for the specific environment.

### Common Variables

- `region`: The AWS region to deploy to
- `vpc_cidr`: The CIDR block for the VPC
- `az_count`: The number of availability zones to use
- `cluster_name`: The name of the EKS cluster
- `kubernetes_version`: The Kubernetes version to use

### Node Groups

The `node_groups` variable is a map of EKS node group configurations. Each node group has the following properties:

- `desired_size`: The desired number of nodes
- `min_size`: The minimum number of nodes
- `max_size`: The maximum number of nodes
- `instance_types`: The instance types to use
- `capacity_type`: The capacity type (ON_DEMAND or SPOT)
- `disk_size`: The disk size in GB
- `labels`: A map of Kubernetes labels to apply to the nodes

### RDS

The RDS variables configure the PostgreSQL database:

- `rds_instance_class`: The instance class to use
- `rds_allocated_storage`: The amount of storage to allocate (in GB)
- `rds_max_allocated_storage`: The maximum amount of storage to allocate (in GB)
- `rds_multi_az`: Whether to deploy in multiple availability zones
- `rds_deletion_protection`: Whether to enable deletion protection
- `rds_skip_final_snapshot`: Whether to skip the final snapshot when the RDS instance is deleted

### ElastiCache

The ElastiCache variables configure the Redis cluster:

- `elasticache_node_type`: The node type to use
- `elasticache_num_cache_clusters`: The number of cache clusters
- `elasticache_automatic_failover_enabled`: Whether to enable automatic failover
- `elasticache_multi_az_enabled`: Whether to enable multi-AZ deployment

## Outputs

Each environment exports the following outputs:

- `vpc_id`: The ID of the VPC
- `private_subnet_ids`: The IDs of the private subnets
- `public_subnet_ids`: The IDs of the public subnets
- `eks_cluster_id`: The ID of the EKS cluster
- `eks_cluster_endpoint`: The endpoint of the EKS cluster
- `eks_cluster_security_group_id`: The security group ID of the EKS cluster
- `eks_node_security_group_id`: The security group ID of the EKS node group
- `eks_oidc_provider_arn`: The ARN of the OIDC provider
- `rds_instance_endpoint`: The endpoint of the RDS instance
- `rds_instance_address`: The address of the RDS instance
- `rds_security_group_id`: The security group ID of the RDS instance
- `elasticache_primary_endpoint_address`: The address of the primary endpoint for ElastiCache
- `elasticache_reader_endpoint_address`: The address of the reader endpoint for ElastiCache
- `elasticache_security_group_id`: The security group ID of the ElastiCache cluster
