# Terraform backend configuration
# Equivalent to Azure Storage Account backend, using S3 instead
terraform {
  backend "s3" {
    # These values should be provided via backend configuration or environment variables
    bucket         = "goprojectskeleton-terraform-state"
    key            = "dev/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "terraform-locks" # For state locking (equivalent to Azure blob lease)
    encrypt        = true
  }
}
