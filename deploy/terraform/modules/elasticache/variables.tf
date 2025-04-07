variable "region" {
  description = "The AWS region to deploy to"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "The environment (dev, staging, prod)"
  type        = string
}

variable "name" {
  description = "The name of the ElastiCache cluster"
  type        = string
}

variable "vpc_id" {
  description = "The ID of the VPC"
  type        = string
}

variable "subnet_ids" {
  description = "The IDs of the subnets where the ElastiCache cluster will be deployed"
  type        = list(string)
}

variable "allowed_security_groups" {
  description = "The security groups allowed to access the ElastiCache cluster"
  type        = list(string)
  default     = []
}

variable "kms_key_arn" {
  description = "The ARN of the KMS key to use for encryption"
  type        = string
  default     = ""
}

variable "node_type" {
  description = "The node type to use"
  type        = string
  default     = "cache.t3.medium"
}

variable "engine_version" {
  description = "The version of the Redis engine to use"
  type        = string
  default     = "7.0"
}

variable "parameter_group_family" {
  description = "The family of the parameter group"
  type        = string
  default     = "redis7"
}

variable "parameters" {
  description = "A map of parameters to apply to the parameter group"
  type        = map(string)
  default     = {}
}

variable "automatic_failover_enabled" {
  description = "Whether to enable automatic failover"
  type        = bool
  default     = true
}

variable "multi_az_enabled" {
  description = "Whether to enable multi-AZ deployment"
  type        = bool
  default     = true
}

variable "num_cache_clusters" {
  description = "The number of cache clusters in the replication group"
  type        = number
  default     = 2
}

variable "maintenance_window" {
  description = "The weekly time range during which system maintenance can occur"
  type        = string
  default     = "sun:04:00-sun:05:00"
}

variable "snapshot_window" {
  description = "The daily time range during which automated snapshots are created"
  type        = string
  default     = "03:00-04:00"
}

variable "snapshot_retention_limit" {
  description = "The number of days to retain snapshots"
  type        = number
  default     = 7
}

variable "apply_immediately" {
  description = "Whether to apply changes immediately or during the next maintenance window"
  type        = bool
  default     = false
}

variable "auto_minor_version_upgrade" {
  description = "Whether to automatically upgrade minor versions of the Redis engine"
  type        = bool
  default     = true
}

variable "notification_topic_arn" {
  description = "The ARN of the SNS topic to notify when changes occur"
  type        = string
  default     = ""
}

variable "tags" {
  description = "A map of tags to add to all resources"
  type        = map(string)
  default     = {}
}
