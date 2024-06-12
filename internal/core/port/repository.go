package port

import "context"

type TransactionRepository interface {
	SaveTransaction(ctx context.Context)
	GetTransactions(ctx context.Context)
}
