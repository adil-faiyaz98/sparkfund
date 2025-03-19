terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

# EKS cluster
resource "aws_eks_cluster" "accounts" {
  name     = "accounts-cluster"
  role_arn = aws_iam_role.eks_cluster.arn
  version  = "1.27"

  vpc_config {
    subnet_ids = var.subnet_ids
  }

  depends_on = [aws_iam_role_policy_attachment.eks_cluster_policy]
}

# EKS node group
resource "aws_eks_node_group" "accounts" {
  cluster_name    = aws_eks_cluster.accounts.name
  node_group_name = "accounts-node-group"
  node_role_arn   = aws_iam_role.eks_node.arn
  subnet_ids      = var.subnet_ids

  scaling_config {
    desired_size = 2
    max_size     = 4
    min_size     = 1
  }

  instance_types = ["t3.medium"]

  depends_on = [
    aws_iam_role_policy_attachment.eks_worker_node_policy,
    aws_iam_role_policy_attachment.eks_cni_policy,
    aws_iam_role_policy_attachment.ecr_read_only
  ]
}

# IAM roles and policies
resource "aws_iam_role" "eks_cluster" {
  name = "eks-cluster-role"
}

resource "aws_iam_role" "eks_node" {
  name = "eks-node-role"
}

resource "aws_iam_role_policy_attachment" "eks_cluster_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.eks_cluster.name
}

resource "aws_iam_role_policy_attachment" "eks_worker_node_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.eks_node.name
}

resource "aws_iam_role_policy_attachment" "eks_cni_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.eks_node.name
}

resource "aws_iam_role_policy_attachment" "ecr_read_only" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.eks_node.name
}

# RDS instance for PostgreSQL
resource "aws_db_instance" "accounts" {
  identifier           = "accounts-db"
  engine              = "postgres"
  engine_version      = "14"
  instance_class      = "db.t3.micro"
  allocated_storage   = 20
  storage_type        = "gp2"
  db_name             = "money_pulse"
  username           = var.db_username
  password           = var.db_password
  skip_final_snapshot = true

  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.accounts.name
}

# RDS security group
resource "aws_security_group" "rds" {
  name        = "accounts-rds-sg"
  description = "Security group for accounts RDS instance"
  vpc_id      = var.vpc_id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.eks.id]
  }
}

# EKS security group
resource "aws_security_group" "eks" {
  name        = "accounts-eks-sg"
  description = "Security group for accounts EKS cluster"
  vpc_id      = var.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# DB subnet group
resource "aws_db_subnet_group" "accounts" {
  name       = "accounts-db-subnet-group"
  subnet_ids = var.subnet_ids
}

# Outputs
output "cluster_endpoint" {
  value = aws_eks_cluster.accounts.endpoint
}

output "cluster_ca_certificate" {
  value = aws_eks_cluster.accounts.certificate_authority[0].data
}

output "db_endpoint" {
  value = aws_db_instance.accounts.endpoint
} 