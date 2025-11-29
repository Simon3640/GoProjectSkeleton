# Lambda Function Module
# Equivalent to Azure Function App module
# Creates a Lambda function with API Gateway integration, IAM roles, and CloudWatch logging

# Create a dummy zip file for Lambda deployment
# This is a placeholder - actual code should be deployed separately
data "archive_file" "lambda_zip" {
  type        = "zip"
  output_path = "/tmp/${var.name_prefix}-${var.function_name}-dummy.zip"

  source {
    content  = "dummy"
    filename = "dummy.txt"
  }
}

# IAM Role for Lambda function
# Equivalent to Azure Managed Identity
resource "aws_iam_role" "lambda" {
  name = "${var.name_prefix}-${var.function_name}-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name        = "${var.name_prefix}-${var.function_name}-lambda-role"
    Environment = var.environment
  }
}

# Attach basic Lambda execution policy
resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Attach VPC access policy (if Lambda needs VPC access)
resource "aws_iam_role_policy_attachment" "lambda_vpc" {
  count      = length(var.subnet_ids) > 0 ? 1 : 0
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# Attach Secrets Manager access policy
resource "aws_iam_role_policy_attachment" "lambda_secrets" {
  role       = aws_iam_role.lambda.name
  policy_arn = var.secrets_manager_policy_arn
}

# CloudWatch Log Group for Lambda
# Equivalent to Azure Application Insights
resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${var.name_prefix}-${var.function_name}"
  retention_in_days = 7

  tags = {
    Name        = "${var.name_prefix}-${var.function_name}-logs"
    Environment = var.environment
  }
}

# Lambda Function
# Equivalent to Azure Function App (azurerm_linux_function_app)
resource "aws_lambda_function" "main" {
  function_name = "${var.name_prefix}-${var.function_name}"
  description   = "Lambda function for ${var.function_name}"

  # Runtime and handler
  runtime     = "provided.al2" # Custom runtime for Go (equivalent to Azure Custom Handler)
  handler     = "bootstrap"
  timeout     = 30
  memory_size = 512

  # Environment variables
  # Equivalent to Azure Function App app_settings
  environment {
    variables = merge(
      {
        # Application
        APP_NAME          = var.project_name
        APP_ENV           = var.environment
        APP_PORT          = var.app_port
        APP_VERSION       = var.app_version
        APP_DESCRIPTION   = var.app_description
        APP_SUPPORT_EMAIL = var.app_support_email

        # Logging
        ENABLE_LOG     = tostring(var.enable_log)
        DEBUG_LOG      = tostring(var.debug_log)
        TEMPLATES_PATH = var.templates_path

        # Database - using Secrets Manager ARN (Lambda will fetch at runtime)
        DB_HOST     = var.db_host
        DB_PORT     = tostring(var.db_port)
        DB_USER     = var.db_user
        DB_NAME     = var.db_name
        DB_SSL      = tostring(var.db_ssl)
        DB_PASSWORD = var.db_password_secret_arn

        # Redis
        REDIS_HOST     = var.redis_host
        REDIS_DB       = tostring(var.redis_db)
        REDIS_TTL      = tostring(var.redis_ttl)
        REDIS_PASSWORD = var.redis_password_secret_arn

        # JWT
        JWT_SECRET      = var.jwt_secret_arn
        JWT_ISSUER      = var.jwt_issuer
        JWT_AUDIENCE    = var.jwt_issuer
        JWT_ACCESS_TTL  = tostring(var.jwt_access_ttl)
        JWT_REFRESH_TTL = tostring(var.jwt_refresh_ttl)
        JWT_CLOCK_SKEW  = tostring(var.jwt_clock_skew)

        # Tokens and OTP
        ONE_TIME_TOKEN_TTL              = tostring(var.one_time_token_ttl)
        ONE_TIME_TOKEN_EMAIL_VERIFY_TTL = tostring(var.one_time_token_email_verify_ttl)
        ONE_TIME_PASSWORD_LENGTH        = tostring(var.one_time_password_length)
        ONE_TIME_PASSWORD_TTL           = tostring(var.one_time_password_ttl)

        # Frontend
        FRONTEND_RESET_PASSWORD_URL   = var.frontend_reset_password_url
        FRONTEND_ACTIVATE_ACCOUNT_URL = var.frontend_activate_account_url

        # Mail
        MAIL_HOST     = var.mail_host
        MAIL_PORT     = tostring(var.mail_port)
        MAIL_FROM     = var.mail_from
        MAIL_PASSWORD = var.mail_password_secret_arn
      },
      var.extra_environment_variables
    )
  }

  # VPC configuration (if needed for RDS/ElastiCache access)
  vpc_config {
    subnet_ids         = var.subnet_ids
    security_group_ids = var.security_group_ids
  }

  # IAM role
  role = aws_iam_role.lambda.arn

  # Package type and deployment
  package_type     = "Zip"
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  # Dead letter queue (optional)
  # dead_letter_config {
  #   target_arn = aws_sqs_queue.dlq.arn
  # }

  # Tracing configuration (equivalent to Application Insights)
  tracing_config {
    mode = "Active" # Enable X-Ray tracing
  }

  tags = {
    Name        = "${var.name_prefix}-${var.function_name}"
    Environment = var.environment
  }

  depends_on = [
    aws_iam_role_policy_attachment.lambda_basic,
    aws_iam_role_policy_attachment.lambda_secrets,
    aws_cloudwatch_log_group.lambda
  ]
}

# Lambda Function URL (for direct HTTP access)
# Equivalent to Azure Function App default hostname
resource "aws_lambda_function_url" "main" {
  function_name      = aws_lambda_function.main.function_name
  authorization_type = "NONE" # Can be changed to AWS_IAM for authentication

  cors {
    allow_credentials = false
    allow_origins     = ["*"]
    allow_methods     = ["*"]
    allow_headers     = ["*"]
    expose_headers    = ["*"]
    max_age           = 3600
  }
}
