variable "namespace" {
  description = "Kubernetes namespace for the accounts service"
  type        = string
  default     = "finance"
}

variable "replicas" {
  description = "Number of replicas for the accounts service"
  type        = number
  default     = 2
}

variable "ecr_repository_url" {
  description = "URL of the ECR repository containing the accounts service image"
  type        = string
}

variable "image_tag" {
  description = "Tag of the accounts service image to deploy"
  type        = string
}

variable "db_host" {
  description = "Hostname of the PostgreSQL database"
  type        = string
}

variable "db_port" {
  description = "Port of the PostgreSQL database"
  type        = string
  default     = "5432"
}

variable "db_user" {
  description = "Username for the PostgreSQL database"
  type        = string
}

variable "db_password" {
  description = "Password for the PostgreSQL database"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "Name of the PostgreSQL database"
  type        = string
} 