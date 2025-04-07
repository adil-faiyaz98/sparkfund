output "instance_id" {
  description = "The ID of the RDS instance"
  value       = aws_db_instance.main.id
}

output "instance_address" {
  description = "The address of the RDS instance"
  value       = aws_db_instance.main.address
}

output "instance_endpoint" {
  description = "The endpoint of the RDS instance"
  value       = aws_db_instance.main.endpoint
}

output "instance_port" {
  description = "The port of the RDS instance"
  value       = aws_db_instance.main.port
}

output "db_name" {
  description = "The name of the database"
  value       = aws_db_instance.main.db_name
}

output "username" {
  description = "The username for the database"
  value       = aws_db_instance.main.username
}

output "security_group_id" {
  description = "The ID of the security group for the RDS instance"
  value       = aws_security_group.rds.id
}

output "secret_arn" {
  description = "The ARN of the secret containing the RDS credentials"
  value       = aws_secretsmanager_secret.rds.arn
}

output "secret_name" {
  description = "The name of the secret containing the RDS credentials"
  value       = aws_secretsmanager_secret.rds.name
}
