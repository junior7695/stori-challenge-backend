package adapter

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"html/template"
	"log"
	"os"
	"stori-challenge/internal/core/domain"
)

const (
	CharSet = "UTF-8"
	Subject = "Transaction Summary"
)

var htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Transaction Summary</title>
</head>
<body>
    <div style="text-align: center;">
        <img src="{{.LogoURL}}" alt="Logo" style="width: 150px;"/>
    </div>
    <h1>Transaction Summary</h1>
    <p>Total balance: {{.TotalBalance}}</p>
    <p>Average debit amount: {{.AverageDebit}}</p>
    <p>Average credit amount: {{.AverageCredit}}</p>
    {{range $month, $count := .TransactionsByMonth}}
    <p>Number of transactions in {{$month}}: {{$count}}</p>
    {{end}}
</body>
</html>
`

type EmailClientAdapter struct {
	sesClient *ses.SES
}

func NewEmailClientAdapter(sesClient *ses.SES) *EmailClientAdapter {
	return &EmailClientAdapter{sesClient}
}

func (adapter *EmailClientAdapter) SendEmailReport(ctx context.Context, reportDocument domain.ReportDocument) {
	body, err := formatEmailBody(reportDocument)
	if err != nil {
		log.Fatalf("failed to format email body: %v", err)
	}

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(os.Getenv("SES_RECIPIENT")),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(body.String()),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(os.Getenv("SES_SENDER")),
	}

	_, err = adapter.sesClient.SendEmailWithContext(ctx, input)

	if err != nil {
		log.Fatalf("failed to send email: %v", err)
	}
}

func formatEmailBody(reportDocument domain.ReportDocument) (body bytes.Buffer, err error) {
	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		return body, fmt.Errorf("failed to parse email template: %w", err)
	}

	data := struct {
		LogoURL             string
		TotalBalance        float64
		AverageDebit        float64
		AverageCredit       float64
		TransactionsByMonth map[string]int
	}{
		LogoURL:             *aws.String(os.Getenv("LOGO_URL")),
		TotalBalance:        reportDocument.TotalBalance,
		AverageDebit:        reportDocument.AverageDebit,
		AverageCredit:       reportDocument.AverageCredit,
		TransactionsByMonth: reportDocument.MonthlyTransactions,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return body, fmt.Errorf("failed to execute email template: %w", err)
	}

	return body, nil
}
