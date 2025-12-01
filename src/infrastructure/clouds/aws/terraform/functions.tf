# Read functions.json file and create Lambda functions
# Equivalent to Azure Functions module creation from JSON
locals {
  functions_json = jsondecode(file("${path.module}/../functions.json"))

  # Convert JSON to map for for_each
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

# Create a Lambda function for each function defined in JSON
# Equivalent to Azure Function App module
module "lambda_functions" {
  source = "./modules/lambda"

  for_each = local.functions_map

  # General variables
  function_name = replace(each.value.name, "-", "") # Clean hyphens for resource names
  name_prefix   = local.name_prefix
  project_name  = var.project_name
  environment   = var.environment

  # Application variables
  app_port         = var.app_port
  app_version      = var.app_version
  app_description  = var.app_description
  app_support_email = var.app_support_email
  enable_log       = var.enable_log
  debug_log        = var.debug_log
  templates_path   = var.templates_path

  # Database variables
  db_host                = aws_db_instance.postgres.address
  db_port                = tonumber(aws_db_instance.postgres.port)
  db_user                = var.db_user
  db_name                = var.db_name
  db_ssl                 = var.db_ssl
  db_password_secret_arn = aws_secretsmanager_secret.db_password.arn

  # Redis variables
  # For non-clustered Redis (num_cache_clusters = 1), use primary_endpoint_address instead of configuration_endpoint_address
  redis_host                = var.use_redis_cache && length(aws_elasticache_replication_group.redis) > 0 ? "${aws_elasticache_replication_group.redis[0].primary_endpoint_address}:${tostring(aws_elasticache_replication_group.redis[0].port)}" : var.redis_host
  redis_db                  = var.redis_db
  redis_ttl                 = var.redis_ttl
  redis_password_secret_arn = var.use_redis_cache && length(aws_secretsmanager_secret.redis_password) > 0 ? aws_secretsmanager_secret.redis_password[0].arn : ""

  # JWT variables
  jwt_secret_arn = aws_secretsmanager_secret.jwt_secret.arn
  jwt_issuer     = var.jwt_issuer
  jwt_access_ttl = var.jwt_access_ttl
  jwt_refresh_ttl = var.jwt_refresh_ttl
  jwt_clock_skew  = var.jwt_clock_skew

  # Token and OTP variables
  one_time_token_ttl              = var.one_time_token_ttl
  one_time_token_email_verify_ttl = var.one_time_token_email_verify_ttl
  one_time_password_length         = var.one_time_password_length
  one_time_password_ttl            = var.one_time_password_ttl

  # Frontend variables
  frontend_reset_password_url   = var.frontend_reset_password_url
  frontend_activate_account_url = var.frontend_activate_account_url

  # Mail variables
  mail_host                = var.mail_host
  mail_port                = var.mail_port
  mail_from                = var.mail_from
  mail_password_secret_arn = aws_secretsmanager_secret.sendgrid_key.arn

  # VPC and networking
  vpc_id             = aws_vpc.main.id
  subnet_ids         = aws_subnet.private[*].id
  security_group_ids = [aws_security_group.lambda.id]

  # IAM
  secrets_manager_policy_arn = aws_iam_policy.secrets_manager_access.arn

  depends_on = [
    aws_db_instance.postgres,
    aws_elasticache_replication_group.redis,
    aws_secretsmanager_secret.db_password,
    aws_secretsmanager_secret.jwt_secret,
    aws_secretsmanager_secret.sendgrid_key,
    aws_iam_policy.secrets_manager_access
  ]
}
