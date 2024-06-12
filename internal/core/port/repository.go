package port

import (
	"context"
	"stori-challenge/internal/core/domain"
)

type TransactionRepository interface {
	SaveTransaction(ctx context.Context, transaction domain.TransactionDto)
	GetTransactions(ctx context.Context) ([]domain.TransactionModel, error)
}
