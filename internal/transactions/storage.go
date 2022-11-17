package transactions

import (
	"context"
)

type Repository interface {
	GainPerPeriod(ctx context.Context, year, month string) (string, error)
	SuccessTransaction(ctx context.Context, idUser, isService, idOrder, cost string) ([]Transactions, error)
}

type SortOptions struct {
	Field, Order string
}
