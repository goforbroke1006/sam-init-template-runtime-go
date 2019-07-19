package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func retrieveEventHandler() {
	log.Println("handle event here...")
}

func main() {
	lambda.Start(retrieveEventHandler)
}
