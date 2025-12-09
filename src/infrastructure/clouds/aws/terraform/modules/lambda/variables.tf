# General variables
variable "function_name" {
  description = "Function name (used for resource naming)"
  type        = string
}

variable "name_prefix" {
  description = "Prefix for resource names"
  type        = string
}

variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
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
}

variable "db_port" {
  description = "Database port"
  type        = number
  default     = 5432
}

variable "db_user" {
  description = "Database user"
  type        = string
}

variable "db_name" {
  description = "Database name"
  type        = string
}

variable "db_ssl" {
  description = "Enable SSL for database"
  type        = bool
  default     = true
}

variable "db_password_secret_arn" {
  description = "ARN of database password secret in Secrets Manager"
  type        = string
}

# Redis variables
variable "redis_host" {
  description = "Redis host"
  type        = string
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

variable "redis_password_secret_arn" {
  description = "ARN of Redis password secret in Secrets Manager"
  type        = string
  default     = ""
}

# JWT variables
variable "jwt_secret_arn" {
  description = "ARN of JWT secret in Secrets Manager"
  type        = string
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

variable "mail_from" {
  description = "Sender email"
  type        = string
  default     = "noreply@goprojectskeleton.com"
}

variable "mail_password_secret_arn" {
  description = "ARN of mail password secret in Secrets Manager"
  type        = string
}

# VPC and networking variables
variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "subnet_ids" {
  description = "Subnet IDs for Lambda VPC configuration"
  type        = list(string)
  default     = []
}

variable "security_group_ids" {
  description = "Security group IDs for Lambda"
  type        = list(string)
  default     = []
}

# IAM variables
variable "secrets_manager_policy_arn" {
  description = "ARN of Secrets Manager access policy"
  type        = string
}

variable "s3_templates_policy_arn" {
  description = "ARN of S3 templates read policy"
  type        = string
}

# Additional variables
variable "extra_environment_variables" {
  description = "Additional environment variables"
  type        = map(string)
  default     = {}
}
