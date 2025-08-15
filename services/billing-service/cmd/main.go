package main

import (
	"billing-service/internal/config"
	"billing-service/internal/rabbitmq"
	redis_db "billing-service/internal/redis"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct {
	RedisService    *redis_db.RedisService
	RabbitMQService *rabbitmq.RabbitMQService
}

func main() {
	var app App

	app.RedisService = redis_db.NewService()
	app.RabbitMQService = rabbitmq.NewService()

	go app.Process()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	app.RedisService.Client.Close()
	time.Sleep(3 * time.Second)
	logrus.Info("Service stopped gracefully")
}

func init() {
	config.LoadConfig()
}
