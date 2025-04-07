provider "aws" {
  region = var.region
}

locals {
  identifier = "${var.environment}-${var.name}"
}

# KMS Key for ElastiCache
resource "aws_kms_key" "elasticache" {
  count                   = var.kms_key_arn == "" ? 1 : 0
  description             = "KMS key for ElastiCache cluster ${local.identifier}"
  deletion_window_in_days = 7
  enable_key_rotation     = true

  tags = var.tags
}

resource "aws_kms_alias" "elasticache" {
  count         = var.kms_key_arn == "" ? 1 : 0
  name          = "alias/${local.identifier}-key"
  target_key_id = aws_kms_key.elasticache[0].key_id
}

# ElastiCache Subnet Group
resource "aws_elasticache_subnet_group" "main" {
  name       = "${local.identifier}-subnet-group"
  subnet_ids = var.subnet_ids

  tags = merge(
    var.tags,
    {
      Name = "${local.identifier}-subnet-group"
    }
  )
}

# ElastiCache Parameter Group
resource "aws_elasticache_parameter_group" "main" {
  name   = "${local.identifier}-parameter-group"
  family = var.parameter_group_family

  dynamic "parameter" {
    for_each = var.parameters
    content {
      name  = parameter.key
      value = parameter.value
    }
  }

  tags = merge(
    var.tags,
    {
      Name = "${local.identifier}-parameter-group"
    }
  )
}

# ElastiCache Security Group
resource "aws_security_group" "elasticache" {
  name        = "${local.identifier}-sg"
  description = "Security group for ElastiCache cluster ${local.identifier}"
  vpc_id      = var.vpc_id

  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = var.allowed_security_groups
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    var.tags,
    {
      Name = "${local.identifier}-sg"
    }
  )
}

# Random password for ElastiCache
resource "random_password" "auth_token" {
  length           = 32
  special          = false
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

# ElastiCache Replication Group
resource "aws_elasticache_replication_group" "main" {
  replication_group_id          = local.identifier
  description                   = "ElastiCache cluster for ${local.identifier}"
  node_type                     = var.node_type
  port                          = 6379
  parameter_group_name          = aws_elasticache_parameter_group.main.name
  subnet_group_name             = aws_elasticache_subnet_group.main.name
  security_group_ids            = [aws_security_group.elasticache.id]
  automatic_failover_enabled    = var.automatic_failover_enabled
  multi_az_enabled              = var.multi_az_enabled
  num_cache_clusters            = var.num_cache_clusters
  at_rest_encryption_enabled    = true
  transit_encryption_enabled    = true
  auth_token                    = random_password.auth_token.result
  kms_key_id                    = var.kms_key_arn != "" ? var.kms_key_arn : aws_kms_key.elasticache[0].arn
  engine_version                = var.engine_version
  maintenance_window            = var.maintenance_window
  snapshot_window               = var.snapshot_window
  snapshot_retention_limit      = var.snapshot_retention_limit
  apply_immediately             = var.apply_immediately
  auto_minor_version_upgrade    = var.auto_minor_version_upgrade
  notification_topic_arn        = var.notification_topic_arn

  tags = merge(
    var.tags,
    {
      Name = local.identifier
    }
  )
}

# Store ElastiCache credentials in AWS Secrets Manager
resource "aws_secretsmanager_secret" "elasticache" {
  name        = "${local.identifier}-credentials"
  description = "ElastiCache credentials for ${local.identifier}"
  kms_key_id  = var.kms_key_arn != "" ? var.kms_key_arn : aws_kms_key.elasticache[0].arn

  tags = var.tags
}

resource "aws_secretsmanager_secret_version" "elasticache" {
  secret_id = aws_secretsmanager_secret.elasticache.id
  secret_string = jsonencode({
    auth_token = random_password.auth_token.result
    host       = aws_elasticache_replication_group.main.primary_endpoint_address
    port       = aws_elasticache_replication_group.main.port
  })
}
