package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"rest_api/internal/config"
	"rest_api/internal/pkg/client/mongobd"
	"rest_api/internal/pkg/logging"
	"rest_api/internal/user"
	"rest_api/internal/user/db"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()
	cfg := config.GetConfig()

	cfgMongo:= cfg.MongoDB
	var mongoDBClient, err = mongobd.NewClient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username,
		cfgMongo.Password, cfgMongo.Database, cfgMongo.AuthDB)
	if err != nil {
		panic(err)
	}
	storage := db.NewStorage(mongoDBClient, cfgMongo.Collection, logger)

	user1:= user.User{
		ID:           "",
		Email:        "istefa92@gmail.com",
		Username:     "Nikita",
		PasswordHash: "12345",
	}

	userID, err := storage.Create(context.Background(), user1)
	if err != nil {
		panic(err)
	}
	logger.Info(userID)

	logger.Info("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)



	start(router, cfg)

}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")
	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		sokcetPath := path.Join(appDir, "app.sokc")
		logger.Debugf("socket path %s", sokcetPath)

		logger.Info("listen unix soket")
		listener, listenErr = net.Listen("unix", sokcetPath)
		logger.Info("server is listening unix sokced %s", sokcetPath)
	} else {
		logger.Info("listen tcp sokcet")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
		logger.Info("server is listening port %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Fatalln(server.Serve(listener))
}
