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
  description = "The name of the RDS instance"
  type        = string
}

variable "vpc_id" {
  description = "The ID of the VPC"
  type        = string
}

variable "subnet_ids" {
  description = "The IDs of the subnets where the RDS instance will be deployed"
  type        = list(string)
}

variable "allowed_security_groups" {
  description = "The security groups allowed to access the RDS instance"
  type        = list(string)
  default     = []
}

variable "kms_key_arn" {
  description = "The ARN of the KMS key to use for encryption"
  type        = string
  default     = ""
}

variable "engine" {
  description = "The database engine to use"
  type        = string
  default     = "postgres"
}

variable "engine_version" {
  description = "The version of the database engine to use"
  type        = string
  default     = "14.7"
}

variable "instance_class" {
  description = "The instance class to use"
  type        = string
  default     = "db.t3.medium"
}

variable "allocated_storage" {
  description = "The amount of storage to allocate (in GB)"
  type        = number
  default     = 20
}

variable "max_allocated_storage" {
  description = "The maximum amount of storage to allocate (in GB)"
  type        = number
  default     = 100
}

variable "storage_type" {
  description = "The type of storage to use"
  type        = string
  default     = "gp3"
}

variable "db_name" {
  description = "The name of the database to create"
  type        = string
}

variable "username" {
  description = "The username for the database"
  type        = string
  default     = "postgres"
}

variable "parameter_group_family" {
  description = "The family of the parameter group"
  type        = string
  default     = "postgres14"
}

variable "parameters" {
  description = "A map of parameters to apply to the parameter group"
  type        = map(string)
  default     = {}
}

variable "skip_final_snapshot" {
  description = "Whether to skip the final snapshot when the RDS instance is deleted"
  type        = bool
  default     = false
}

variable "deletion_protection" {
  description = "Whether to enable deletion protection"
  type        = bool
  default     = true
}

variable "backup_retention_period" {
  description = "The number of days to retain backups"
  type        = number
  default     = 7
}

variable "backup_window" {
  description = "The daily time range during which automated backups are created"
  type        = string
  default     = "03:00-04:00"
}

variable "maintenance_window" {
  description = "The weekly time range during which system maintenance can occur"
  type        = string
  default     = "sun:04:00-sun:05:00"
}

variable "multi_az" {
  description = "Whether to deploy the RDS instance in multiple availability zones"
  type        = bool
  default     = true
}

variable "apply_immediately" {
  description = "Whether to apply changes immediately or during the next maintenance window"
  type        = bool
  default     = false
}

variable "auto_minor_version_upgrade" {
  description = "Whether to automatically upgrade minor versions of the database engine"
  type        = bool
  default     = true
}

variable "performance_insights_enabled" {
  description = "Whether to enable Performance Insights"
  type        = bool
  default     = true
}

variable "performance_insights_retention_period" {
  description = "The retention period for Performance Insights (in days)"
  type        = number
  default     = 7
}

variable "tags" {
  description = "A map of tags to add to all resources"
  type        = map(string)
  default     = {}
}
