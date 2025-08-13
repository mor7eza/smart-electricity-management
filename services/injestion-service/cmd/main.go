package main

import (
	"fmt"
	"injestion-service/internal/config"
	mqtt_broker "injestion-service/internal/mqtt"
	"os"
	"os/signal"
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

	mqqtQuit := make(chan bool)
	go mqtt_broker.RunMqttClient(
		fmt.Sprintf("tcp://%s:%d", app.Config.MQTT.Address, app.Config.MQTT.Port),
		app.Config.MQTT.ClientID,
		app.Config.MQTT.Topic,
		mqqtQuit,
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	mqqtQuit <- true
	time.Sleep(3 * time.Second)
	logrus.Info("Service stopped gracefully")
}
