output "function_app_id" {
  description = "ID del Function App"
  value       = azurerm_linux_function_app.main.id
}

output "function_app_name" {
  description = "Nombre del Function App"
  value       = azurerm_linux_function_app.main.name
}

output "function_app_default_hostname" {
  description = "Hostname por defecto del Function App"
  value       = azurerm_linux_function_app.main.default_hostname
}

output "application_insights_id" {
  description = "ID del Application Insights"
  value       = azurerm_application_insights.main.id
}

output "application_insights_instrumentation_key" {
  description = "Instrumentation key de Application Insights"
  value       = azurerm_application_insights.main.instrumentation_key
  sensitive   = true
}

output "storage_account_name" {
  description = "Nombre del Storage Account"
  value       = azurerm_storage_account.functions.name
}

output "service_plan_id" {
  description = "ID del Service Plan"
  value       = azurerm_service_plan.functions.id
}

output "managed_identity_principal_id" {
  description = "Principal ID de la Managed Identity"
  value       = azurerm_user_assigned_identity.function_app.principal_id
}

output "managed_identity_client_id" {
  description = "Client ID de la Managed Identity"
  value       = azurerm_user_assigned_identity.function_app.client_id
}
