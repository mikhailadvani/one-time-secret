package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mikhailadvani/one-time-secret/pkg/api"
	configf "github.com/mikhailadvani/one-time-secret/pkg/config"
)

var getEndpointBase = "/api/secret"
var config = configf.LoadConfig()

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	baseURL := fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com", request.RequestContext.APIID, config.AwsRegion)
	responseURLPrefix := fmt.Sprintf("%s%s", baseURL, getEndpointBase)
	decoder := json.NewDecoder(strings.NewReader(request.Body))
	var createSecretRequest api.CreateSecretRequest
	err := decoder.Decode(&createSecretRequest)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	secret, err := api.CreateSecret(createSecretRequest, responseURLPrefix)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: secret.Status,
		}, nil
	}
	jsonBody, err := json.Marshal(secret)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonBody),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
