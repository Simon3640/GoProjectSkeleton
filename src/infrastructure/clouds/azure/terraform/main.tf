locals {
  # Limpiar nombre del proyecto (sin guiones ni guiones bajos)
  project_name_clean = replace(replace(var.project_name, "-", ""), "_", "")

  env_prefix = substr(var.environment, 0, 3)

  name_prefix = "${local.project_name_clean}-${local.env_prefix}"

  # Key Vault: 3-24 caracteres, solo alfanuméricos y guiones
  key_vault_name = "${local.name_prefix}-kv"

  # ACR: 5-50 caracteres, solo minúsculas y números
  acr_name = substr("${local.name_prefix}-acr", 0, 50)

  # Container App PostgreSQL: 2-32 caracteres, minúsculas, alfanuméricos o '-', sin '--'
  postgres_app_name = substr("${local.name_prefix}-pg", 0, 32)
}

resource "azurerm_resource_group" "rg" {
  name     = var.project_name
  location = var.location
}
