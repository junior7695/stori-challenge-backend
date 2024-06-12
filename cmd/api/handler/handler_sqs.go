package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func handlerSqs(ctx context.Context, sqsEvent events.SQSEvent) {
	sessionAws := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))
	dynamoService := dynamodb.New(sessionAws)

	for _, message := range sqsEvent.Records {
		body := message.Body
		parts := strings.Split(body, ",")
		if len(parts) != 3 {
			log.Printf("Invalid message format: %s", body)
			continue
		}

		id := parts[0]
		date := parts[1]
		amount := parts[2]

		_, err := dynamoService.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
			Item: map[string]*dynamodb.AttributeValue{
				"Id": {
					S: aws.String(id),
				},
				"Date": {
					S: aws.String(date),
				},
				"Amount": {
					N: aws.String(amount),
				},
			},
		})
		if err != nil {
			log.Fatalf("Failed to save to DynamoDB: %v", err)
		}
	}
}

func main() {
	lambda.Start(handlerSqs)
}
