package domain

import "time"

type TransactionDto struct {
	Id     string `json:"id"`
	Date   string `json:"date"`
	Amount string `json:"amount"`
}

type TransactionModel struct {
	Id     string    `json:"id"`
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}
