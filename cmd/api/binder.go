package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
	"stori-challenge/internal/adapter"
	"stori-challenge/internal/core/service"
)

func BindTransactionService() *service.TransactionService {
	// AWS services
	sessionAws := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))
	s3Service := s3.New(sessionAws)
	sqsService := sqs.New(sessionAws)
	dynamoDb := dynamodb.New(sessionAws)
	sesService := ses.New(sessionAws)

	// Internal core services
	repository := adapter.NewTransactionRepository(dynamoDb)
	reportGenerator := adapter.NewReportGenerator()
	EmailClient := adapter.NewEmailClientAdapter(sesService)

	return service.NewTransactionService(s3Service, sqsService, repository, reportGenerator, EmailClient)
}
