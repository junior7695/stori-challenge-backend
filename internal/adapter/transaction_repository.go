package adapter

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"os"
	"stori-challenge/internal/core/domain"
	"strconv"
	"time"
)

type TransactionRepository struct {
	dynamoDbClient *dynamodb.DynamoDB
}

func NewTransactionRepository(dynamoDbClient *dynamodb.DynamoDB) *TransactionRepository {
	return &TransactionRepository{dynamoDbClient}
}

func (r *TransactionRepository) SaveTransaction(ctx context.Context, transaction domain.TransactionDto) {
	_, err := r.dynamoDbClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(transaction.Id),
			},
			"Date": {
				S: aws.String(transaction.Date),
			},
			"Amount": {
				N: aws.String(transaction.Amount),
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to save to DynamoDB: %v", err)
	}
}

func (r *TransactionRepository) GetTransactions(ctx context.Context) (transactions []domain.TransactionModel, err error) {
	result, err := r.dynamoDbClient.Scan(&dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
	})
	if err != nil {
		log.Fatalf("Failed to scan DynamoDB: %v", err)
	}

	for _, item := range result.Items {
		date := *item["Date"].S
		amountStr := *item["Amount"].N
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			log.Printf("Invalid amount: %v", amountStr)
			continue
		}

		dateParsed, err := time.Parse("2006-01-02", date)
		if err != nil {
			log.Printf("Invalid date: %v", date)
			continue
		}

		transactions = append(transactions, domain.TransactionModel{
			Id:     *item["Id"].S,
			Date:   dateParsed,
			Amount: amount,
		})
	}

	return transactions, nil
}
