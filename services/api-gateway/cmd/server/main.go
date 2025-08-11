package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"service-gateway/internal/config"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct {
	Config config.Config
}

func main() {
	app := App{
		Config: config.LoadConfig(),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", app.Config.Server.Address, app.Config.Server.Port),
		Handler: nil,
	}

	go func() {
		logrus.WithFields(logrus.Fields{
			"Address": app.Config.Server.Address,
			"Port":    app.Config.Server.Port,
		}).Info("API-Gateway service started")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logrus.Warn("Server is shutting down")
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("server shutdown failed: %v", err)
	}

	logrus.Info("Server stopped gracefully")
}
