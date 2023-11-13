package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type Response struct {
	Message string `json:"message"`
}

func Decrypt(encrypted string, key []byte) ([]byte, error) {
	enc, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesGCM.NonceSize()
	if len(enc) < nonceSize {
		return nil, errors.New("data not encrypted or not encrypted with a nonce")
	}
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{"Access-Control-Allow-Origin": "*"}
	str := event.Body
	strsplit := strings.Split(str, ":")
	if len(strsplit) < 5 {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid input", Headers: headers}, nil
	}

	encryptedKey, err := base64.StdEncoding.DecodeString(strsplit[2])
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: err.Error(), Headers: headers}, nil
	}

	kmsKeyId, err := base64.StdEncoding.DecodeString(strsplit[3])
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: err.Error(), Headers: headers}, nil
	}

	kmsKeyRegion, err := base64.StdEncoding.DecodeString(strsplit[4])
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: err.Error(), Headers: headers}, nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(string(kmsKeyRegion)),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error(), Headers: headers}, nil
	}

	svc := kms.New(sess, &aws.Config{Region: aws.String(string(kmsKeyRegion))})
	input := &kms.DecryptInput{
		CiphertextBlob: encryptedKey,
		KeyId:          aws.String(string(kmsKeyId)),
	}

	result, err := svc.Decrypt(input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error(), Headers: headers}, nil
	}

	decrypted, err := Decrypt(strsplit[1], result.Plaintext)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error(), Headers: headers}, nil
	}

	response := Response{
		Message: string(decrypted),
	}
	responseBody, err := json.Marshal(&response)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error marshalling the response body", Headers: headers}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    headers}, nil
}

func main() {
	lambda.Start(handler)
}
