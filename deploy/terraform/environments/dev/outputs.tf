output "vpc_id" {
  description = "The ID of the VPC"
  value       = module.vpc.vpc_id
}

output "private_subnet_ids" {
  description = "The IDs of the private subnets"
  value       = module.vpc.private_subnet_ids
}

output "public_subnet_ids" {
  description = "The IDs of the public subnets"
  value       = module.vpc.public_subnet_ids
}

output "eks_cluster_id" {
  description = "The ID of the EKS cluster"
  value       = module.eks.cluster_id
}

output "eks_cluster_endpoint" {
  description = "The endpoint of the EKS cluster"
  value       = module.eks.cluster_endpoint
}

output "eks_cluster_security_group_id" {
  description = "The security group ID of the EKS cluster"
  value       = module.eks.cluster_security_group_id
}

output "eks_node_security_group_id" {
  description = "The security group ID of the EKS node group"
  value       = module.eks.node_security_group_id
}

output "eks_oidc_provider_arn" {
  description = "The ARN of the OIDC provider"
  value       = module.eks.oidc_provider_arn
}

output "rds_instance_endpoint" {
  description = "The endpoint of the RDS instance"
  value       = module.rds.instance_endpoint
}

output "rds_instance_address" {
  description = "The address of the RDS instance"
  value       = module.rds.instance_address
}

output "rds_security_group_id" {
  description = "The security group ID of the RDS instance"
  value       = module.rds.security_group_id
}

output "elasticache_primary_endpoint_address" {
  description = "The address of the primary endpoint for ElastiCache"
  value       = module.elasticache.primary_endpoint_address
}

output "elasticache_reader_endpoint_address" {
  description = "The address of the reader endpoint for ElastiCache"
  value       = module.elasticache.reader_endpoint_address
}

output "elasticache_security_group_id" {
  description = "The security group ID of the ElastiCache cluster"
  value       = module.elasticache.security_group_id
}
