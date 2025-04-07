variable "region" {
  description = "The AWS region to deploy to"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "The environment (dev, staging, prod)"
  type        = string
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

variable "vpc_id" {
  description = "The ID of the VPC"
  type        = string
}

variable "subnet_ids" {
  description = "The IDs of the subnets where the EKS cluster will be deployed"
  type        = list(string)
}

variable "kms_key_arn" {
  description = "The ARN of the KMS key to use for encryption"
  type        = string
  default     = ""
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

variable "tags" {
  description = "A map of tags to add to all resources"
  type        = map(string)
  default     = {}
}
