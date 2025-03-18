# Money Pulse Infrastructure - Backend Configuration

terraform {
  backend "s3" {
    # These values must be provided via CLI or a separate terraform.tfvars file
    # bucket         = "money-pulse-tf-state"
    # key            = "terraform.tfstate"
    # region         = "us-west-2"
    # dynamodb_table = "money-pulse-tf-locks"
    # encrypt        = true
  }
}