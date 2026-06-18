package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"

	_ "decrypt-lambda/provider/awskms"
)

type Response struct {
	Message string `json:"message"`
}

var (
	headers = map[string]string{
		"Access-Control-Allow-Origin":  "https://app.formal.ai",
		"Access-Control-Allow-Methods": "OPTIONS,POST",
	}
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Debug().Msgf("Received event: %s", event.Body)

	obj, err := parseJWE(event.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse JWE from request body")
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "request body is not a valid JWE", Headers: headers}, nil
	}

	decrypted, err := decryptJWE(ctx, obj)
	if err != nil {
		log.Error().Err(err).Msg("failed to decrypt JWE")
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "could not decrypt the JWE", Headers: headers}, nil
	}

	responseBody, err := json.Marshal(&Response{Message: decrypted})
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal decrypt response")
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "internal error", Headers: headers}, nil
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
