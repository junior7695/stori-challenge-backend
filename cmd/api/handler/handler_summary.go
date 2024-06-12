package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ses"
)

func handlerSummaryTransaction(ctx context.Context) {
	sessionAws := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))

	dynamoService := dynamodb.New(sessionAws)
	sesService := ses.New(sessionAws)

	result, err := dynamoService.Scan(&dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
	})
	if err != nil {
		log.Fatalf("Failed to scan DynamoDB: %v", err)
	}

	var totalBalance float64
	var creditCount, debitCount int
	monthlyCredits := make(map[string][]float64)
	monthlyDebits := make(map[string][]float64)

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

		month := dateParsed.Format("2006-01")

		if amount > 0 {
			monthlyCredits[month] = append(monthlyCredits[month], amount)
			creditCount++
		} else {
			monthlyDebits[month] = append(monthlyDebits[month], amount)
			debitCount++
		}

		totalBalance += amount
	}

	emailBody := formatEmailBody(totalBalance, creditCount, debitCount, monthlyCredits, monthlyDebits)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(os.Getenv("SES_RECIPIENT")),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(emailBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("Monthly Transaction Summary"),
			},
		},
		Source: aws.String(os.Getenv("SES_SENDER")),
	}

	_, err = sesService.SendEmail(input)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}
}

func formatEmailBody(totalBalance float64, creditCount, debitCount int, monthlyCredits, monthlyDebits map[string][]float64) string {
	monthlySummary := "<ul>"
	for month, credits := range monthlyCredits {
		creditSum := 0.0
		for _, credit := range credits {
			creditSum += credit
		}
		avgCredit := creditSum / float64(len(credits))

		debits := monthlyDebits[month]
		debitSum := 0.0
		for _, debit := range debits {
			debitSum += debit
		}
		avgDebit := debitSum / float64(len(debits))

		monthlySummary += fmt.Sprintf("<li>%s: Avg Credit: %.2f, Avg Debit: %.2f</li>", month, avgCredit, avgDebit)
	}
	monthlySummary += "</ul>"

	return fmt.Sprintf(`
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; }
				h1 { color: #333; }
				p { color: #666; }
			</style>
		</head>
		<body>
			<h1>Transaction Summary</h1>
			<p>Total Balance: %.2f</p>
			<p>Number of Credit Transactions: %d</p>
			<p>Number of Debit Transactions: %d</p>
			<h2>Monthly Summary</h2>
			%s
		</body>
		</html>
	`, totalBalance, creditCount, debitCount, monthlySummary)
}

func main() {
	lambda.Start(handlerSummaryTransaction)
}
