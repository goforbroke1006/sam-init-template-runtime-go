package main

import (
	"errors"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	UndefinedTimezone = errors.New("undefined timezone")
	NotNullTimezone   = errors.New("timezone is required")
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	timezone, exists := request.QueryStringParameters["timezone"]
	var timeStr string
	if exists {
		tz, err := time.LoadLocation(timezone)
		if nil != err {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       UndefinedTimezone.Error(),
			}, err
		}
		timeStr = time.Now().In(tz).Format(time.RFC3339)
	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       NotNullTimezone.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		Body:       timeStr,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
