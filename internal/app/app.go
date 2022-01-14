package app

import (
	"fmt"
	"github.com/dkischenko/chat/internal/config"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"time"
)

func Run(router *httprouter.Router, logger *logger.Logger, config *config.Config) {
	logger.Info("start application")
	logger.Info("listen TCP")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Listen.Ip, config.Listen.Port))

	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Infof("server listening address %s:%s", config.Listen.Ip, config.Listen.Port)
	log.Fatal(server.Serve(listener))
}
