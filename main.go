package main

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
}

var (
	headers = map[string]string{
		"Access-Control-Allow-Origin":  "https://app.joinformal.com",
		"Access-Control-Allow-Methods": "OPTIONS,POST",
	}
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	parsed, err := parseFormalEncryptedData(event.Body)
	if err != nil {
		slog.Error("Error happen in parseFormalEncryptedData", err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: err.Error(), Headers: headers}, nil
	}

	dataKey, err := decryptDataKey(parsed.KmsKeyRegion, parsed.KmsKeyId, []byte(parsed.EncryptedKey))
	if err != nil {
		slog.Error("Error happen in decryptDataKey", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error(), Headers: headers}, nil
	}

	decrypted, err := decryptString(parsed.EncryptedData, dataKey)
	if err != nil {
		slog.Error("Error happen in decryptString", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error(), Headers: headers}, nil
	}

	response := Response{Message: string(decrypted)}
	responseBody, err := json.Marshal(&response)
	if err != nil {
		slog.Error("Error happen in json.Marshal", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error marshalling the response body", Headers: headers}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(handler)
}
