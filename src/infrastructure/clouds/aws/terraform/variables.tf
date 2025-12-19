# General variables
variable "project_name" {
  description = "Project name"
  type        = string
  default     = "go-project-skeleton"
}

variable "aws_region" {
  description = "AWS region (equivalent to Azure location)"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name (development, staging, production)"
  type        = string
  default     = "development"
}

# Application variables
variable "app_port" {
  description = "Application port"
  type        = string
  default     = "8080"
}

variable "app_version" {
  description = "Application version"
  type        = string
  default     = "0.0.1"
}

variable "app_description" {
  description = "Application description"
  type        = string
  default     = "Gorm Go Skeleton Template"
}

variable "app_support_email" {
  description = "Support email"
  type        = string
  default     = "support@goprojectskeleton.com"
}

variable "enable_log" {
  description = "Enable logging"
  type        = bool
  default     = true
}

variable "debug_log" {
  description = "Enable debug logging"
  type        = bool
  default     = true
}

variable "templates_path" {
  description = "Templates path"
  type        = string
  default     = "src/application/shared/templates/"
}

# Database variables
variable "db_host" {
  description = "Database host"
  type        = string
  default     = "localhost"
}

variable "db_port" {
  description = "Database port"
  type        = string
  default     = "5432"
}

variable "db_user" {
  description = "Database user"
  type        = string
  default     = "goprojectskeleton"
}

variable "db_password" {
  description = "Database password. MUST be a secure password. Generate with: openssl rand -base64 32"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "goprojectskeleton"
}

variable "db_ssl" {
  description = "Enable SSL for database"
  type        = bool
  default     = true
}

# Redis variables
variable "redis_host" {
  description = "Redis host"
  type        = string
  default     = "localhost:6379"
}

variable "redis_password" {
  description = "Redis password (not used directly, obtained from ElastiCache when use_redis_cache=true)"
  type        = string
  sensitive   = true
  default     = ""
}

variable "redis_db" {
  description = "Redis database number"
  type        = number
  default     = 0
}

variable "redis_ttl" {
  description = "Redis TTL in seconds"
  type        = number
  default     = 300
}

variable "use_redis_cache" {
  description = "Use AWS ElastiCache Redis (equivalent to Azure Redis Cache)"
  type        = bool
  default     = true
}

# JWT variables
variable "jwt_secret_key" {
  description = "JWT secret key. MUST be a secure random string (minimum 32 characters recommended). Generate with: openssl rand -hex 32"
  type        = string
  sensitive   = true
}

variable "jwt_issuer" {
  description = "JWT issuer"
  type        = string
  default     = "goprojectskeleton"
}

variable "jwt_access_ttl" {
  description = "JWT access token TTL in seconds"
  type        = number
  default     = 3600
}

variable "jwt_refresh_ttl" {
  description = "JWT refresh token TTL in seconds"
  type        = number
  default     = 86400
}

variable "jwt_clock_skew" {
  description = "JWT clock skew in seconds"
  type        = number
  default     = 60
}

# Token and OTP variables
variable "one_time_token_ttl" {
  description = "One-time token TTL in minutes"
  type        = number
  default     = 15
}

variable "one_time_token_email_verify_ttl" {
  description = "Email verification token TTL in minutes"
  type        = number
  default     = 60
}

variable "one_time_password_length" {
  description = "One-time password length"
  type        = number
  default     = 6
}

variable "one_time_password_ttl" {
  description = "One-time password TTL in minutes"
  type        = number
  default     = 10
}

# Frontend variables
variable "frontend_reset_password_url" {
  description = "Frontend reset password URL"
  type        = string
  default     = "http://localhost:3000/reset-password"
}

variable "frontend_activate_account_url" {
  description = "Frontend activate account URL"
  type        = string
  default     = "http://localhost:3000/activate-account"
}

# Mail variables
variable "mail_host" {
  description = "Mail server host"
  type        = string
  default     = "smtp.gmail.com"
}

variable "mail_port" {
  description = "Mail server port"
  type        = number
  default     = 587
}

variable "mail_password" {
  description = "Mail server password"
  type        = string
  sensitive   = true
  default     = ""
}

variable "mail_from" {
  description = "Sender email"
  type        = string
  default     = "noreply@goprojectskeleton.com"
}

# VPC variables (AWS specific)
variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "Availability zones for multi-AZ deployment"
  type        = list(string)
  default     = []
}

# Observability variables
variable "observability_enabled" {
  description = "Enable observability"
  type        = bool
  default     = true
}

variable "observability_backend" {
  description = "Observability backend"
  type        = string
  default     = "opentelemetry"
}

variable "otlp_endpoint" {
  description = "OTLP endpoint"
  type        = string
  default     = "http://localhost:4318"
}

variable "observability_sampling_rate" {
  description = "Observability sampling rate"
  type        = number
  default     = 0.1
}
