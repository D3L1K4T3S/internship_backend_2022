package db

import (
	"context"
	"internship_bachend_2022/internal/user"
	"internship_bachend_2022/pkg/client/postgreSQL"
	"internship_bachend_2022/pkg/logging"
)

type repository struct {
	client postgreSQL.Client
	logger *logging.Logger
}

func (repository *repository) Create(ctx context.Context, user user.User) (string, error) {
	query := `INSERT INTO users (id,name) VALUES ($1,$2) RETURNING id`
	repository.client.QueryRow(ctx, query, id, name)
}
func (repository *repository) FindALl(ctx context.Context) (users []user.User, err error) {

}
func (repository *repository) FindOne(ctx context.Context, id string) (user.User, error) {

}
func (repository *repository) Update(ctx context.Context, user user.User) error {

}
func (repository *repository) Delete(ctx context.Context, id string) error {

}

func NewRepository(client postgreSQL.Client, logger *logging.Logger) user.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
