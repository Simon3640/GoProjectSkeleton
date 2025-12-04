# Outputs for VPC
output "vpc_id" {
  description = "VPC ID (equivalent to Azure Resource Group)"
  value       = aws_vpc.main.id
}

output "vpc_cidr" {
  description = "VPC CIDR block"
  value       = aws_vpc.main.cidr_block
}

# Outputs for RDS PostgreSQL
output "postgres_endpoint" {
  description = "RDS PostgreSQL endpoint (equivalent to Azure PostgreSQL FQDN)"
  value       = aws_db_instance.postgres.address
}

output "postgres_port" {
  description = "RDS PostgreSQL port"
  value       = aws_db_instance.postgres.port
}

output "postgres_connection_string" {
  description = "PostgreSQL connection string (without password for security)"
  value       = "postgresql://${var.db_user}@${aws_db_instance.postgres.address}:${aws_db_instance.postgres.port}/${var.db_name}?sslmode=${var.db_ssl ? "require" : "disable"}"
  sensitive   = false
}

# Outputs for ElastiCache Redis
output "redis_endpoint" {
  description = "ElastiCache Redis endpoint (equivalent to Azure Redis Cache hostname)"
  value       = var.use_redis_cache && length(aws_elasticache_replication_group.redis) > 0 ? aws_elasticache_replication_group.redis[0].configuration_endpoint_address : null
}

output "redis_port" {
  description = "ElastiCache Redis port"
  value       = var.use_redis_cache && length(aws_elasticache_replication_group.redis) > 0 ? aws_elasticache_replication_group.redis[0].port : 6379
}

output "redis_type" {
  description = "Type of Redis used"
  value       = var.use_redis_cache ? "AWS ElastiCache Redis" : "Local/None"
}

output "redis_primary_key" {
  description = "Redis password (stored in Secrets Manager)"
  value       = "See in Secrets Manager: ${local.name_prefix}/redis-password"
  sensitive   = true
}

# Outputs for Secrets Manager
output "secrets_manager_arn" {
  description = "Secrets Manager ARN (equivalent to Azure Key Vault URI)"
  value       = aws_secretsmanager_secret.db_password.arn
}

# Outputs for Lambda Functions
output "lambda_functions" {
  description = "Map of all Lambda functions created from functions.json (equivalent to Azure Function Apps)"
  value = {
    for k, v in module.lambda_functions : k => {
      name            = v.lambda_function_name
      arn             = v.lambda_function_arn
      invoke_url      = v.lambda_function_url
      cloudwatch_logs = v.cloudwatch_log_group_name
    }
  }
}

output "lambda_function_names" {
  description = "List of all Lambda function names"
  value       = [for k, v in module.lambda_functions : v.lambda_function_name]
}

# Outputs for API Gateway HTTP API
output "http_api_id" {
  description = "ID of the HTTP API Gateway"
  value       = aws_apigatewayv2_api.http_api.id
}

output "http_api_endpoint" {
  description = "Base URL of the HTTP API Gateway"
  value       = aws_apigatewayv2_api.http_api.api_endpoint
}

output "http_api_routes" {
  description = "Map of routes configured in the HTTP API"
  value = {
    for k, v in aws_apigatewayv2_route.lambda_route : k => {
      route_key = v.route_key
    }
  }
}

# Outputs for S3 Templates Bucket
output "s3_templates_bucket_name" {
  description = "Name of the S3 bucket for templates"
  value       = aws_s3_bucket.templates.id
}

output "s3_templates_bucket_arn" {
  description = "ARN of the S3 bucket for templates"
  value       = aws_s3_bucket.templates.arn
}

output "templates_path" {
  description = "S3 path for templates (to be used as TEMPLATES_PATH environment variable)"
  value       = "s3://${aws_s3_bucket.templates.bucket}/templates/"
}

output "s3_templates_policy_arn" {
  description = "ARN of the IAM policy for reading templates from S3"
  value       = aws_iam_policy.s3_templates_read.arn
}
