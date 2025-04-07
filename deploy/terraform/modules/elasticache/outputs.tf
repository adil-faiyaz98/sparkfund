output "replication_group_id" {
  description = "The ID of the ElastiCache replication group"
  value       = aws_elasticache_replication_group.main.id
}

output "primary_endpoint_address" {
  description = "The address of the primary endpoint"
  value       = aws_elasticache_replication_group.main.primary_endpoint_address
}

output "reader_endpoint_address" {
  description = "The address of the reader endpoint"
  value       = aws_elasticache_replication_group.main.reader_endpoint_address
}

output "port" {
  description = "The port of the ElastiCache cluster"
  value       = aws_elasticache_replication_group.main.port
}

output "security_group_id" {
  description = "The ID of the security group for the ElastiCache cluster"
  value       = aws_security_group.elasticache.id
}

output "secret_arn" {
  description = "The ARN of the secret containing the ElastiCache credentials"
  value       = aws_secretsmanager_secret.elasticache.arn
}

output "secret_name" {
  description = "The name of the secret containing the ElastiCache credentials"
  value       = aws_secretsmanager_secret.elasticache.name
}
