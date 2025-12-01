# Azure Container Apps Environment
# Requerido para Container Apps (PostgreSQL y aplicaciones)
resource "azurerm_container_app_environment" "main" {
  # Nombre debe ser: minúsculas, alfanuméricos o '-', 2-32 caracteres, sin '--'
  name                       = substr("${local.name_prefix}-env", 0, 32)
  location                   = var.location
  resource_group_name        = azurerm_resource_group.rg.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.main.id
}

# Log Analytics Workspace para Container Apps
resource "azurerm_log_analytics_workspace" "main" {
  # Nombre debe ser único globalmente, 4-63 caracteres, alfanuméricos y guiones
  name                = substr("${local.name_prefix}-logs", 0, 63)
  location            = var.location
  resource_group_name = azurerm_resource_group.rg.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
