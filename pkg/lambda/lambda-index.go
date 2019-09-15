package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mikhailadvani/one-time-secret/pkg/api"
)

var createEndpoint = "/api/v1/secret"

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"content-type": "text/html"},
		Body:       api.Index(createEndpoint),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
