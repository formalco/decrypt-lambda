package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
)

type Response struct {
	Message string `json:"message"`
}

var (
	headers = map[string]string{
		"Access-Control-Allow-Origin":  "https://app.joinformal.com, https://app.datadoghq.com, https://app.datadoghq.eu",
		"Access-Control-Allow-Methods": "OPTIONS,POST",
	}
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Debug().Msgf("Received event: %s", event.Body)

	parsed, err := parseFormalEncryptedData(event.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error happen in parseFormalEncryptedData")
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: err.Error(), Headers: headers}, nil
	}

	log.Debug().Msgf("Parsed body: %+v", parsed)

	dataKey, err := decryptDataKey(parsed.KmsKeyRegion, parsed.KmsKeyId, []byte(parsed.EncryptedKey))
	if err != nil {
		log.Error().Err(err).Msg("Error happen in decryptDataKey")
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error(), Headers: headers}, nil
	}

	decrypted, err := decryptString(parsed.EncryptedData, dataKey)
	if err != nil {
		log.Error().Err(err).Msg("Error happen in decryptString")
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error(), Headers: headers}, nil
	}

	log.Debug().Msgf("Decrypted data: %s", decrypted)

	response := Response{Message: decrypted}
	responseBody, err := json.Marshal(&response)
	if err != nil {
		log.Error().Err(err).Msg("Error happen in json.Marshal")
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error marshalling the response body", Headers: headers}, nil
	}

	log.Debug().Msgf("Response body: %s", responseBody)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(handler)
}
