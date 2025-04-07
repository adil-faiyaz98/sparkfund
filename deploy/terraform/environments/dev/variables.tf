variable "region" {
  description = "The AWS region to deploy to"
  type        = string
  default     = "us-east-1"
}

variable "vpc_cidr" {
  description = "The CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "az_count" {
  description = "The number of availability zones to use"
  type        = number
  default     = 3
}

variable "cluster_name" {
  description = "The name of the EKS cluster"
  type        = string
  default     = "sparkfund"
}

variable "kubernetes_version" {
  description = "The Kubernetes version to use"
  type        = string
  default     = "1.27"
}

variable "node_groups" {
  description = "Map of EKS node group configurations"
  type        = map(any)
  default = {
    default = {
      desired_size = 2
      min_size     = 1
      max_size     = 3
      instance_types = ["t3.medium"]
      capacity_type  = "ON_DEMAND"
      disk_size      = 20
      labels = {
        "role" = "default"
      }
    }
  }
}

variable "rds_instance_class" {
  description = "The instance class to use for RDS"
  type        = string
  default     = "db.t3.medium"
}

variable "rds_allocated_storage" {
  description = "The amount of storage to allocate for RDS (in GB)"
  type        = number
  default     = 20
}

variable "rds_max_allocated_storage" {
  description = "The maximum amount of storage to allocate for RDS (in GB)"
  type        = number
  default     = 100
}

variable "rds_multi_az" {
  description = "Whether to deploy RDS in multiple availability zones"
  type        = bool
  default     = false
}

variable "rds_deletion_protection" {
  description = "Whether to enable deletion protection for RDS"
  type        = bool
  default     = false
}

variable "rds_skip_final_snapshot" {
  description = "Whether to skip the final snapshot when RDS is deleted"
  type        = bool
  default     = true
}

variable "elasticache_node_type" {
  description = "The node type to use for ElastiCache"
  type        = string
  default     = "cache.t3.small"
}

variable "elasticache_num_cache_clusters" {
  description = "The number of cache clusters for ElastiCache"
  type        = number
  default     = 2
}

variable "elasticache_automatic_failover_enabled" {
  description = "Whether to enable automatic failover for ElastiCache"
  type        = bool
  default     = true
}

variable "elasticache_multi_az_enabled" {
  description = "Whether to enable multi-AZ for ElastiCache"
  type        = bool
  default     = false
}
