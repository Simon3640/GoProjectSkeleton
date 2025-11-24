# Outputs del Resource Group
output "resource_group_name" {
  description = "Nombre del Resource Group"
  value       = azurerm_resource_group.rg.name
}

output "resource_group_location" {
  description = "Ubicación del Resource Group"
  value       = azurerm_resource_group.rg.location
}

# Outputs de PostgreSQL
output "postgres_fqdn" {
  description = "FQDN del PostgreSQL Flexible Server"
  value       = azurerm_postgresql_flexible_server.postgres.fqdn
}

output "postgres_connection_string" {
  description = "Connection string para PostgreSQL (sin contraseña por seguridad)"
  value       = "postgresql://${var.db_user}@${azurerm_postgresql_flexible_server.postgres.fqdn}:5432/${var.db_name}?sslmode=require"
  sensitive   = false
}

# Outputs de Redis
output "redis_hostname" {
  description = "Hostname de Redis Cache"
  value       = azurerm_redis_cache.redis.hostname
}

output "redis_port" {
  description = "Puerto de Redis (SSL para Cache)"
  value       = azurerm_redis_cache.redis.ssl_port
}

output "redis_type" {
  description = "Tipo de Redis utilizado"
  value       = "Azure Redis Cache"
}

output "redis_primary_key" {
  description = "Primary access key de Redis (almacenado en Key Vault)"
  value       = "Ver en Key Vault: redis-password"
  sensitive   = true
}

# Outputs de Key Vault
output "key_vault_name" {
  description = "Nombre del Key Vault"
  value       = azurerm_key_vault.main.name
}

output "key_vault_uri" {
  description = "URI del Key Vault"
  value       = azurerm_key_vault.main.vault_uri
}

# Outputs de SendGrid
output "sendgrid_username" {
  description = "Username de SendGrid"
  value       = var.mail_from
}

output "sendgrid_key_info" {
  description = "API Key de SendGrid (almacenado en Key Vault)"
  value       = "Ver en Key Vault: sendgrid-api-key"
  sensitive   = true
}

# Outputs de Azure Functions (para CI/CD)
output "function_app_name" {
  description = "Nombre del Function App para despliegues desde CI/CD"
  value       = azurerm_linux_function_app.main.name
}

output "function_app_url" {
  description = "URL base del Function App"
  value       = "https://${azurerm_linux_function_app.main.default_hostname}"
}
