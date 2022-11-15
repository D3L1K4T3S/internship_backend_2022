package main

import (
	"github.com/julienschmidt/httprouter"
	"internship_bachend_2022/internal/user"
	"internship_bachend_2022/pkg/logging"
	"net"
	"net/http"
	"time"
)

func main() {

	logger := logging.GetLogger()
	logger.Info("create router")

	router := httprouter.New()
	logger.Info("register user handler")

	handler := user.NewHandler(logger)
	handler.Register(router)

	startServer(router)
}

func startServer(router *httprouter.Router) {
	logger := logging.GetLogger()
	logger.Info("start application")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("server is listening port 8080")
	logger.Fatal(server.Serve(listener))
}
