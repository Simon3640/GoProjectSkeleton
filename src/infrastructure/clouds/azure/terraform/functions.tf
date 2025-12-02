# Leer el archivo JSON con las funciones
locals {
  functions_json = jsondecode(file("${path.module}/../functions.json"))

  # Convertir el JSON a un mapa para usar con for_each
  functions_map = {
    for fn in local.functions_json : fn.name => {
      name           = fn.name
      path           = fn.path
      handler        = fn.handler
      route          = fn.route
      method         = fn.method
      authLevel      = fn.authLevel
      needsAuth      = try(fn.needsAuth, false)
      needsQuery     = try(fn.needsQuery, false)
      hasPathParams  = try(fn.hasPathParams, false)
      pathParamName  = try(fn.pathParamName, "")
    }
  }
}

# Crear un Function App para cada función definida en el JSON
module "functions" {
  source = "./modules/functions"

  for_each = local.functions_map

  # Variables generales
  function_name      = replace(each.value.name, "-", "") # Limpiar guiones para nombres de recursos
  resource_group_name = azurerm_resource_group.rg.name
  location           = var.location
  name_prefix        = local.name_prefix
  project_name_clean = local.project_name_clean
  env_prefix         = local.env_prefix
  project_name       = var.project_name
  environment        = var.environment

  # Variables de aplicación
  app_port        = var.app_port
  app_version     = var.app_version
  app_description = var.app_description
  app_support_email = var.app_support_email
  enable_log      = var.enable_log
  debug_log       = var.debug_log
  templates_path  = var.templates_path

  # Variables de base de datos
  db_host                 = azurerm_postgresql_flexible_server.postgres.fqdn
  db_port                 = "5432"
  db_user                 = var.db_user
  db_name                 = var.db_name
  db_ssl                  = true
  db_password_secret_uri  = "@Microsoft.KeyVault(SecretUri=${azurerm_key_vault_secret.db_password.id})"

  # Variables de Redis
  redis_host              = "${azurerm_redis_cache.redis.hostname}:${azurerm_redis_cache.redis.ssl_port}"
  redis_db                = var.redis_db
  redis_ttl               = var.redis_ttl
  redis_password_secret_uri = "@Microsoft.KeyVault(SecretUri=${azurerm_key_vault_secret.redis_password.id})"

  # Variables de JWT
  jwt_secret_uri   = "@Microsoft.KeyVault(SecretUri=${azurerm_key_vault_secret.jwt_secret.id})"
  jwt_issuer       = var.jwt_issuer
  jwt_access_ttl   = var.jwt_access_ttl
  jwt_refresh_ttl  = var.jwt_refresh_ttl
  jwt_clock_skew   = var.jwt_clock_skew

  # Variables de tokens y OTP
  one_time_token_ttl              = var.one_time_token_ttl
  one_time_token_email_verify_ttl = var.one_time_token_email_verify_ttl
  one_time_password_length         = var.one_time_password_length
  one_time_password_ttl            = var.one_time_password_ttl

  # Variables de frontend
  frontend_reset_password_url   = var.frontend_reset_password_url
  frontend_activate_account_url = var.frontend_activate_account_url

  # Variables de mail
  mail_host                = var.mail_host
  mail_port                = var.mail_port
  mail_from                = var.mail_from
  mail_password_secret_uri  = "@Microsoft.KeyVault(SecretUri=${azurerm_key_vault_secret.sendgrid_key.id})"

  # Variables de Key Vault
  key_vault_id = azurerm_key_vault.main.id
  tenant_id    = data.azurerm_client_config.current.tenant_id

  depends_on = [
    azurerm_postgresql_flexible_server.postgres,
    azurerm_redis_cache.redis,
    azurerm_key_vault.main,
    azurerm_key_vault_secret.db_password,
    azurerm_key_vault_secret.redis_password,
    azurerm_key_vault_secret.jwt_secret,
    azurerm_key_vault_secret.sendgrid_key
  ]
}
