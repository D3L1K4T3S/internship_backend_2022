package postgresql

import (
	"context"
	"internship_bachend_2022/internal/user"
	"internship_bachend_2022/pkg/client/postgreSQL"
	"internship_bachend_2022/pkg/logging"
	"strings"
)

type repository struct {
	client postgreSQL.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (repository *repository) GetTransactions(ctx context.Context, id string) (string, error) {

}
func (repository *repository) SuccessTransaction(ctx context.Context, idUser, isService, idOrder, cost string) error {

}

func NewRepository(client postgreSQL.Client, logger *logging.Logger) user.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
