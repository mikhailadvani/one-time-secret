package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mikhailadvani/one-time-secret/pkg/api"
	configf "github.com/mikhailadvani/one-time-secret/pkg/config"
)

var getEndpointBase = "/api/secret"
var config = configf.LoadConfig()

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	responseURLPrefix := fmt.Sprintf("%s%s", config.BaseURL, getEndpointBase)
	createSecretRequest := api.CreateSecretRequest{Content: request.Body}
	secret, err := api.CreateSecret(createSecretRequest, responseURLPrefix)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: secret.Status,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"content-type": "text/html"},
		Body:       secret.URL,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
