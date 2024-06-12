package main

import (
	"context"
	"os"
	"stori-challenge/internal/core/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func handlerS3(ctx context.Context, event events.S3Event) {
	sessionAws := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))

	s3Service := s3.New(sessionAws)
	sqsService := sqs.New(sessionAws)

	transactionService := service.NewTransactionService(s3Service, sqsService)

	transactionService.ReadS3Files(ctx, event)
}

func main() {
	lambda.Start(handlerS3)
}
