# AWS Lambda via Terraform

This directory contains Terraform configurations to deploy the decryptor function with an API Gateway in a private subnet within your VPC.

This deployment will put the API Gateway and Lambda in a private subnet. We recommend accessing the resulting API Gateway decryptor URI via a VPN.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html)
- Go (to build the Lambda binary)
- An AWS VPC with a private subnet.
- AWS Credentials with the ability to deploy API Gateways, Lambdas, EC2 instances and the associated networking.

## Files

- `main.tf` - Main Terraform configuration (Lambda, API Gateway, IAM roles)
- `variables.tf` - Variable definitions
- `outputs.tf` - Output definitions
- `terraform.tfvars.example` - Example variables file

## Deployment steps

1. First, copy the example variables and update them.

```bash
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your preferred values
```

2. From the repo root, build and zip the Lambda binary, then apply:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bootstrap .
zip bootstrap.zip bootstrap
cd deploy/aws-lambda-terraform
terraform init && terraform apply
```
## Resources Created

- **Lambda Function**: `decryptor` (or custom name) deployed in private subnets
- **IAM Role**: `decryptor-role` with KMS decrypt and VPC access permissions
- **API Gateway**: Private REST API with POST /decrypt endpoint
- **VPC Endpoint**: Interface endpoint for API Gateway execute-api service
- **Security Groups**:
  - Lambda security group with egress to all
  - API Gateway VPC endpoint security group with ingress on port 443 from VPC CIDR
- **CORS Configuration**: Configured for https://app.formal.ai
- **CloudWatch Log Group**: For Lambda function logs

## VPC Configuration

This deployment creates a **private API Gateway** accessible only from within the VPC. The Lambda function runs in private subnets and connects to API Gateway through a VPC endpoint.

### VPC Requirements

- **VPC ID**: An existing VPC where resources will be deployed
- **Private Subnets**: At least 2 private subnets (recommended for high availability)
  - Subnets should have routes to a NAT Gateway if the Lambda needs internet access
  - Subnets should be in different Availability Zones for resilience
- **VPC Endpoints**: The terraform configuration automatically creates the required API Gateway execute-api endpoint

### Network Architecture

```
Client (within VPC) → VPC Endpoint (execute-api) → Private API Gateway → Lambda (in private subnet)
```

The API Gateway is not accessible from the public internet. We recommend requiring access through a VPN so that users can access the API Gateway from their browsers.

## Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `aws_region` | AWS region to deploy resources | `us-east-1` | No |
| `function_name` | Name of the Lambda function | `decryptor` | No |
| `stage_name` | API Gateway stage name | `prod` | No |
| `kms_key_arn` | ARN for KMS key we're using to decrypt | - | Yes |
| `vpc_id` | VPC ID where Lambda and API Gateway will be deployed | - | Yes |
| `private_subnet_ids` | List of private subnet IDs (recommend 2+) | - | Yes |
| `log_retention_days` | CloudWatch log retention in days | `14` | No |

## Outputs

| Output | Description |
|--------|-------------|
| `api_gateway_url` | Full URL of the API endpoint including /decrypt path |
| `api_gateway_id` | API Gateway REST API ID |
| `lambda_function_name` | Lambda function name |
| `lambda_function_arn` | Lambda function ARN |
| `lambda_role_arn` | Lambda IAM role ARN |
| `cloudwatch_log_group_name` | CloudWatch Log Group name |
| `vpc_endpoint_id` | VPC Endpoint ID for API Gateway |
| `vpc_endpoint_dns_entries` | DNS entries for the VPC endpoint |
| `vpc_endpoint_private_ips` | Private IP addresses of the VPC endpoint|
| `access_instructions` | Instructions for accessing the private API |
