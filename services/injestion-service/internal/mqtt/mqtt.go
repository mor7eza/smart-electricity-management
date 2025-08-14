package mqtt_broker

import (
	"context"
	redis_db "injestion-service/internal/redis"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	keepAlive     = 30 * time.Second
	pingTimeout   = 1 * time.Second
	retryInterval = 10 * time.Second
)

func RunMqttClient(brokerURL, clientID, topic string, rdb *redis.Client, quit chan bool) {
	var ctx = context.Background()

	opts := mqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(clientID).
		SetDefaultPublishHandler(func(c mqtt.Client, m mqtt.Message) {
			redis_db.Publish(ctx, rdb, "telemetry", m.Payload())
		}).
		SetOrderMatters(false). // Allow out-of-order handling for speed
		SetKeepAlive(keepAlive).
		SetPingTimeout(pingTimeout)

	client := mqtt.NewClient(opts)
	for {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			logrus.Errorf("failed to connect to MQTT broker: %v", token.Error())
			logrus.Infof("Retrying in %v...", retryInterval)
			time.Sleep(retryInterval)
			continue
		}

		if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
			client.Disconnect(250)
			logrus.Errorf("failed to subscribe to MQTT broker: %v", token.Error())
			logrus.Infof("Retrying in %v...", retryInterval)
			time.Sleep(retryInterval)
			continue
		}
		break
	}

	logrus.Info("Connected to MQTT broker")

	<-quit
	client.Disconnect(250)
	logrus.Info("MQTT broker disconnected")
}
