## Purpose
 Deploy a Lambda function within your infrastructure to enable Formal HTTP Sidecar to decrypt sensitive data.

## How to deploy the lambda?
Run the following commands:

```
make deploy
sls deploy
```

Then copy the URL of the lambda function and paste it into Formal console as the `Decryption API URL`
