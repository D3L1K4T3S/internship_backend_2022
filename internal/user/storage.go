package user

import (
	"context"
	"internship_bachend_2022/internal/transactions"
)

type Repository interface {
	AddFunds(ctx context.Context, id, balance string) (string, error)
	CreateUser(ctx context.Context, id, balance string) (string, error)
	CreateOrder(ctx context.Context, idOrder, idUser, idService, cost string) (string, error)
	CreateTransaction(ctx context.Context, idOrder, tp, amount string) (string, error)
	CreateDeposit(ctx context.Context, idUser, amount string) (string, error)
	GetBalance(ctx context.Context, id string) (User, error)
	ExistUserId(ctx context.Context, id string) (bool, error)
	ExistOrderId(ctx context.Context, id string) (bool, error)
	GetTransactions(ctx context.Context, id string, options Options) ([]transactions.Transactions, error)
	DeleteUser(ctx context.Context, id string) error
	RevenueRecognition(ctx context.Context, idUser, idOrder, amount string) error
}

type Options struct {
	Field, Order, List, Records string
}
