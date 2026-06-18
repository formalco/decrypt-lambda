# Decrypt Lambda

Reference decryptor for Formal's field-level log encryption. Formal encrypts sensitive log fields client-side and never holds the private key, so you deploy this Lambda in your own infrastructure to let users decrypt those fields on demand from their browser.

Once deployed, use the API Gateway URL as the `decryptor_uri` on a Formal encryption key.

**Note: we highly encourage making sure the API Gateway is only accessible via a VPN to prevent users outside of your organization from making requests to the /decrypt endpoint.**

## How it works

Encrypted fields are [JWE](https://datatracker.ietf.org/doc/html/rfc7516) compact tokens (`alg=RSA-OAEP-256`, `enc=A256GCM`): a fresh AES-256-GCM content key is wrapped per record with the RSA public half of your KMS key. The JWE `kid` is a key URI of the form `<scheme>://<keyID>`, e.g. `aws-kms://arn:aws:kms:us-east-1:123456789012:key/abcd`.

On each request the Lambda:

1. parses the JWE and reads the `kid`;
2. picks a provider from the URI scheme and asks it to unwrap the content key (for AWS, `kms:Decrypt` with `RSAES_OAEP_SHA_256`);
3. decrypts the content with the unwrapped key and returns the plaintext.

The Lambda holds no private key material; only the KMS unwrap can recover the content key.

## Deployment

There are three deployment methods: Terraform, Serverless (via Cloudformation), and Docker.

### Deploying via Terraform (Recommended)

To deploy via Terraform, we recommend incorporating the configuration template provided in the `terraform` directory into your Terraform setup.
To deploy the configuration as-is, run `make deploy-terraform` with your AWS credentials and with the Terraform CLI installed. This deployment deploys the API Gateway and Lambda in a *private* subnet within your VPC.

### Deploying via Serverless

To deploy via Serverless, run `make deploy-sls` with your Serverless credentials. Note: you will need a Serverless licesnse, AWS Account, and the Serverless CLI installed. This deployment deploys the API Gateway and Lambda *publicly.*

### Deploying via Docker

To deploy via Docker, use the provided Dockerfile to build and push to an ECR repo.
The resulting container image can be used to deploy a lambda function [as a container image](https://docs.aws.amazon.com/lambda/latest/dg/go-image.html#go-image-other).

## Tests

`go test ./...` covers JWE parsing and decryption with a stubbed KMS unwrap.

An integration test drives the full path (provider registry through to AWS KMS `Decrypt`) against a real KMS. Point it at any KMS reachable through the default AWS credential chain, or at a local [LocalStack](https://www.localstack.cloud) KMS:

```bash
go test -tags integration -run TestDecryptViaKMS ./...
```

The endpoint defaults to `http://localhost:4566`; override it with `TEST_KMS_ENDPOINT`.
