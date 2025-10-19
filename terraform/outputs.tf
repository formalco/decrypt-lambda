output "api_gateway_url" {
  description = "URL of the API Gateway endpoint"
  value       = "${aws_api_gateway_stage.decrypt_stage.invoke_url}/decrypt"
}

output "api_gateway_id" {
  description = "ID of the API Gateway REST API"
  value       = aws_api_gateway_rest_api.decrypt_api.id
}

output "lambda_function_name" {
  description = "Name of the Lambda function"
  value       = aws_lambda_function.decrypt.function_name
}

output "lambda_function_arn" {
  description = "ARN of the Lambda function"
  value       = aws_lambda_function.decrypt.arn
}

output "lambda_role_arn" {
  description = "ARN of the Lambda IAM role"
  value       = aws_iam_role.decrypt_lambda_role.arn
}

output "cloudwatch_log_group_name" {
  description = "Name of the CloudWatch Log Group"
  value       = aws_cloudwatch_log_group.decrypt_lambda_logs.name
}

output "vpc_endpoint_id" {
  description = "ID of the API Gateway VPC Endpoint"
  value       = aws_vpc_endpoint.apigw_endpoint.id
}

output "vpc_endpoint_dns_entries" {
  description = "DNS entries for the VPC endpoint"
  value       = aws_vpc_endpoint.apigw_endpoint.dns_entry
}

output "vpc_endpoint_private_ips" {
  description = "Private IP addresses of the VPC endpoint ENIs"
  value       = [for eni in data.aws_network_interface.vpc_endpoint_enis : eni.private_ip]
}

output "access_instructions" {
  description = "Instructions for accessing the private API"
  value       = <<-EOT
    This is a PRIVATE API Gateway, accessible only from within the VPC.

    API Gateway URL: ${aws_api_gateway_stage.decrypt_stage.invoke_url}/decrypt
    API Gateway ID: ${aws_api_gateway_rest_api.decrypt_api.id}
    VPC Endpoint ID: ${aws_vpc_endpoint.apigw_endpoint.id}
    VPC Endpoint Private IPs: ${join(", ", [for eni in data.aws_network_interface.vpc_endpoint_enis : eni.private_ip])}

    === From within VPC (EC2/ECS) ===
    curl -X POST ${aws_api_gateway_stage.decrypt_stage.invoke_url}/decrypt \
      -H "Content-Type: application/json" \
      -d '{"encrypted_data": "your-encrypted-data"}'

    DNS Resolution:
    - Private DNS is enabled on the VPC endpoint
    - The API Gateway hostname will resolve to private IPs within your VPC
  EOT
}
