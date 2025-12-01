# Storage Account para Azure Functions
# Azure Functions requiere un Storage Account para almacenar el código y los datos de runtime
resource "azurerm_storage_account" "functions" {
  # Nombre debe ser único globalmente, 3-24 caracteres, solo minúsculas y números
  name                     = substr("${var.project_name}${local.env_prefix}funcst", 0, 24)
  resource_group_name      = azurerm_resource_group.rg.name
  location                 = var.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  # Habilitar HTTPS solo
  min_tls_version = "TLS1_2"

  # Blob soft delete para recuperación
  blob_properties {
    delete_retention_policy {
      days = 7
    }
  }
}

# App Service Plan para Azure Functions
# Usamos Consumption Plan (Y1) para serverless
resource "azurerm_service_plan" "functions" {
  # Nombre debe ser único, 1-60 caracteres
  name                = substr("${local.name_prefix}-func-plan", 0, 60)
  resource_group_name = azurerm_resource_group.rg.name
  location            = var.location
  os_type             = "Linux"
  sku_name            = "Y1" # Consumption Plan (serverless)
}

# Managed Identity para Function App
resource "azurerm_user_assigned_identity" "function_app" {
  name                = substr("${local.name_prefix}-func-id", 0, 128)
  location            = var.location
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_application_insights" "main" {
  name                = "${local.name_prefix}-appi"
  location            = var.location
  resource_group_name = azurerm_resource_group.rg.name
  application_type    = "web"

  # Habilitar retención de datos extendida (opcional, hasta 730 días)
  retention_in_days = 90

  # Habilitar sampling para capturar todos los logs (100% = sin sampling)
  sampling_percentage = 100.0

  # Ignorar cambios en workspace_id para evitar el error
  # "workspace_id can not be removed after set"
  # Si el Application Insights ya tiene un workspace_id configurado,
  # Terraform no intentará removerlo
  lifecycle {
    ignore_changes = [workspace_id]
  }
}

# Azure Function App
resource "azurerm_linux_function_app" "main" {
  # Nombre debe ser único globalmente, 2-60 caracteres, alfanuméricos y guiones
  name                = substr("${local.name_prefix}-func", 0, 60)
  resource_group_name = azurerm_resource_group.rg.name
  location            = var.location
  service_plan_id     = azurerm_service_plan.functions.id

  storage_account_name       = azurerm_storage_account.functions.name
  storage_account_access_key = azurerm_storage_account.functions.primary_access_key

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.function_app.id]
  }

  site_config {
    # Configuración para Custom Handler
    # No especificar linux_fx_version ni application_stack para Custom Handler

    # Configuración de CORS
    # Incluye el portal de Azure para poder probar funciones desde el portal
    cors {
      allowed_origins = [
        "*",                           # Permitir todos los orígenes (ajustar según necesidades de seguridad)
        "https://portal.azure.com",    # Portal de Azure para testing
        "https://functions.azure.com", # Portal de Azure Functions
      ]
      support_credentials = false
    }

    # Límites de tiempo y memoria
    app_scale_limit     = 200 # Máximo de instancias en Consumption Plan
    use_32_bit_worker   = false
    ftps_state          = "Disabled"
    http2_enabled       = true
    minimum_tls_version = "1.2"
  }

  # Variables de entorno de aplicación
  app_settings = {
    # Aplicación
    "APP_NAME"          = var.project_name
    "APP_ENV"           = var.environment
    "APP_PORT"          = var.app_port
    "APP_VERSION"       = var.app_version
    "APP_DESCRIPTION"   = var.app_description
    "APP_SUPPORT_EMAIL" = var.app_support_email

    # Logging
    "ENABLE_LOG"     = tostring(var.enable_log)
    "DEBUG_LOG"      = tostring(var.debug_log)
    "TEMPLATES_PATH" = var.templates_path

    # Application Insights - Configuración completa para capturar todos los logs
    "APPLICATION_INSIGHTS_CONNECTION_STRING" = azurerm_application_insights.main.connection_string
    "APPLICATION_INSIGHTS_INSTRUMENTATION_KEY" = azurerm_application_insights.main.instrumentation_key

    # Habilitar logging completo de Azure Functions
    "APPINSIGHTS_INSTRUMENTATIONKEY" = azurerm_application_insights.main.instrumentation_key

    # Configuración de logging para capturar todos los niveles
    "AzureWebJobsDashboard" = ""  # Deshabilitar dashboard antiguo, usar Application Insights
    "AzureFunctionsJobHost__Logging__LogLevel__Default" = "Information"
    "AzureFunctionsJobHost__Logging__LogLevel__Host" = "Information"
    "AzureFunctionsJobHost__Logging__LogLevel__Function" = "Information"
    "AzureFunctionsJobHost__Logging__LogLevel__Host.Results" = "Information"
    "AzureFunctionsJobHost__Logging__LogLevel__Microsoft" = "Warning"
    "AzureFunctionsJobHost__Logging__LogLevel__Worker" = "Information"

    # Habilitar telemetría detallada
    "APPINSIGHTS_ENABLE_AGENT" = "true"
    "APPINSIGHTS_PROACTIVE_SAMPLING_ENABLED" = "false"  # Deshabilitar sampling proactivo para capturar todo

    # Base de datos - PostgreSQL Flexible Server
    "DB_HOST"     = azurerm_postgresql_flexible_server.postgres.fqdn
    "DB_PORT"     = "5432"
    "DB_USER"     = var.db_user
    "DB_NAME"     = var.db_name
    "DB_SSL"      = "true"  # Flexible Server requiere SSL por defecto
    "DB_PASSWORD" = "@Microsoft.KeyVault(SecretUri=${azurerm_key_vault_secret.db_password.id})"

    # Redis - Azure Redis Cache
    "REDIS_HOST"     = "${azurerm_redis_cache.redis.hostname}:${azurerm_redis_cache.redis.ssl_port}"
    "REDIS_DB"       = tostring(var.redis_db)
    "REDIS_TTL"      = tostring(var.redis_ttl)
    "REDIS_PASSWORD" = "@Microsoft.KeyVault(SecretUri=${azurerm_key_vault_secret.redis_password.id})"

    # JWT
    "JWT_SECRET_KEY"  = "@Microsoft.KeyVault(SecretUri=${azurerm_key_vault_secret.jwt_secret.id})"
    "JWT_ISSUER"      = var.jwt_issuer
    "JWT_AUDIENCE"    = var.jwt_issuer
    "JWT_ACCESS_TTL"  = tostring(var.jwt_access_ttl)
    "JWT_REFRESH_TTL" = tostring(var.jwt_refresh_ttl)
    "JWT_CLOCK_SKEW"  = tostring(var.jwt_clock_skew)

    # Tokens y OTP
    "ONE_TIME_TOKEN_TTL"              = tostring(var.one_time_token_ttl)
    "ONE_TIME_TOKEN_EMAIL_VERIFY_TTL" = tostring(var.one_time_token_email_verify_ttl)
    "ONE_TIME_PASSWORD_LENGTH"        = tostring(var.one_time_password_length)
    "ONE_TIME_PASSWORD_TTL"           = tostring(var.one_time_password_ttl)

    # Frontend
    "FRONTEND_RESET_PASSWORD_URL"   = var.frontend_reset_password_url
    "FRONTEND_ACTIVATE_ACCOUNT_URL" = var.frontend_activate_account_url

    # Mail
    "MAIL_HOST"     = var.mail_host
    "MAIL_PORT"     = tostring(var.mail_port)
    "MAIL_FROM"     = var.mail_from
    "MAIL_PASSWORD" = "@Microsoft.KeyVault(SecretUri=${azurerm_key_vault_secret.sendgrid_key.id})"

    # Azure Functions específico - Configuración para despliegues desde CI/CD
    "FUNCTIONS_WORKER_RUNTIME"       = "custom"
    "FUNCTIONS_EXTENSION_VERSION"    = "~4"
    "WEBSITE_RUN_FROM_PACKAGE"       = "1"
    "SCM_DO_BUILD_DURING_DEPLOYMENT" = "true"
    "ENABLE_ORYX_BUILD"              = "true"

    # Managed Identity - Necesario para que Azure Functions resuelva referencias de Key Vault
    "AZURE_CLIENT_ID" = azurerm_user_assigned_identity.function_app.client_id
  }

  depends_on = [
    azurerm_storage_account.functions,
    azurerm_service_plan.functions,
    azurerm_key_vault_access_policy.function_app,
    azurerm_application_insights.main
  ]
}

# Key Vault Access Policy para Function App
resource "azurerm_key_vault_access_policy" "function_app" {
  key_vault_id = azurerm_key_vault.main.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.function_app.principal_id

  secret_permissions = [
    "Get",
    "List"
  ]

  depends_on = [
    azurerm_key_vault.main,
    azurerm_user_assigned_identity.function_app
  ]
}
