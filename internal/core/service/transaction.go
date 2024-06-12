package service

import (
	"context"
	"encoding/csv"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
	"stori-challenge/internal/core/domain"
	"stori-challenge/internal/core/port"
	"strings"
)

type TransactionService struct {
	s3Service       *s3.S3
	sqsService      *sqs.SQS
	repository      port.TransactionRepository
	reportGenerator port.ReportGenerator
	emailClient     port.EmailClient
}

func NewTransactionService(
	s3 *s3.S3,
	sqs *sqs.SQS,
	repository port.TransactionRepository,
	reportGenerator port.ReportGenerator,
	emailClient port.EmailClient,
) *TransactionService {
	return &TransactionService{
		s3Service:       s3,
		sqsService:      sqs,
		repository:      repository,
		reportGenerator: reportGenerator,
		emailClient:     emailClient,
	}
}

func (srv *TransactionService) ReadS3Files(ctx context.Context, event events.S3Event) {
	for _, record := range event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		if key != *aws.String(os.Getenv("BUCKET_CSV")) {
			continue
		}

		obj, err := srv.s3Service.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			log.Fatalf("Failed to get object from S3: %v", err)
		}
		defer obj.Body.Close()

		r := csv.NewReader(obj.Body)
		records, err := r.ReadAll()
		if err != nil {
			log.Fatalf("Failed to read CSV: %v", err)
		}

		for _, record := range records {
			messageBody := record[0] + "," + record[1] + "," + record[2]
			_, err := srv.sqsService.SendMessage(&sqs.SendMessageInput{
				MessageBody: aws.String(messageBody),
				QueueUrl:    aws.String(os.Getenv("SQS_QUEUE_URL")),
			})
			if err != nil {
				log.Fatalf("Failed to send message to SQS: %v", err)
			}
		}
	}
}

func (srv *TransactionService) SaveTransactions(ctx context.Context, sqsEvent events.SQSEvent) {
	for _, message := range sqsEvent.Records {
		body := message.Body
		parts := strings.Split(body, ",")
		if len(parts) != 3 {
			log.Fatalf("Invalid message format: %s", body)
		}

		transactionDto := domain.TransactionDto{
			Id:     parts[0],
			Date:   parts[1],
			Amount: parts[2],
		}

		srv.repository.SaveTransaction(ctx, transactionDto)
	}
}

func (srv *TransactionService) SendMonthlyReport(ctx context.Context) {
	transactions, err := srv.repository.GetTransactions(ctx)
	if err != nil {
		log.Fatalf("Failed to get transactions: %v", err)
	}

	reportDocument := srv.reportGenerator.GenerateMonthlyReport(transactions)

	srv.emailClient.SendEmailReport(ctx, reportDocument)
}
