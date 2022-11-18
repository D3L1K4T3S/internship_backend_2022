package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/pgconn"
	"internship_bachend_2022/internal/apperror"
	"internship_bachend_2022/internal/transactions"
	"internship_bachend_2022/internal/user"
	"internship_bachend_2022/pkg/client/postgreSQL"
	"internship_bachend_2022/pkg/logging"
	"strconv"
	"strings"
)

const (
	Debiting       = "Debiting"
	Refund         = "Refund"
	Done           = "Done"
	NotEnoughFunds = "Not enough funds"
)

type repository struct {
	client postgreSQL.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (repository *repository) AddFunds(ctx context.Context, id, balance string) (string, error) {
	query := `UPDATE users SET balance = balance + $1 WHERE id = $2`
	repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))
	_ = repository.client.QueryRow(ctx, query, balance, id)

	_, err := repository.CreateDeposit(ctx, id, balance)
	if err != nil {
		repository.logger.Info(err)
		return "", err
	}

	return id, nil

}

func (repository *repository) CreateUser(ctx context.Context, id, balance string) (string, error) {

	query := `INSERT INTO users (id,balance,reserved) VALUES ($1,$2,$3) RETURNING id`
	repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))

	if err := repository.client.QueryRow(ctx, query, id, balance, "0").Scan(&id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			repository.logger.Infof("Detail: %s, Code: %s, Where: %s", pgErr.Detail, pgErr.Code, pgErr.Where)
			return "", err
		}
		return "", err
	}

	_, err := repository.CreateDeposit(ctx, id, balance)
	if err != nil {
		repository.logger.Info(err)
		return "", err
	}

	return id, nil
}

func (repository *repository) ExistUserId(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT * FROM users WHERE id = $1)`
	repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))
	var exist bool
	if err := repository.client.QueryRow(ctx, query, id).Scan(&exist); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			repository.logger.Infof("Detail: %s, Code: %s, Where: %s", pgErr.Detail, pgErr.Code, pgErr.Where)
			return exist, err
		}
		return exist, err
	}
	return exist, nil
}

func (repository *repository) ExistOrderId(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT * FROM orders WHERE id = $1)`
	repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))
	var exist bool
	if err := repository.client.QueryRow(ctx, query, id).Scan(&exist); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			repository.logger.Infof("Detail: %s, Code: %s, Where: %s", pgErr.Detail, pgErr.Code, pgErr.Where)
			return exist, err
		}
		return exist, err
	}
	return exist, nil
}

func (repository *repository) CreateDeposit(ctx context.Context, idUser, amount string) (string, error) {
	query := `INSERT INTO deposits (user_id, amount) VALUES ($1,$2) RETURNING id`
	repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))

	var id string
	if err := repository.client.QueryRow(ctx, query, idUser, amount).Scan(&id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			repository.logger.Infof("Detail: %s, Code: %s, Where: %s", pgErr.Detail, pgErr.Code, pgErr.Where)
			return "", err
		}
		return "", err
	}
	return id, nil
}

func (repository *repository) CreateTransaction(ctx context.Context, idOrder, tp, amount string) (string, error) {
	query := `INSERT INTO transactions (order_id, type, amount) VALUES ($1,$2,$3) RETURNING id`
	repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))

	var id string
	if err := repository.client.QueryRow(ctx, query, idOrder, tp, amount).Scan(&id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			repository.logger.Infof("Detail: %s, Code: %s, Where: %s", pgErr.Detail, pgErr.Code, pgErr.Where)
			return "", err
		}
		return "", err
	}
	return id, nil
}

func (repository *repository) CreateOrder(ctx context.Context, idOrder, idUser, idService, cost string) (string, error) {

	var balance, cst float64
	var bal, query string
	cst, _ = strconv.ParseFloat(cost, 64)

	getBalance := `SELECT balance FROM users WHERE id = $1`
	repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(getBalance)))

	if err := repository.client.QueryRow(ctx, getBalance, idUser).Scan(&bal); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			repository.logger.Infof("Detail: %s, Code: %s, Where: %s", pgErr.Detail, pgErr.Code, pgErr.Where)
		}
	}

	balance, _ = strconv.ParseFloat(bal, 64)

	if balance-cst < 0 {

		query = `INSERT INTO orders (id, user_id, service_id, status, cost) VALUES ($1,$2,$3,$4,$5) RETURNING id`
		repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))
		if err := repository.client.QueryRow(ctx, query, idOrder, idUser, idService, NotEnoughFunds, cost).Scan(&idOrder); err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				repository.logger.Infof("Detail: %s, Code: %s, Where: %s", pgErr.Detail, pgErr.Code, pgErr.Where)
				return "", nil
			}
			return "", err
		}

		_, err := repository.CreateTransaction(ctx, idOrder, Debiting, cost)
		if err != nil {
			repository.logger.Info(err)
			return "", err
		}

		repository.logger.Infof("not enough money user id = %v. Refund amount", idUser)
		_, err = repository.CreateTransaction(ctx, idOrder, Refund, cost)
		if err != nil {
			repository.logger.Info(err)
			return "", err
		}

		return idOrder, nil
	} else {
		query = `UPDATE users SET balance = balance - $1, reserved = reserved + $1 WHERE id = $2`
		repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))
		_ = repository.client.QueryRow(ctx, query, cst, idUser)

		query = `INSERT INTO orders (id, user_id, service_id, cost) VALUES ($1,$2,$3,$4) RETURNING id`
		repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))
		if err := repository.client.QueryRow(ctx, query, idOrder, idUser, idService, cost).Scan(&idOrder); err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				repository.logger.Infof("Detail: %s, Code: %s, Where: %s", pgErr.Detail, pgErr.Code, pgErr.Where)
				return "", nil
			}
			return "", err
		}

		_, err := repository.CreateTransaction(ctx, idOrder, Debiting, cost)
		if err != nil {
			repository.logger.Info(err)
			return "", err
		}

		return idOrder, nil
	}
}

func (repository *repository) GetBalance(ctx context.Context, id string) (user.User, error) {
	query := `SELECT balance, reserved FROM users WHERE id = $1`
	repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))

	usr := user.User{Id: id, Balance: "0", Reserved: "0"}
	if err := repository.client.QueryRow(ctx, query, id).Scan(&usr.Balance, &usr.Reserved); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			repository.logger.Infof("Detail: %s, Code: %s, Where: %s", pgErr.Detail, pgErr.Code, pgErr.Where)
			return usr, err
		}
		return usr, err
	}
	return usr, nil
}

func (repository *repository) RevenueRecognition(ctx context.Context, idUser, idOrder, amount string) error {

	getReserved := `SELECT reserved FROM users WHERE idUser = $1`
	var res float64
	if err := repository.client.QueryRow(ctx, getReserved, idUser).Scan(&res); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			repository.logger.Infof("Detail: %s, Code: %s, Where: %s", pgErr.Detail, pgErr.Code, pgErr.Where)
			return err
		}
		return err
	}

	amt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		repository.logger.Info(err)
		return err
	}

	if res-amt < 0 {
		return apperror.IncorrectRequest
	} else {
		getTrans := `SELECT pass FROM transactions WHERE order_id = $1 and type = $2`
		repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(getTrans)))
		var pass bool
		if err := repository.client.QueryRow(ctx, getTrans, idOrder, Debiting).Scan(&pass); err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				repository.logger.Infof("Detail: %s, Code: %s, Where: %s", pgErr.Detail, pgErr.Code, pgErr.Where)
				return err
			}
			return err
		}

		if pass {
			changeStatus := `UPDATE orders SET status = $1 WHERE idUser = $2`
			repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(changeStatus)))
			_ = repository.client.QueryRow(ctx, changeStatus, Done, idOrder)
			query := `UPDATE users SET reserved = reserved - $1 WHERE idUser = $2 RETURNING idUser`
			repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))
			_ = repository.client.QueryRow(ctx, query, amount, idUser)
		} else {
			return apperror.TransactionNotPass
		}
		return nil
	}
}

func (repository *repository) GetTransactions(ctx context.Context, id string) ([]transactions.Transactions, error) {
	//TODO: Реализовать пагинацию и сортировку
	//TODO: Изменить запрос на вывод транзакция конкретного пользователя

	var trns []transactions.Transactions

	query := `SELECT 
    	transactions.type, 
    	transactions.time_trans, 
    	transactions.amount, 
    	transactions.comment
	FROM orders LEFT JOIN transactions ON orders.transaction_id = transactions.id 
	WHERE orders.user_id = $1;`

	repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))

	rows, err := repository.client.Query(ctx, query, id)
	if err != nil {
		repository.logger.Fatal(err)
		return trns, err
	}

	for rows.Next() {
		var tmp transactions.Transactions
		err := rows.Scan(&tmp.Type, &tmp.TimeTrans, &tmp.Amount, &tmp.Comment)
		if err != nil {
			repository.logger.Fatal(err)
			return trns, err
		}
		trns = append(trns, tmp)
	}

	return trns, nil
}

func (repository *repository) DeleteUser(ctx context.Context, id string) error {
	return nil
}

func NewRepository(client postgreSQL.Client, logger *logging.Logger) user.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
