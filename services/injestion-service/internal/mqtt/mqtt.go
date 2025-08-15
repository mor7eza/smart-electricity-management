package mqtt_broker

import (
	"context"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MqttService struct {
	Client mqtt.Client
}

const (
	keepAlive     = 30 * time.Second
	pingTimeout   = 1 * time.Second
	retryInterval = 10 * time.Second
)

func NewService(rdb *redis.Client) *MqttService {
	var (
		url      = viper.GetString("MQTT_URL")
		clientID = viper.GetString("MQTT_CLIENT_ID")
		ctx      = context.Background()
	)

	opts := mqtt.NewClientOptions().
		AddBroker(url).
		SetClientID(clientID).
		SetDefaultPublishHandler(func(c mqtt.Client, m mqtt.Message) {
			rdb.Publish(ctx, "telemetry", m.Payload())
		}).
		SetOrderMatters(false). // Allow out-of-order handling for speed
		SetKeepAlive(keepAlive).
		SetPingTimeout(pingTimeout)

	return &MqttService{
		Client: mqtt.NewClient(opts),
	}
}

func (ms *MqttService) Run() {
	var (
		topic = viper.GetString("MQTT_TOPIC")
	)
	for {
		if token := ms.Client.Connect(); token.Wait() && token.Error() != nil {
			logrus.Errorf("failed to connect to MQTT broker: %v", token.Error())
			logrus.Infof("Retrying in %v...", retryInterval)
			time.Sleep(retryInterval)
			continue
		}

		if token := ms.Client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
			ms.Client.Disconnect(250)
			logrus.Errorf("failed to subscribe to MQTT broker: %v", token.Error())
			logrus.Infof("Retrying in %v...", retryInterval)
			time.Sleep(retryInterval)
			continue
		}
		break
	}

	logrus.Info("Connected to MQTT broker")
}
