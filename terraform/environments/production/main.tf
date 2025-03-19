terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket         = "money-pulse-terraform-state"
    key            = "production/terraform.tfstate"
    region         = "us-west-2"
    encrypt        = true
  }
}

provider "aws" {
  region = "us-west-2"
}

module "accounts" {
  source = "../../modules/accounts"

  vpc_id     = "vpc-87654321"
  subnet_ids = ["subnet-8765", "subnet-4321"]

  db_username = "postgres"
  db_password = var.db_password
} 