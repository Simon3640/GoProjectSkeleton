# Key Vault para almacenar secretos
# Nota: El nombre debe tener entre 3-24 caracteres, solo alfanuméricos y guiones
resource "azurerm_key_vault" "main" {
  name                = local.key_vault_name
  location            = var.location
  resource_group_name = azurerm_resource_group.rg.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  soft_delete_retention_days = 7
  purge_protection_enabled    = false

  # Permitir acceso desde App Service
  network_acls {
    default_action = "Allow"
    bypass         = "AzureServices"
  }
}

# Data source para obtener configuración del cliente actual
data "azurerm_client_config" "current" {}

# Access Policy para el usuario/entidad de servicio actual (para que Terraform pueda crear secretos)
resource "azurerm_key_vault_access_policy" "current_user" {
  key_vault_id = azurerm_key_vault.main.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  secret_permissions = [
    "Get",
    "List",
    "Set",
    "Delete",
    "Recover",
    "Backup",
    "Restore"
  ]

  depends_on = [azurerm_key_vault.main]
}

# Key Vault Secrets
resource "azurerm_key_vault_secret" "db_password" {
  name         = "db-password"
  value        = var.db_password
  key_vault_id = azurerm_key_vault.main.id

  depends_on = [
    azurerm_key_vault_access_policy.current_user
  ]
}

# Redis Password - Usa la primary_access_key de Azure Redis Cache
resource "azurerm_key_vault_secret" "redis_password" {
  name         = "redis-password"
  # Usar la primary_access_key de Azure Redis Cache
  value        = azurerm_redis_cache.redis.primary_access_key
  key_vault_id = azurerm_key_vault.main.id

  depends_on = [
    azurerm_key_vault_access_policy.current_user,
    azurerm_redis_cache.redis
  ]
}

resource "azurerm_key_vault_secret" "jwt_secret" {
  name         = "jwt-secret"
  value        = var.jwt_secret_key
  key_vault_id = azurerm_key_vault.main.id

  depends_on = [
    azurerm_key_vault_access_policy.current_user
  ]
}

resource "azurerm_key_vault_secret" "sendgrid_key" {
  name         = "sendgrid-api-key"
  value        = var.mail_password
  key_vault_id = azurerm_key_vault.main.id

  depends_on = [
    azurerm_key_vault_access_policy.current_user
  ]
}
