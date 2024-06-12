package adapter

import (
	"stori-challenge/internal/core/domain"
)

type ReportGeneratorAdapter struct{}

func NewReportGenerator() *ReportGeneratorAdapter {
	return &ReportGeneratorAdapter{}
}

func (adapter *ReportGeneratorAdapter) GenerateMonthlyReport(transactions []domain.TransactionModel) domain.ReportDocument {
	reportDocument := domain.ReportDocument{
		TotalBalance:        0,
		AverageCredit:       0,
		AverageDebit:        0,
		MonthlyTransactions: make(map[string]int),
	}

	quantityCredits := 0
	quantityDebits := 0
	totalCredit := 0.0
	totalDebit := 0.0

	for _, transaction := range transactions {
		month := transaction.Date.Month()
		reportDocument.MonthlyTransactions[month.String()]++
		if transaction.Amount > 0 {
			quantityCredits++
			totalCredit += transaction.Amount
		} else {
			quantityDebits++
			totalDebit += transaction.Amount
		}

		reportDocument.TotalBalance += transaction.Amount
	}

	if quantityDebits > 0 {
		reportDocument.AverageDebit = totalDebit / float64(quantityDebits)
	}
	if quantityCredits > 0 {
		reportDocument.AverageCredit = totalCredit / float64(quantityCredits)
	}

	return reportDocument
}
