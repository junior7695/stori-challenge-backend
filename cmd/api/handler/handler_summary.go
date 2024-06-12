package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"stori-challenge/cmd/api"
)

func handlerSummaryTransaction(ctx context.Context) {
	transactionService := api.BindTransactionService()

	transactionService.SendMonthlyReport(ctx)
}

func main() {
	lambda.Start(handlerSummaryTransaction)
}
