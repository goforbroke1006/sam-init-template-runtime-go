package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	username, exist := request.QueryStringParameters["username"]

	var message string

	if exist {
		message = "Hello, " + username + "!"
	} else {
		message = "Anon, please add 'username' GET-param to URL and you will see MAGIC!! ;D"
	}

	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
