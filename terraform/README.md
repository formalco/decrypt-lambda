# Terraform Deployment Guide

This directory contains Terraform configurations to deploy the decrypt Lambda function with API Gateway, equivalent to the serverless configuration.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- AWS CLI configured with appropriate credentials
- The `bootstrap` binary compiled and ready for deployment

## Files

- `main.tf` - Main Terraform configuration (Lambda, API Gateway, IAM roles)
- `variables.tf` - Variable definitions
- `outputs.tf` - Output definitions
- `terraform.tfvars.example` - Example variables file

## Deployment Steps

### 1. Prepare the Lambda Package

Before deploying, you need to package the `bootstrap` binary:

```bash
# Compile your Go code if not already done
# GOOS=linux GOARCH=arm64 go build -o bootstrap main.go crypto.go

# Create a zip file for Lambda
zip bootstrap.zip bootstrap
```

### 2. Initialize Terraform

```bash
terraform init
```

### 3. Configure Variables (Optional)

Copy the example variables file and customize if needed:

```bash
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your preferred values
```

### 4. Plan the Deployment

Review what resources will be created:

```bash
terraform plan
```

### 5. Apply the Configuration

Deploy the infrastructure:

```bash
terraform apply
```

Type `yes` when prompted to confirm.

### 6. Get the API Endpoint

After deployment, the API Gateway URL will be displayed:

```bash
terraform output api_gateway_url
```

## Resources Created

- **Lambda Function**: `decrypt-lambda` (or custom name)
- **IAM Role**: `decrypt-lambda-role` with KMS decrypt permissions
- **API Gateway**: REST API with POST /decrypt endpoint
- **CORS Configuration**: Configured for https://app.joinformal.com (or custom origin)
- **CloudWatch Log Group**: For Lambda function logs

## Updating the Lambda Function

After making changes to your code:

1. Rebuild the bootstrap binary
2. Recreate the zip file: `zip bootstrap.zip bootstrap`
3. Run `terraform apply` to update the Lambda function

## Cleanup

To destroy all resources:

```bash
terraform destroy
```

## Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `aws_region` | AWS region to deploy resources | `us-east-1` |
| `function_name` | Name of the Lambda function | `decrypt-lambda` |
| `stage_name` | API Gateway stage name | `prod` |
| `kms_key_arn` | ARN for KMS key we're using to decrypt | `` |
| `log_retention_days` | CloudWatch log retention in days | `14` |

## Outputs

| Output | Description |
|--------|-------------|
| `api_gateway_url` | Full URL of the API endpoint including /decrypt path |
| `api_gateway_id` | API Gateway REST API ID |
| `lambda_function_name` | Lambda function name |
| `lambda_function_arn` | Lambda function ARN |
| `lambda_role_arn` | Lambda IAM role ARN |
| `cloudwatch_log_group_name` | CloudWatch Log Group name |
