package adapter

import "github.com/aws/aws-sdk-go/service/dynamodb"

type TransactionRepository struct {
	dynamoDbClient *dynamodb.DynamoDB
}

func NewTransactionRepository(dynamoDbClient *dynamodb.DynamoDB) *TransactionRepository {
	return &TransactionRepository{dynamoDbClient}
}

func (r *TransactionRepository) SaveTransaction() {

}
