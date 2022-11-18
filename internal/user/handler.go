package user

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"internship_bachend_2022/internal/apperror"
	"internship_bachend_2022/internal/handlers"
	"internship_bachend_2022/internal/orders"
	"internship_bachend_2022/pkg/logging"
	"net/http"
)

const (
	userURL         = "/user/"
	usersURL        = "/users/"
	orderURL        = "/order/"
	transactionsURL = "/trans/"
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
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(handler.GetBalance))
	//ПРоверить через метод PUT (был POST)
	router.HandlerFunc(http.MethodPut, usersURL, apperror.Middleware(handler.AddingFunds))
	router.HandlerFunc(http.MethodPost, orderURL, apperror.Middleware(handler.CreateOrder))
	router.HandlerFunc(http.MethodPatch, userURL, apperror.Middleware(handler.RevenueRecognition))
	router.HandlerFunc(http.MethodDelete, userURL, apperror.Middleware(handler.DeleteUser))
	router.HandlerFunc(http.MethodGet, transactionsURL, apperror.Middleware(handler.GetTransactions))
}

// GetBalance godoc
// @description Get user balance by id
// @param       id
// @Summary GetBalance
// @Tags User
// @Success 200
// @Failure 404
// @Router getBalance[get]
func (handler *handler) GetBalance(writer http.ResponseWriter, request *http.Request) error {

	var id string
	id = request.URL.Query().Get("id")
	if id == "" {
		return apperror.IncorrectRequest
	}

	usr, err := handler.repository.GetBalance(context.TODO(), id)
	if err != nil {
		handler.logger.Info(err)
		return apperror.NotFoundUser
	}

	jsonUser, err := json.Marshal(usr)
	if err != nil {
		handler.logger.Info(err)
	}

	_, err = writer.Write(jsonUser)
	if err != nil {
		handler.logger.Info(err)
	}

	return nil
}

// AddingFunds godoc
// @Summary AddingFunds
// @Description Add funds user by id
// @Param       id
// @Tags User
// @Success 200
// @Failure 404
// @Router addingFunds [post]
func (handler *handler) AddingFunds(writer http.ResponseWriter, request *http.Request) error {
	writer.Header().Set("Content-Type", "application/json")

	var user User
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&user)
	if err != nil {
		handler.logger.Info(err)
		return err
	}
	exist, err := handler.repository.ExistUserId(context.TODO(), user.Id)
	if err != nil {
		handler.logger.Info(err)
		return err
	}

	if exist {
		_, err := handler.repository.AddFunds(context.TODO(), user.Id, user.Balance)
		if err != nil {
			handler.logger.Info(err)
			return err
		}
		return nil
	} else {
		_, err := handler.repository.CreateUser(context.TODO(), user.Id, user.Balance)
		if err != nil {
			handler.logger.Info(err)
			return err
		}
		return nil
	}
}

// CreateOrder godoc
// @Summary CreateOrder
// @Description Create order by user_id, service_id, amount
// @Param       user_id, service_id, amount
// @Tags User
// @Success 200
// @Failure 404
// @Router createOrder [post]
func (handler *handler) CreateOrder(writer http.ResponseWriter, request *http.Request) error {

	decoder := json.NewDecoder(request.Body)
	var order orders.Orders
	err := decoder.Decode(&order)
	if err != nil {
		handler.logger.Info(err)
		return err
	}

	exist, err := handler.repository.ExistOrderId(context.TODO(), order.Id)
	if err != nil {
		handler.logger.Info(err)
		return err
	}

	if exist {
		return apperror.OrderExist
	} else {
		_, err := handler.repository.CreateOrder(context.TODO(), order.Id, order.UserId, order.ServiceId, order.Cost)
		if err != nil {
			handler.logger.Info(err)
			return err
		}
	}

	return nil
}

// RevenueRecognition godoc
// @Summary RevenueRecognition
// @Description Debiting funds from a temporary account to the company's bank
// @Param       user_id, service_id, order_id, amount
// @Tags User
// @Success 200
// @Failure 404
// @Router revenueRecognition[patch]
func (handler *handler) RevenueRecognition(writer http.ResponseWriter, request *http.Request) error {

	var order orders.Orders
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&order)
	if err != nil {
		handler.logger.Info(err)
		return err
	}

	err = handler.repository.RevenueRecognition(context.TODO(), order.UserId, order.Id, order.Cost)
	if err != nil {
		handler.logger.Info(err)
		return err
	}

	return nil
}

// GetTransactions godoc
// @Summary GetTransactions
// @Description get transactions by user id
// @Param       id, sort, list
// @Tags User
// @Success 200
// @Failure 404
// @Router getTransactions[get]
func (handler *handler) GetTransactions(writer http.ResponseWriter, request *http.Request) error {
	//var id, sort, list string
	//id = request.URL.Query().Get("id")
	//sort = request.URL.Query().Get("sort")
	//list = request.URL.Query().Get("list")
	//
	//if id == "" {
	//	return apperror.IncorrectRequest
	//}

	return nil
}

// DeleteUser godoc
// @Summary DeleteUser
// @Description Delete user from db
// @Param       user_id,
// @Tags User
// @Success 200
// @Failure 404
// @Router revenueRecognition[delete]
func (handler *handler) DeleteUser(writer http.ResponseWriter, request *http.Request) error {

	decoder := json.NewDecoder(request.Body)
	var user User
	err := decoder.Decode(&user)
	if err != nil {
		handler.logger.Info(err)
		return err
	}
	err = handler.repository.DeleteUser(context.TODO(), user.Id)
	if err != nil {
		handler.logger.Info(err)
		return err
	}

	return nil
}
