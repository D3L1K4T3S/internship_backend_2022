package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "internship_bachend_2022/docs"
	"internship_bachend_2022/internal/config"
	"internship_bachend_2022/internal/orders"
	ord "internship_bachend_2022/internal/orders/db/postgresql"
	"internship_bachend_2022/internal/user"
	usr "internship_bachend_2022/internal/user/db/postgresql"
	"internship_bachend_2022/pkg/client/postgreSQL"
	"internship_bachend_2022/pkg/logging"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func main() {

	logger := logging.GetLogger()
	logger.Info("create router")

	router := httprouter.New()

	logger.Info("swagger docs initializing")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	cfg := config.GetConfig()

	postgreSQLClient, err := postgreSQL.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatal(err)
	}

	repositoryUser := usr.NewRepository(postgreSQLClient, logger)
	handlerUser := user.NewHandler(repositoryUser, logger)
	handlerUser.Register(router)

	repositoryOrders := ord.NewRepository(postgreSQLClient, logger)
	handlerOrders := orders.NewHandler(repositoryOrders, logger)
	handlerOrders.Register(router)

	startServer(router, cfg)
}

func startServer(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listerError error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect socket")
		applicationDirectory, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("create socket")
		socketPath := path.Join(applicationDirectory, cfg.Listen.SocketFile)
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
