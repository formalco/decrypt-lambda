terraform {
  required_version = ">= 0.14.9"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

resource "aws_iam_role" "decrypt_lambda_role" {
  name = "${var.function_name}-role"

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
}

resource "aws_iam_role_policy" "kms_decrypt_policy" {
  name = "kms-decrypt-policy"
  role = aws_iam_role.decrypt_lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "kms:Decrypt"
        ]
        Resource = var.kms_key_arn
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.decrypt_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_function" "decrypt" {
  filename         = "../bootstrap.zip"
  function_name    = var.function_name
  role            = aws_iam_role.decrypt_lambda_role.arn
  handler         = "bootstrap"
  source_code_hash = filebase64sha256("../bootstrap.zip")
  runtime         = "provided.al2"
  architectures   = ["arm64"]

  environment {
    variables = {}
  }
}

resource "aws_cloudwatch_log_group" "decrypt_lambda_logs" {
  name              = "/aws/lambda/${var.function_name}"
  retention_in_days = var.log_retention_days
}

resource "aws_api_gateway_rest_api" "decrypt_api" {
  name        = "${var.function_name}-api"
  description = "API Gateway for decrypt Lambda function"
}

resource "aws_api_gateway_resource" "decrypt_resource" {
  rest_api_id = aws_api_gateway_rest_api.decrypt_api.id
  parent_id   = aws_api_gateway_rest_api.decrypt_api.root_resource_id
  path_part   = "decrypt"
}

resource "aws_api_gateway_method" "decrypt_post" {
  rest_api_id   = aws_api_gateway_rest_api.decrypt_api.id
  resource_id   = aws_api_gateway_resource.decrypt_resource.id
  http_method   = "POST"
  authorization = "NONE"
}
resource "aws_api_gateway_method_response" "decrypt_post_response" {
  rest_api_id = aws_api_gateway_rest_api.decrypt_api.id
  resource_id = aws_api_gateway_resource.decrypt_resource.id
  http_method = aws_api_gateway_method.decrypt_post.http_method
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = true
  }
}
resource "aws_api_gateway_integration" "decrypt_lambda_integration" {
  rest_api_id             = aws_api_gateway_rest_api.decrypt_api.id
  resource_id             = aws_api_gateway_resource.decrypt_resource.id
  http_method             = aws_api_gateway_method.decrypt_post.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.decrypt.invoke_arn
}

resource "aws_api_gateway_integration_response" "decrypt_integration_response" {
  rest_api_id = aws_api_gateway_rest_api.decrypt_api.id
  resource_id = aws_api_gateway_resource.decrypt_resource.id
  http_method = aws_api_gateway_method.decrypt_post.http_method
  status_code = aws_api_gateway_method_response.decrypt_post_response.status_code

  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = "'https://app.joinformal.com'"
  }

  depends_on = [aws_api_gateway_integration.decrypt_lambda_integration]
}

# Support the OPTION method for CORS preflight
resource "aws_api_gateway_method" "decrypt_options" {
  rest_api_id   = aws_api_gateway_rest_api.decrypt_api.id
  resource_id   = aws_api_gateway_resource.decrypt_resource.id
  http_method   = "OPTIONS"
  authorization = "NONE"
}

resource "aws_api_gateway_method_response" "decrypt_options_response" {
  rest_api_id = aws_api_gateway_rest_api.decrypt_api.id
  resource_id = aws_api_gateway_resource.decrypt_resource.id
  http_method = aws_api_gateway_method.decrypt_options.http_method
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true
    "method.response.header.Access-Control-Allow-Methods" = true
    "method.response.header.Access-Control-Allow-Origin"  = true
  }
}

resource "aws_api_gateway_integration" "decrypt_options_integration" {
  rest_api_id = aws_api_gateway_rest_api.decrypt_api.id
  resource_id = aws_api_gateway_resource.decrypt_resource.id
  http_method = aws_api_gateway_method.decrypt_options.http_method
  type        = "MOCK"

  request_templates = {
    "application/json" = "{\"statusCode\": 200}"
  }
}
resource "aws_api_gateway_integration_response" "decrypt_options_integration_response" {
  rest_api_id = aws_api_gateway_rest_api.decrypt_api.id
  resource_id = aws_api_gateway_resource.decrypt_resource.id
  http_method = aws_api_gateway_method.decrypt_options.http_method
  status_code = aws_api_gateway_method_response.decrypt_options_response.status_code

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token'"
    "method.response.header.Access-Control-Allow-Methods" = "'POST,OPTIONS'"
    "method.response.header.Access-Control-Allow-Origin"  = "'https://app.joinformal.com'"
  }

  depends_on = [aws_api_gateway_integration.decrypt_options_integration]
}

resource "aws_api_gateway_deployment" "decrypt_deployment" {
  rest_api_id = aws_api_gateway_rest_api.decrypt_api.id

  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.decrypt_resource.id,
      aws_api_gateway_method.decrypt_post.id,
      aws_api_gateway_method.decrypt_options.id,
      aws_api_gateway_integration.decrypt_lambda_integration.id,
      aws_api_gateway_integration.decrypt_options_integration.id,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [
    aws_api_gateway_integration.decrypt_lambda_integration,
    aws_api_gateway_integration.decrypt_options_integration,
  ]
}

resource "aws_api_gateway_stage" "decrypt_stage" {
  deployment_id = aws_api_gateway_deployment.decrypt_deployment.id
  rest_api_id   = aws_api_gateway_rest_api.decrypt_api.id
  stage_name    = var.stage_name
}

resource "aws_lambda_permission" "api_gateway_invoke" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.decrypt.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.decrypt_api.execution_arn}/*/*"
}
