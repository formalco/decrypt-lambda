## Purpose
 Deploy a Lambda function within your infrastructure to enable Formal HTTP Sidecar to decrypt sensitive data.

## How to deploy the lambda?

You can deploy this AWS Lambda function in two ways:
1. Using the serverless framework
To deploy using the Serverless framework, run the following commands:
```
npm i -E serverless@4.21.1 -g
make deploy
sls deploy
```

Once deployed you will be able to use the API url provided by AWS API Gateway to configure the Formal frontend to call the lambda and decrypt encrypted log payload.

2. Using Docker and the provided Dockerfile
To deploy using Docker, you should follow those steps:
1. Build the docker image
2. Push the docker image to your docker registry (ECR)
3. Create a lambda function using the docker image as the runtime
4. Create an API Gateway to expose the lambda function



