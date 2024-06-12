package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"stori-challenge/cmd/api"
)

func handlerSqs(ctx context.Context, sqsEvent events.SQSEvent) {
	transactionService := api.BindTransactionService()

	transactionService.SaveTransactions(ctx, sqsEvent)
}

func main() {
	lambda.Start(handlerSqs)
}
