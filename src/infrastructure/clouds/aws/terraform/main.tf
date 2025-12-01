# Local values for naming and resource organization
# Equivalent to Azure's locals block and resource group naming
locals {
  # Clean project name (remove hyphens and underscores for resource naming)
  project_name_clean = replace(replace(var.project_name, "-", ""), "_", "")

  # Environment prefix (first 3 characters)
  env_prefix = substr(var.environment, 0, 3)

  # Name prefix for resources (equivalent to Azure's name_prefix)
  name_prefix = "${local.project_name_clean}-${local.env_prefix}"

  # Common tags applied to all resources (equivalent to Azure resource group tags)
  common_tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "Terraform"
  }
}
