package port

import "stori-challenge/internal/core/domain"

type ReportGenerator interface {
	GenerateMonthlyReport(transactions []domain.TransactionModel) domain.ReportDocument
}
