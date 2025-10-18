# Serverless Deployment Guide

## Prerequisites

- A Serverless license
- The Serverless CLI (`npm i -E serverless@4.21.1 -g`)
- AWS Credentials with the ability to deploy API Gateways, Lambdas, EC2 instances and the associated networking.


## A note about private deployments

This configuration deploys an AWS Lambda and API Gateway on the public internet. Instead, we recommend modifying this configuration to deploy this in a private subnet and require access to the endpoint via a VPN.

## Deployment steps

To deploy using the Serverless framework, run the following commands:
```
npm i -E serverless@4.21.1 -g
make deploy-sls
```

Note: you will need a license from Serverless to use the CLI. This will create a cloudformation stack in your AWS environment.
