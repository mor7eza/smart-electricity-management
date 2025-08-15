package main

import (
	"injestion-service/internal/config"
	mqtt_broker "injestion-service/internal/mqtt"
	redis_db "injestion-service/internal/redis"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct {
	RedisService redis_db.RedisService
	MqttService  mqtt_broker.MqttService
}

func main() {
	app := App{}
	app.RedisService = *redis_db.NewService()

	app.MqttService = *mqtt_broker.NewService(app.RedisService.Client)

	go app.MqttService.Run()

	logrus.Info("Service Started Successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	app.MqttService.Client.Disconnect(250)
	app.RedisService.Client.Close()
	time.Sleep(3 * time.Second)
	logrus.Info("Service stopped gracefully")
}

func init() {
	config.LoadConfig()
}
