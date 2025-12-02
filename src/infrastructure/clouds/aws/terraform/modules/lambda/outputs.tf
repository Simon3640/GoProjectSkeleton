output "lambda_function_arn" {
  description = "ARN of the Lambda function"
  value       = aws_lambda_function.main.arn
}

output "lambda_function_name" {
  description = "Name of the Lambda function"
  value       = aws_lambda_function.main.function_name
}

output "lambda_function_url" {
  description = "URL of the Lambda function (equivalent to Azure Function App default hostname)"
  value       = aws_lambda_function_url.main.function_url
}

output "lambda_function_invoke_arn" {
  description = "Invoke ARN of the Lambda function"
  value       = aws_lambda_function.main.invoke_arn
}

output "iam_role_arn" {
  description = "ARN of the IAM role (equivalent to Azure Managed Identity)"
  value       = aws_iam_role.lambda.arn
}

output "cloudwatch_log_group_name" {
  description = "CloudWatch Log Group name (equivalent to Application Insights)"
  value       = aws_cloudwatch_log_group.lambda.name
}
