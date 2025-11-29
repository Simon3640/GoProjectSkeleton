package main

import (
	"context"
	"log"

	"gormgoskeleton/aws"
	handlers "gormgoskeleton/src/infrastructure/handlers"

	"github.com/aws/aws-lambda-go/lambda"
)

var initialized bool

func init() {
	if !initialized {
		log.Println("Initializing AWS Infrastructure")
		aws.InitializeInfrastructure()
		initialized = true
		log.Println("AWS Infrastructure initialized successfully")
	}
}

func handler(ctx context.Context, event map[string]interface{}) (aws.LambdaResponse, error) {

	return aws.HandleLambdaEventWithAuth(event, func(hc handlers.HandlerContext) {
		handlers.CreatePassword(hc)
	})

}

func main() {
	lambda.Start(handler)
}
