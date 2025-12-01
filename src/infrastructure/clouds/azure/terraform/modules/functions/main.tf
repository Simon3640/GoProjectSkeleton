# Storage Account para Azure Functions
resource "azurerm_storage_account" "functions" {
  name                     = substr("${var.project_name_clean}${var.env_prefix}${var.function_name}st", 0, 24)
  resource_group_name      = var.resource_group_name
  location                 = var.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  min_tls_version         = "TLS1_2"

  blob_properties {
    delete_retention_policy {
      days = 7
    }
  }
}

# App Service Plan para Azure Functions
resource "azurerm_service_plan" "functions" {
  name                = substr("${var.name_prefix}-${var.function_name}-plan", 0, 60)
  resource_group_name = var.resource_group_name
  location            = var.location
  os_type             = "Linux"
  sku_name            = "Y1" # Consumption Plan (serverless)
}

# Managed Identity para Function App
resource "azurerm_user_assigned_identity" "function_app" {
  name                = substr("${var.name_prefix}-${var.function_name}-id", 0, 128)
  location            = var.location
  resource_group_name = var.resource_group_name
}

# Application Insights para esta función
resource "azurerm_application_insights" "main" {
  name                = "${var.name_prefix}-${var.function_name}-appi"
  location            = var.location
  resource_group_name = var.resource_group_name
  application_type    = "web"
  retention_in_days   = 90
  sampling_percentage = 100.0

  lifecycle {
    ignore_changes = [workspace_id]
  }
}

# Azure Function App
resource "azurerm_linux_function_app" "main" {
  name                = substr("${var.name_prefix}-${var.function_name}-func", 0, 60)
  resource_group_name = var.resource_group_name
  location            = var.location
  service_plan_id     = azurerm_service_plan.functions.id

  storage_account_name       = azurerm_storage_account.functions.name
  storage_account_access_key = azurerm_storage_account.functions.primary_access_key

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.function_app.id]
  }

  site_config {
    cors {
      allowed_origins = [
        "*",
        "https://portal.azure.com",
        "https://functions.azure.com",
      ]
      support_credentials = false
    }

    app_scale_limit     = 200
    use_32_bit_worker   = false
    ftps_state          = "Disabled"
    http2_enabled       = true
    minimum_tls_version = "1.2"
  }

  app_settings = merge(
    {
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

      # Application Insights
      "APPLICATION_INSIGHTS_CONNECTION_STRING" = azurerm_application_insights.main.connection_string
      "APPLICATION_INSIGHTS_INSTRUMENTATION_KEY" = azurerm_application_insights.main.instrumentation_key
      "APPINSIGHTS_INSTRUMENTATIONKEY"         = azurerm_application_insights.main.instrumentation_key
      "AzureWebJobsDashboard"                  = ""
      "AzureFunctionsJobHost__Logging__LogLevel__Default" = "Information"
      "AzureFunctionsJobHost__Logging__LogLevel__Host"     = "Information"
      "AzureFunctionsJobHost__Logging__LogLevel__Function" = "Information"
      "AzureFunctionsJobHost__Logging__LogLevel__Host.Results" = "Information"
      "AzureFunctionsJobHost__Logging__LogLevel__Microsoft"    = "Warning"
      "AzureFunctionsJobHost__Logging__LogLevel__Worker"        = "Information"
      "APPINSIGHTS_ENABLE_AGENT"                = "true"
      "APPINSIGHTS_PROACTIVE_SAMPLING_ENABLED"  = "false"

      # Base de datos
      "DB_HOST"     = var.db_host
      "DB_PORT"     = var.db_port
      "DB_USER"     = var.db_user
      "DB_NAME"     = var.db_name
      "DB_SSL"      = tostring(var.db_ssl)
      "DB_PASSWORD" = var.db_password_secret_uri

      # Redis
      "REDIS_HOST"     = var.redis_host
      "REDIS_DB"       = tostring(var.redis_db)
      "REDIS_TTL"      = tostring(var.redis_ttl)
      "REDIS_PASSWORD" = var.redis_password_secret_uri

      # JWT
      "JWT_SECRET_KEY"  = var.jwt_secret_uri
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
      "MAIL_PASSWORD" = var.mail_password_secret_uri

      # Azure Functions específico
      "FUNCTIONS_WORKER_RUNTIME"       = "custom"
      "FUNCTIONS_EXTENSION_VERSION"    = "~4"
      "WEBSITE_RUN_FROM_PACKAGE"       = "1"
      "SCM_DO_BUILD_DURING_DEPLOYMENT" = "true"
      "ENABLE_ORYX_BUILD"              = "true"

      # Managed Identity
      "AZURE_CLIENT_ID" = azurerm_user_assigned_identity.function_app.client_id
    },
    var.extra_app_settings
  )

  depends_on = [
    azurerm_storage_account.functions,
    azurerm_service_plan.functions,
    azurerm_application_insights.main
  ]
}

# Key Vault Access Policy para Function App
resource "azurerm_key_vault_access_policy" "function_app" {
  key_vault_id = var.key_vault_id
  tenant_id    = var.tenant_id
  object_id    = azurerm_user_assigned_identity.function_app.principal_id

  secret_permissions = [
    "Get",
    "List"
  ]

  depends_on = [
    azurerm_user_assigned_identity.function_app
  ]
}
