package postgresql

import (
	"context"
	"fmt"
	"internship_bachend_2022/internal/orders"
	"internship_bachend_2022/internal/user/db/postgresql"
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

func (repository *repository) GetServiceTotal(ctx context.Context, year, month string) ([]orders.ServiceTotal, error) {
	query := `SELECT services.name, SUM(orders.cost) 
			  FROM services LEFT JOIN orders ON services.id = orders.service_id 
		      WHERE orders.status = $1 and EXTRACT(YEAR FROM orders.date) = $2 and EXTRACT(MONTH FROM orders.date) = $3 
		      GROUP BY (services.name)`
	repository.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))

	var serviceTotal []orders.ServiceTotal
	rows, err := repository.client.Query(ctx, query, postgresql.Done, year, month)
	if err != nil {
		repository.logger.Info(err)
		return serviceTotal, err
	}
	for rows.Next() {
		var tmp orders.ServiceTotal
		rows.Scan(&tmp.ServiceName, &tmp.TotalAmount)
		serviceTotal = append(serviceTotal, tmp)
	}
	return serviceTotal, nil
}

func NewRepository(client postgreSQL.Client, logger *logging.Logger) orders.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
