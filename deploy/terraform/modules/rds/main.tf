provider "aws" {
  region = var.region
}

locals {
  identifier = "${var.environment}-${var.name}"
}

# KMS Key for RDS
resource "aws_kms_key" "rds" {
  count                   = var.kms_key_arn == "" ? 1 : 0
  description             = "KMS key for RDS instance ${local.identifier}"
  deletion_window_in_days = 7
  enable_key_rotation     = true

  tags = var.tags
}

resource "aws_kms_alias" "rds" {
  count         = var.kms_key_arn == "" ? 1 : 0
  name          = "alias/${local.identifier}-key"
  target_key_id = aws_kms_key.rds[0].key_id
}

# RDS Subnet Group
resource "aws_db_subnet_group" "main" {
  name       = "${local.identifier}-subnet-group"
  subnet_ids = var.subnet_ids

  tags = merge(
    var.tags,
    {
      Name = "${local.identifier}-subnet-group"
    }
  )
}

# RDS Parameter Group
resource "aws_db_parameter_group" "main" {
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

# RDS Security Group
resource "aws_security_group" "rds" {
  name        = "${local.identifier}-sg"
  description = "Security group for RDS instance ${local.identifier}"
  vpc_id      = var.vpc_id

  ingress {
    from_port       = 5432
    to_port         = 5432
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

# Random password for RDS
resource "random_password" "password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

# RDS Instance
resource "aws_db_instance" "main" {
  identifier             = local.identifier
  engine                 = var.engine
  engine_version         = var.engine_version
  instance_class         = var.instance_class
  allocated_storage      = var.allocated_storage
  max_allocated_storage  = var.max_allocated_storage
  storage_type           = var.storage_type
  storage_encrypted      = true
  kms_key_id             = var.kms_key_arn != "" ? var.kms_key_arn : aws_kms_key.rds[0].arn
  db_name                = var.db_name
  username               = var.username
  password               = random_password.password.result
  port                   = 5432
  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name
  parameter_group_name   = aws_db_parameter_group.main.name
  publicly_accessible    = false
  skip_final_snapshot    = var.skip_final_snapshot
  deletion_protection    = var.deletion_protection
  backup_retention_period = var.backup_retention_period
  backup_window          = var.backup_window
  maintenance_window     = var.maintenance_window
  multi_az               = var.multi_az
  apply_immediately      = var.apply_immediately
  auto_minor_version_upgrade = var.auto_minor_version_upgrade
  copy_tags_to_snapshot  = true
  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]
  performance_insights_enabled = var.performance_insights_enabled
  performance_insights_retention_period = var.performance_insights_enabled ? var.performance_insights_retention_period : null
  performance_insights_kms_key_id = var.performance_insights_enabled ? (var.kms_key_arn != "" ? var.kms_key_arn : aws_kms_key.rds[0].arn) : null

  tags = merge(
    var.tags,
    {
      Name = local.identifier
    }
  )
}

# Store RDS credentials in AWS Secrets Manager
resource "aws_secretsmanager_secret" "rds" {
  name        = "${local.identifier}-credentials"
  description = "RDS credentials for ${local.identifier}"
  kms_key_id  = var.kms_key_arn != "" ? var.kms_key_arn : aws_kms_key.rds[0].arn

  tags = var.tags
}

resource "aws_secretsmanager_secret_version" "rds" {
  secret_id = aws_secretsmanager_secret.rds.id
  secret_string = jsonencode({
    username = aws_db_instance.main.username
    password = random_password.password.result
    engine   = aws_db_instance.main.engine
    host     = aws_db_instance.main.address
    port     = aws_db_instance.main.port
    dbname   = aws_db_instance.main.db_name
  })
}
