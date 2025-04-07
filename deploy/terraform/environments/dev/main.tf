provider "aws" {
  region = var.region
}

terraform {
  backend "s3" {
    bucket         = "sparkfund-terraform-state"
    key            = "dev/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "sparkfund-terraform-locks"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
    }
  }

  required_version = ">= 1.0.0"
}

locals {
  environment = "dev"
  tags = {
    Environment = local.environment
    Project     = "SparkFund"
    ManagedBy   = "Terraform"
  }
}

# VPC
module "vpc" {
  source = "../../modules/vpc"

  region       = var.region
  environment  = local.environment
  vpc_cidr     = var.vpc_cidr
  az_count     = var.az_count
  cluster_name = var.cluster_name
  tags         = local.tags
}

# EKS
module "eks" {
  source = "../../modules/eks"

  region             = var.region
  environment        = local.environment
  cluster_name       = var.cluster_name
  kubernetes_version = var.kubernetes_version
  vpc_id             = module.vpc.vpc_id
  subnet_ids         = module.vpc.private_subnet_ids
  node_groups        = var.node_groups
  tags               = local.tags
}

# RDS
module "rds" {
  source = "../../modules/rds"

  region                 = var.region
  environment            = local.environment
  name                   = "postgres"
  vpc_id                 = module.vpc.vpc_id
  subnet_ids             = module.vpc.private_subnet_ids
  allowed_security_groups = [module.eks.node_security_group_id]
  db_name                = "sparkfund"
  instance_class         = var.rds_instance_class
  allocated_storage      = var.rds_allocated_storage
  max_allocated_storage  = var.rds_max_allocated_storage
  multi_az               = var.rds_multi_az
  deletion_protection    = var.rds_deletion_protection
  skip_final_snapshot    = var.rds_skip_final_snapshot
  tags                   = local.tags
}

# ElastiCache
module "elasticache" {
  source = "../../modules/elasticache"

  region                    = var.region
  environment               = local.environment
  name                      = "redis"
  vpc_id                    = module.vpc.vpc_id
  subnet_ids                = module.vpc.private_subnet_ids
  allowed_security_groups   = [module.eks.node_security_group_id]
  node_type                 = var.elasticache_node_type
  num_cache_clusters        = var.elasticache_num_cache_clusters
  automatic_failover_enabled = var.elasticache_automatic_failover_enabled
  multi_az_enabled          = var.elasticache_multi_az_enabled
  tags                      = local.tags
}

# Kubernetes provider configuration
data "aws_eks_cluster" "cluster" {
  name = module.eks.cluster_id
}

data "aws_eks_cluster_auth" "cluster" {
  name = module.eks.cluster_id
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.cluster.token
}

provider "helm" {
  kubernetes {
    host                   = data.aws_eks_cluster.cluster.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority[0].data)
    token                  = data.aws_eks_cluster_auth.cluster.token
  }
}

# Kubernetes Namespace
resource "kubernetes_namespace" "sparkfund" {
  metadata {
    name = "sparkfund"

    labels = {
      name        = "sparkfund"
      environment = local.environment
    }
  }
}

# Kubernetes Secrets for RDS
resource "kubernetes_secret" "rds" {
  metadata {
    name      = "rds-credentials"
    namespace = kubernetes_namespace.sparkfund.metadata[0].name
  }

  data = {
    host     = module.rds.instance_address
    port     = module.rds.instance_port
    username = module.rds.username
    password = jsondecode(data.aws_secretsmanager_secret_version.rds.secret_string)["password"]
    database = module.rds.db_name
  }
}

# Kubernetes Secrets for ElastiCache
resource "kubernetes_secret" "elasticache" {
  metadata {
    name      = "redis-credentials"
    namespace = kubernetes_namespace.sparkfund.metadata[0].name
  }

  data = {
    host       = module.elasticache.primary_endpoint_address
    port       = module.elasticache.port
    password   = jsondecode(data.aws_secretsmanager_secret_version.elasticache.secret_string)["auth_token"]
  }
}

# Data sources for secrets
data "aws_secretsmanager_secret_version" "rds" {
  secret_id = module.rds.secret_name
}

data "aws_secretsmanager_secret_version" "elasticache" {
  secret_id = module.elasticache.secret_name
}

# Helm release for AWS Load Balancer Controller
resource "helm_release" "aws_load_balancer_controller" {
  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  namespace  = "kube-system"
  version    = "1.6.0"

  set {
    name  = "clusterName"
    value = module.eks.cluster_id
  }

  set {
    name  = "serviceAccount.create"
    value = "true"
  }

  set {
    name  = "serviceAccount.name"
    value = "aws-load-balancer-controller"
  }

  set {
    name  = "serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
    value = module.eks.aws_load_balancer_controller_role_arn
  }
}

# Helm release for Cluster Autoscaler
resource "helm_release" "cluster_autoscaler" {
  name       = "cluster-autoscaler"
  repository = "https://kubernetes.github.io/autoscaler"
  chart      = "cluster-autoscaler"
  namespace  = "kube-system"
  version    = "9.29.0"

  set {
    name  = "autoDiscovery.clusterName"
    value = module.eks.cluster_id
  }

  set {
    name  = "awsRegion"
    value = var.region
  }

  set {
    name  = "rbac.serviceAccount.create"
    value = "true"
  }

  set {
    name  = "rbac.serviceAccount.name"
    value = "cluster-autoscaler"
  }

  set {
    name  = "rbac.serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
    value = module.eks.cluster_autoscaler_role_arn
  }
}

# Helm release for Cert Manager
resource "helm_release" "cert_manager" {
  name       = "cert-manager"
  repository = "https://charts.jetstack.io"
  chart      = "cert-manager"
  namespace  = "cert-manager"
  version    = "1.13.1"
  create_namespace = true

  set {
    name  = "installCRDs"
    value = "true"
  }
}
