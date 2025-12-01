// API Gateway HTTP API for Lambda functions
// This gateway exposes all Lambda functions defined in functions.json

resource "aws_apigatewayv2_api" "http_api" {
  name          = "${local.name_prefix}-http-api"
  protocol_type = "HTTP"

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-http-api"
    }
  )
}

# Default stage with auto deploy
resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.http_api.id
  name        = "$default"
  auto_deploy = true

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-http-api-stage"
    }
  )
}

# Integration for each Lambda function
resource "aws_apigatewayv2_integration" "lambda" {
  for_each = local.functions_map

  api_id                = aws_apigatewayv2_api.http_api.id
  integration_type      = "AWS_PROXY"
  integration_method    = "POST"
  payload_format_version = "2.0"

  # Use the Lambda function ARN from the module
  integration_uri = module.lambda_functions[each.key].lambda_function_arn
}

# Route for each Lambda function
# We prefix routes with /api to keep a clean namespace (e.g., /api/health-check)
resource "aws_apigatewayv2_route" "lambda_route" {
  for_each = local.functions_map

  api_id = aws_apigatewayv2_api.http_api.id

  # Example: "GET /api/health-check"
  route_key = "${upper(each.value.method)} /api/${each.value.route}"

  target = "integrations/${aws_apigatewayv2_integration.lambda[each.key].id}"
}

# Lambda permission to allow API Gateway to invoke the functions
resource "aws_lambda_permission" "apigw_invoke" {
  for_each = local.functions_map

  statement_id  = "AllowAPIGatewayInvoke-${each.key}"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_functions[each.key].lambda_function_name
  principal     = "apigateway.amazonaws.com"

  # Allow any route/stage of this API to invoke the function
  source_arn = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}
