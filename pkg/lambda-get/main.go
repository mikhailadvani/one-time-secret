package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mikhailadvani/one-time-secret/pkg/api"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	secretID := request.PathParameters["secretID"]
	secret, err := api.GetSecret(secretID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: secret.Status,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "text/html"},
		Body:       secret.Content,
		StatusCode: secret.Status,
	}, nil
}

func main() {
	lambda.Start(handler)
}
