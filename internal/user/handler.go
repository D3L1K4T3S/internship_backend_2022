package user

import (
	"github.com/julienschmidt/httprouter"
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
	router.GET(usersURL, handler.GetList)
	router.POST(usersURL, handler.CreateUser)
	router.GET(userURL, handler.GetUserByUUID)
	router.PUT(userURL, handler.UpdateUser)
	router.PATCH(userURL, handler.PartiallyUpdateUser)
	router.DELETE(userURL, handler.DeleteUser)
}

func (handler *handler) GetList(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writer.WriteHeader(200)
	writer.Write([]byte("list of users"))
}

func (handler *handler) CreateUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writer.WriteHeader(201)
	writer.Write([]byte("create user"))
}

func (handler *handler) GetUserByUUID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writer.WriteHeader(200)
	writer.Write([]byte("get user by id"))
}

func (handler *handler) UpdateUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writer.WriteHeader(204)
	writer.Write([]byte("update user"))
}

func (handler *handler) PartiallyUpdateUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writer.WriteHeader(204)
	writer.Write([]byte("partially update user"))
}

func (handler *handler) DeleteUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writer.WriteHeader(204)
	writer.Write([]byte("delete user"))
}
