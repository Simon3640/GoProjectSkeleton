# Azure Container Registry para almacenar imágenes Docker (OPCIONAL)
# Solo necesario si usas Container Apps con imágenes personalizadas
# Para Azure Functions serverless NO es necesario
resource "azurerm_container_registry" "main" {
  count               = var.create_container_registry ? 1 : 0
  name                = lower(local.acr_name)  # ACR requiere minúsculas
  resource_group_name = azurerm_resource_group.rg.name
  location            = var.location
  sku                 = "Basic"  # Para producción usar Standard o Premium
  admin_enabled       = true     # Habilitar admin user para simplificar autenticación
}

# Outputs del Container Registry (solo si se crea)
output "container_registry_name" {
  description = "Nombre del Container Registry"
  value       = var.create_container_registry && length(azurerm_container_registry.main) > 0 ? azurerm_container_registry.main[0].name : null
}

output "container_registry_login_server" {
  description = "Login server del Container Registry"
  value       = var.create_container_registry && length(azurerm_container_registry.main) > 0 ? azurerm_container_registry.main[0].login_server : null
}

output "container_registry_admin_username" {
  description = "Admin username del Container Registry"
  value       = var.create_container_registry && length(azurerm_container_registry.main) > 0 ? azurerm_container_registry.main[0].admin_username : null
  sensitive   = true
}

output "container_registry_admin_password" {
  description = "Admin password del Container Registry"
  value       = var.create_container_registry && length(azurerm_container_registry.main) > 0 ? azurerm_container_registry.main[0].admin_password : null
  sensitive   = true
}
