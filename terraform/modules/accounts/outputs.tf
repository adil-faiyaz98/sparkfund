output "cluster_endpoint" {
  description = "Endpoint for the EKS cluster"
  value       = aws_eks_cluster.accounts.endpoint
}

output "cluster_ca_certificate" {
  description = "Base64 encoded certificate data required to communicate with the cluster"
  value       = aws_eks_cluster.accounts.certificate_authority[0].data
}

output "cluster_name" {
  description = "Name of the EKS cluster"
  value       = aws_eks_cluster.accounts.name
}

output "db_endpoint" {
  description = "Endpoint for the RDS instance"
  value       = aws_db_instance.accounts.endpoint
}

output "db_name" {
  description = "Name of the database"
  value       = aws_db_instance.accounts.db_name
}

output "db_username" {
  description = "Username for the database"
  value       = aws_db_instance.accounts.username
} 