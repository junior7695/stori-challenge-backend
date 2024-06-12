package adapter

import (
	"stori-challenge/internal/core/domain"
	"testing"
	"time"
)

func TestGenerateMonthlyReport(t *testing.T) {
	transactions := []domain.TransactionModel{
		{Amount: 100.0, Date: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{Amount: -50.0, Date: time.Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC)},
		{Amount: 200.0, Date: time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC)},
		{Amount: -100.0, Date: time.Date(2023, time.February, 2, 0, 0, 0, 0, time.UTC)},
	}

	expectedReport := domain.ReportDocument{
		TotalBalance:  150.0,
		AverageCredit: 150.0,
		AverageDebit:  -75.0,
		MonthlyTransactions: map[string]int{
			time.January.String():  2,
			time.February.String(): 2,
		},
	}

	reportGenerator := NewReportGenerator()
	report := reportGenerator.GenerateMonthlyReport(transactions)

	if report.TotalBalance != expectedReport.TotalBalance {
		t.Errorf("expected total balance %v, got %v", expectedReport.TotalBalance, report.TotalBalance)
	}
	if report.AverageCredit != expectedReport.AverageCredit {
		t.Errorf("expected average credit %v, got %v", expectedReport.AverageCredit, report.AverageCredit)
	}
	if report.AverageDebit != expectedReport.AverageDebit {
		t.Errorf("expected average debit %v, got %v", expectedReport.AverageDebit, report.AverageDebit)
	}
	for month, count := range expectedReport.MonthlyTransactions {
		if report.MonthlyTransactions[month] != count {
			t.Errorf("expected %d transactions in %s, got %d", count, month, report.MonthlyTransactions[month])
		}
	}
}
