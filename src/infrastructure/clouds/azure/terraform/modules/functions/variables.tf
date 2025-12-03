# Variables generales
variable "function_name" {
  description = "Nombre de la función (usado para nombres de recursos)"
  type        = string
}

variable "resource_group_name" {
  description = "Nombre del resource group"
  type        = string
}

variable "location" {
  description = "Ubicación de los recursos"
  type        = string
}

variable "name_prefix" {
  description = "Prefijo para nombres de recursos"
  type        = string
}

variable "project_name_clean" {
  description = "Nombre del proyecto limpio (sin guiones)"
  type        = string
}

variable "env_prefix" {
  description = "Prefijo del entorno (primeros 3 caracteres)"
  type        = string
}

variable "project_name" {
  description = "Nombre del proyecto"
  type        = string
}

variable "environment" {
  description = "Entorno del proyecto"
  type        = string
}

# Variables de aplicación
variable "app_port" {
  description = "Puerto de la aplicación"
  type        = string
  default     = "8080"
}

variable "app_version" {
  description = "Versión de la aplicación"
  type        = string
  default     = "0.0.1"
}

variable "app_description" {
  description = "Descripción de la aplicación"
  type        = string
  default     = "Gorm Go Skeleton Template"
}

variable "app_support_email" {
  description = "Email de soporte"
  type        = string
  default     = "support@goprojectskeleton.com"
}

variable "enable_log" {
  description = "Habilitar logging"
  type        = bool
  default     = true
}

variable "debug_log" {
  description = "Habilitar debug logging"
  type        = bool
  default     = true
}

variable "templates_path" {
  description = "Ruta de templates"
  type        = string
  default     = "src/application/shared/templates/"
}

# Variables de base de datos
variable "db_host" {
  description = "Host de la base de datos"
  type        = string
}

variable "db_port" {
  description = "Puerto de la base de datos"
  type        = string
  default     = "5432"
}

variable "db_user" {
  description = "Usuario de la base de datos"
  type        = string
}

variable "db_name" {
  description = "Nombre de la base de datos"
  type        = string
}

variable "db_ssl" {
  description = "Habilitar SSL para la base de datos"
  type        = bool
  default     = true
}

variable "db_password_secret_uri" {
  description = "URI del secreto de contraseña de DB en Key Vault"
  type        = string
}

# Variables de Redis
variable "redis_host" {
  description = "Host de Redis"
  type        = string
}

variable "redis_db" {
  description = "Base de datos de Redis"
  type        = number
  default     = 0
}

variable "redis_ttl" {
  description = "TTL de Redis"
  type        = number
  default     = 300
}

variable "redis_password_secret_uri" {
  description = "URI del secreto de contraseña de Redis en Key Vault"
  type        = string
}

# Variables de JWT
variable "jwt_secret_uri" {
  description = "URI del secreto JWT en Key Vault"
  type        = string
}

variable "jwt_issuer" {
  description = "Issuer de JWT"
  type        = string
  default     = "goprojectskeleton"
}

variable "jwt_access_ttl" {
  description = "TTL del token de acceso"
  type        = number
  default     = 3600
}

variable "jwt_refresh_ttl" {
  description = "TTL del token de refresh"
  type        = number
  default     = 86400
}

variable "jwt_clock_skew" {
  description = "Clock skew de JWT"
  type        = number
  default     = 60
}

# Variables de tokens y OTP
variable "one_time_token_ttl" {
  description = "TTL del token de un solo uso"
  type        = number
  default     = 15
}

variable "one_time_token_email_verify_ttl" {
  description = "TTL del token de verificación de email"
  type        = number
  default     = 60
}

variable "one_time_password_length" {
  description = "Longitud del OTP"
  type        = number
  default     = 6
}

variable "one_time_password_ttl" {
  description = "TTL del OTP"
  type        = number
  default     = 10
}

# Variables de frontend
variable "frontend_reset_password_url" {
  description = "URL de reset de contraseña del frontend"
  type        = string
  default     = "http://localhost:3000/reset-password"
}

variable "frontend_activate_account_url" {
  description = "URL de activación de cuenta del frontend"
  type        = string
  default     = "http://localhost:3000/activate-account"
}

# Variables de mail
variable "mail_host" {
  description = "Host del servidor de mail"
  type        = string
  default     = "localhost"
}

variable "mail_port" {
  description = "Puerto del servidor de mail"
  type        = number
  default     = 1025
}

variable "mail_from" {
  description = "Email remitente"
  type        = string
  default     = "noreply@example.com"
}

variable "mail_password_secret_uri" {
  description = "URI del secreto de contraseña de mail en Key Vault"
  type        = string
}

# Variables de Key Vault
variable "key_vault_id" {
  description = "ID del Key Vault"
  type        = string
}

variable "tenant_id" {
  description = "ID del tenant de Azure"
  type        = string
}

# Variables adicionales
variable "extra_app_settings" {
  description = "Configuraciones adicionales de app_settings"
  type        = map(string)
  default     = {}
}
