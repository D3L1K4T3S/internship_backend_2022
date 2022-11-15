package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"internship_bachend_2022/internal/config"
	"internship_bachend_2022/internal/user"
	"internship_bachend_2022/pkg/logging"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

const socket = "sock"
const applicationSocket = "application.sock"

func main() {

	logger := logging.GetLogger()
	logger.Info("create router")

	router := httprouter.New()
	logger.Info("register user handler")

	cfg := config.GetConfig()

	handler := user.NewHandler(logger)
	handler.Register(router)

	startServer(router, cfg)
}

func startServer(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listerError error

	if cfg.Listen.Type == socket {
		logger.Info("detect socket")
		applicationDirectory, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("create socket")
		socketPath := path.Join(applicationDirectory, applicationSocket)
		logger.Debugf("socket path: %s", socketPath)

		logger.Info("server in listening unix socket")
		listener, listerError = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket %s", socketPath)
	} else {
		logger.Info("detect tcp")
		listener, listerError = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listerError != nil {
		logger.Fatal(listerError)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
