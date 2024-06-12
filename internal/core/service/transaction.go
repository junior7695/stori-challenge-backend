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
)

type TransactionService struct {
	s3Service  *s3.S3
	sqsService *sqs.SQS
}

func NewTransactionService(s3 *s3.S3, sqs *sqs.SQS) *TransactionService {
	return &TransactionService{
		s3Service:  s3,
		sqsService: sqs,
	}
}

func (srv *TransactionService) ReadS3Files(ctx context.Context, event events.S3Event) {
	for _, record := range event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

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

func (srv *TransactionService) SaveTransactions() {

}
