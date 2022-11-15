package user

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"internship_bachend_2022/internal/apperror"
	"internship_bachend_2022/internal/handlers"
	"internship_bachend_2022/pkg/logging"
	"net/http"
)

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)

type handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
	}
}

func (handler *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(handler.GetList))
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(handler.CreateUser))
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(handler.GetUserByUUID))
	router.HandlerFunc(http.MethodPut, userURL, apperror.Middleware(handler.UpdateUser))
	router.HandlerFunc(http.MethodPatch, userURL, apperror.Middleware(handler.PartiallyUpdateUser))
	router.HandlerFunc(http.MethodDelete, userURL, apperror.Middleware(handler.DeleteUser))
}

func (handler *handler) GetList(writer http.ResponseWriter, request *http.Request) error {
	//writer.WriteHeader(200)
	//writer.Write([]byte("list of users"))

	return apperror.ErrorNotFound
}

func (handler *handler) CreateUser(writer http.ResponseWriter, request *http.Request) error {
	//writer.WriteHeader(201)
	//writer.Write([]byte("create user"))

	return fmt.Errorf("this is test")
}

func (handler *handler) GetUserByUUID(writer http.ResponseWriter, request *http.Request) error {
	//writer.WriteHeader(200)
	//writer.Write([]byte("get user by id"))

	return apperror.NewApplicationError(nil, "test", "test", "1234")
}

func (handler *handler) UpdateUser(writer http.ResponseWriter, request *http.Request) error {
	writer.WriteHeader(204)
	writer.Write([]byte("update user"))

	return nil
}

func (handler *handler) PartiallyUpdateUser(writer http.ResponseWriter, request *http.Request) error {
	writer.WriteHeader(204)
	writer.Write([]byte("partially update user"))

	return nil
}

func (handler *handler) DeleteUser(writer http.ResponseWriter, request *http.Request) error {
	writer.WriteHeader(204)
	writer.Write([]byte("delete user"))

	return nil
}
