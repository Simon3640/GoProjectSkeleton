variable "project_name" {
  description = "Nombre del proyecto"
  type        = string
  default     = "go-project-skeleton"
}

variable "location" {
  description = "Ubicación del proyecto (eastus y westus2 pueden estar restringidas, usar centralus, westeurope, etc.)"
  type        = string
  default     = "centralus"
}

variable "environment" {
  description = "Entorno del proyecto"
  type        = string
  default     = "development"
}

variable "create_container_registry" {
  description = "Crear Azure Container Registry (solo necesario para Container Apps con imágenes personalizadas, NO necesario para Functions serverless)"
  type        = bool
  default     = false  # Por defecto false ya que Functions no necesita ACR
}

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
  description = "Email de soporte de la aplicación"
  type        = string
  default     = "support@gormgoskeleton.com"
}

variable "enable_log" {
  description = "Habilitar logging"
  type        = bool
  default     = true
}

variable "debug_log" {
  description = "Habilitar logging de depuración"
  type        = bool
  default     = true
}

variable "templates_path" {
  description = "Ruta de los templates"
  type        = string
  default     = "src/application/shared/templates/"
}

variable "db_host" {
  description = "Host de la base de datos"
  type        = string
  default     = "localhost"
}

variable "db_port" {
  description = "Puerto de la base de datos"
  type        = string
  default     = "5432"
}

variable "db_user" {
  description = "Usuario de la base de datos"
  type        = string
  default     = "gormgoskeleton"
}

variable "db_password" {
  description = "Contraseña de la base de datos. DEBE ser una contraseña segura. Generar con: openssl rand -base64 32"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "Nombre de la base de datos"
  type        = string
  default     = "gormgoskeleton"
}

variable "db_ssl" {
  description = "Habilitar SSL para la base de datos"
  type        = bool
  default     = false
}

variable "redis_host" {
  description = "Host de Redis"
  type        = string
  default     = "localhost:6379"
}

variable "redis_password" {
  description = "Contraseña de Redis. No se usa directamente, se obtiene de Azure Redis Cache primary_access_key cuando use_redis_cache=true"
  type        = string
  sensitive   = true
  default     = ""  # No se usa, se obtiene automáticamente de Azure Redis Cache
}

variable "jwt_secret_key" {
  description = "Clave secreta de JWT. DEBE ser una cadena aleatoria segura (mínimo 32 caracteres recomendado). Generar con: openssl rand -hex 32"
  type        = string
  sensitive   = true
}

variable "jwt_issuer" {
  description = "Issuer de JWT"
  type        = string
  default     = "gormgoskeleton"
}

variable "redis_db" {
  description = "Base de datos de Redis"
  type        = number
  default     = 0
}

variable "redis_ttl" {
  description = "Tiempo de vida de Redis"
  type        = number
  default     = 300
}


variable "jwt_access_ttl" {
  description = "Tiempo de vida del token de acceso"
  type        = number
  default     = 3600
}

variable "jwt_refresh_ttl" {
  description = "Tiempo de vida del token de refresco"
  type        = number
  default     = 86400
}

variable "jwt_clock_skew" {
  description = "Desviación de reloj de JWT"
  type        = number
  default     = 60
}

variable "one_time_token_ttl" {
  description = "Tiempo de vida del token de uno tiempo"
  type        = number
  default     = 15
}

variable "one_time_token_email_verify_ttl" {
  description = "Tiempo de vida del token de verificación de email"
  type        = number
  default     = 60
}

variable "one_time_password_length" {
  description = "Longitud de la contraseña de uno tiempo"
  type        = number
  default     = 6
}

variable "one_time_password_ttl" {
  description = "Tiempo de vida de la contraseña de uno tiempo"
  type        = number
  default     = 10
}

variable "frontend_reset_password_url" {
  description = "URL del frontend para reset de contraseña"
  type        = string
  default     = "http://localhost:3000/reset-password"
}

variable "frontend_activate_account_url" {
  description = "URL del frontend para activación de cuenta"
  type        = string
  default     = "http://localhost:3000/activate-account"
}

variable "mail_host" {
  description = "Host del servidor de correo"
  type        = string
  default     = "smtp.gmail.com"
}

variable "mail_port" {
  description = "Puerto del servidor de correo"
  type        = number
  default     = 587
}

variable "mail_password" {
  description = "Contraseña del servidor de correo"
  type        = string
  sensitive   = true
  default     = ""
}

variable "mail_from" {
  description = "Email remitente"
  type        = string
  default     = "noreply@gormgoskeleton.com"
}
