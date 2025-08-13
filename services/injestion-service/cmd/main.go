package main

import (
	"fmt"
	"injestion-service/internal/config"
	mqtt_broker "injestion-service/internal/mqtt"
	"injestion-service/internal/rabbitmq"
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

	mqttQuit := make(chan bool)
	mqttMessages := make(chan []byte, 100)

	go mqtt_broker.RunMqttClient(
		fmt.Sprintf("tcp://%s:%d", app.Config.MQTT.Address, app.Config.MQTT.Port),
		app.Config.MQTT.ClientID,
		app.Config.MQTT.Topic,
		mqttMessages,
		mqttQuit,
	)

	go rabbitmq.Publisher(app.Config.RabbitMQ.URL, mqttMessages)

	logrus.Info("Service Started Successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	mqttQuit <- true
	time.Sleep(3 * time.Second)
	logrus.Info("Service stopped gracefully")
}
