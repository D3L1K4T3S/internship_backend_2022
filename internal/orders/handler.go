package orders

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"internship_bachend_2022/internal/apperror"
	"internship_bachend_2022/internal/handlers"
	"internship_bachend_2022/pkg/logging"
	"net/http"
)

const (
	totalURL = "/total/"
)

type handler struct {
	logger     *logging.Logger
	repository Repository
}

func NewHandler(repository Repository, logger *logging.Logger) handlers.Handler {
	return &handler{
		repository: repository,
		logger:     logger,
	}
}

func (handler *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, totalURL, apperror.Middleware(handler.GetRevenue))
}

func (handler *handler) GetRevenue(writer http.ResponseWriter, request *http.Request) error {
	writer.Header().Set("Content-Type", "text/csvFile")

	year := request.URL.Query().Get("year")
	month := request.URL.Query().Get("month")

	revenue, err := handler.repository.GetServiceTotal(context.TODO(), year, month)
	if err != nil {
		handler.logger.Info(err)
		return err
	}

	pathFile := "files/total.csv"
	var header []string
	header = append(header, "Service name", "Total amount")
	err = CreateData(pathFile, header, revenue)
	if err != nil {
		handler.logger.Info(err)
		return err
	}
	_, err = writer.Write([]byte(fmt.Sprintf("File saved in %v", "internship_backend_2022/"+pathFile)))
	if err != nil {
		handler.logger.Info(err)
		return err
	}
	return nil
}
