# Decrypt Lambda

Deploy a Lambda function within your infrastructure to enable users to decrypt sensitive data from their browser.

Once deployed you will be able to use the API URL provided by the AWS API Gateway and add the URL as a decryptor URI for
a Formal Encryption Key.

**Note: we highly encourage making sure the API Gateway is only accessible via a VPN to prevent users outside of your organization from making requests to the /decrypt endpoint.**

There are three deployment methods: Terraform, Serverless (via Cloudformation), and Docker.

## Deploying via Terraform

To deploy via Terraform, we recommend incorporating the configuration template provided in the `terraform` directory into your Terraform setup.
To deploy the configuration as-is, run `make deploy-terraform` with your AWS credentials and with the Terraform CLI installed.

## Deploying via Serverless.

To deploy via Serverless, run `make deploy-sls` with your Serverless credentials. Note: you will need a Serverless licesnse and the Serverless CLI installed.


## Deploying via Docker

To deploy via Docker, use the provided Dockerfile to build and push to an ECR repo.
The resulting container image can be used to deploy a lambda function [as a container image](https://docs.aws.amazon.com/lambda/latest/dg/go-image.html#go-image-other).