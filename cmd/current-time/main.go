package main

import (
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	timeStr := time.Now().Format(time.RFC3339)
	return events.APIGatewayProxyResponse{
		Body:       timeStr,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
