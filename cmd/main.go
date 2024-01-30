package main

import (
	"context"
	"device-worker-service/config"
	"device-worker-service/httpserver"
	"device-worker-service/internal/handler"
	"device-worker-service/internal/service"
	"device-worker-service/pkg/logger"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var env = "config/cloud_config.yml"

//	@title	Device worker service
//
// @host    192.168.0.105:8087
func main() {

	zapLogger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	appConfig, err := config.NewConfig(env)
	if err != nil {
		zapLogger.Error("error init config", zap.Error(err))
		return
	}

	ctx, contextCancel := context.WithCancel(context.Background())

	services := service.NewServices(ctx, appConfig, zapLogger)
	handler := handler.NewHandler(services, zapLogger)
	server := httpserver.NewServer(handler.InitRoutes())

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			zapLogger.Error("error server start", zap.Error(err))
			return
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)
	<-stop
	contextCancel()

	start := time.Now()
	for i := 5; i > 0; i-- {
		fmt.Println(i)

		if time.Since(start) > 20 {
			break
		}

		time.Sleep(1 * time.Second)
	}

	zapLogger.Info("application closed")
}
