package domain

type ReportDocument struct {
	TotalBalance        float64        `json:"total_balance"`
	AverageDebit        float64        `json:"average_debit"`
	AverageCredit       float64        `json:"average_credit"`
	MonthlyTransactions map[string]int `json:"monthly_transactions"`
}
