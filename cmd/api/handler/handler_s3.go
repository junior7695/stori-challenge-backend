package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"stori-challenge/cmd/api"
)

func handlerS3(ctx context.Context, event events.S3Event) {
	transactionService := api.BindTransactionService()

	transactionService.ReadS3Files(ctx, event)
}

func main() {
	lambda.Start(handlerS3)
}
