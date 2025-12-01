# AWS Secrets Manager
# Equivalent to Azure Key Vault (azurerm_key_vault and azurerm_key_vault_secret)

# Secrets Manager Secret for database password
resource "aws_secretsmanager_secret" "db_password" {
  name                    = "${local.name_prefix}/db-password"
  description             = "Database password for ${var.project_name}"
  recovery_window_in_days = 7 # Equivalent to Azure soft delete retention

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-db-password"
    }
  )
}

resource "aws_secretsmanager_secret_version" "db_password" {
  secret_id     = aws_secretsmanager_secret.db_password.id
  secret_string = var.db_password
}

# Secrets Manager Secret for Redis password
# Will be updated with ElastiCache auth token after creation
resource "aws_secretsmanager_secret" "redis_password" {
  count                   = var.use_redis_cache ? 1 : 0
  name                    = "${local.name_prefix}/redis-password"
  description             = "Redis password for ${var.project_name}"
  recovery_window_in_days = 7

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-redis-password"
    }
  )
}

resource "aws_secretsmanager_secret_version" "redis_password" {
  count         = var.use_redis_cache ? 1 : 0
  secret_id     = aws_secretsmanager_secret.redis_password[0].id
  secret_string = var.use_redis_cache && length(aws_elasticache_replication_group.redis) > 0 ? aws_elasticache_replication_group.redis[0].auth_token : var.redis_password

  depends_on = [aws_elasticache_replication_group.redis]
}

# Secrets Manager Secret for JWT secret
resource "aws_secretsmanager_secret" "jwt_secret" {
  name                    = "${local.name_prefix}/jwt-secret"
  description             = "JWT secret key for ${var.project_name}"
  recovery_window_in_days = 7

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-jwt-secret"
    }
  )
}

resource "aws_secretsmanager_secret_version" "jwt_secret" {
  secret_id     = aws_secretsmanager_secret.jwt_secret.id
  secret_string = var.jwt_secret_key
}

# Secrets Manager Secret for SendGrid API key
resource "aws_secretsmanager_secret" "sendgrid_key" {
  name                    = "${local.name_prefix}/sendgrid-api-key"
  description             = "SendGrid API key for ${var.project_name}"
  recovery_window_in_days = 7

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-sendgrid-key"
    }
  )
}

resource "aws_secretsmanager_secret_version" "sendgrid_key" {
  secret_id     = aws_secretsmanager_secret.sendgrid_key.id
  secret_string = var.mail_password
}

# IAM Policy for Lambda functions to access Secrets Manager
# Equivalent to Azure Key Vault Access Policy
resource "aws_iam_policy" "secrets_manager_access" {
  name        = "${local.name_prefix}-secrets-manager-access"
  description = "Policy for Lambda functions to access Secrets Manager"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue",
          "secretsmanager:DescribeSecret"
        ]
        Resource = [
          aws_secretsmanager_secret.db_password.arn,
          aws_secretsmanager_secret.jwt_secret.arn,
          aws_secretsmanager_secret.sendgrid_key.arn,
          var.use_redis_cache ? aws_secretsmanager_secret.redis_password[0].arn : ""
        ]
      }
    ]
  })

  # Note: Tags removed to avoid iam:TagPolicy permission requirement
  # If you have iam:TagPolicy permission, you can uncomment the tags block below
  # tags = merge(
  #   local.common_tags,
  #   {
  #     Name = "${local.name_prefix}-secrets-manager-policy"
  #   }
  # )
}
