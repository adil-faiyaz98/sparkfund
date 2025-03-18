# Money Pulse Infrastructure - Main Configuration

provider "aws" {
  region = var.region
}

locals {
  project_name = var.project_name
  environment  = var.environment
  tags = {
    Project     = local.project_name
    Environment = local.environment
    ManagedBy   = "Terraform"
  }
}

# VPC Module
module "vpc" {
  source = "./modules/vpc"

  project_name      = local.project_name
  environment       = local.environment
  vpc_cidr          = var.vpc_cidr
  azs               = var.availability_zones
  private_subnets   = var.private_subnets
  public_subnets    = var.public_subnets
  database_subnets  = var.database_subnets
  tags              = local.tags
}

# EKS Module
module "eks" {
  source = "./modules/eks"

  project_name        = local.project_name
  environment         = local.environment
  vpc_id              = module.vpc.vpc_id
  private_subnet_ids  = module.vpc.private_subnet_ids
  public_subnet_ids   = module.vpc.public_subnet_ids
  eks_version         = var.eks_version
  node_instance_types = var.node_instance_types
  node_desired_size   = var.node_desired_size
  node_min_size       = var.node_min_size
  node_max_size       = var.node_max_size
  tags                = local.tags
}

# RDS PostgreSQL Module
module "rds" {
  source = "./modules/rds"

  project_name        = local.project_name
  environment         = local.environment
  vpc_id              = module.vpc.vpc_id
  database_subnet_ids = module.vpc.database_subnet_ids
  allowed_cidr_blocks = module.vpc.private_subnet_cidr_blocks
  db_name             = var.db_name
  db_username         = var.db_username
  db_password         = var.db_password
  db_instance_class   = var.db_instance_class
  tags                = local.tags
}

# Monitoring Module (CloudWatch)
module "monitoring" {
  source = "./modules/monitoring"

  project_name  = local.project_name
  environment   = local.environment
  eks_cluster_name = module.eks.cluster_name
  tags          = local.tags
}