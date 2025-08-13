package mqtt_broker

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func RunMqttClient(brokerURL, clientID, topic string, out chan []byte, quit chan bool) {
	opts := mqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(clientID).
		SetDefaultPublishHandler(func(c mqtt.Client, m mqtt.Message) {
			out <- m.Payload()
		}).
		SetKeepAlive(60 * time.Second).
		SetPingTimeout(1 * time.Second)

	client := mqtt.NewClient(opts)
	for {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			logrus.Errorf("error connecting to MQTT server: %v", token.Error())
			logrus.Info("trying to connect in 10 seconds...")
			time.Sleep(10 * time.Second)
			continue
		}

		if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
			client.Disconnect(250)
			logrus.Errorf("error subscribing to MQTT topic: %v", token.Error())
			logrus.Info("trying to reconnect and subscribe in 10 seconds...")
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}

	logrus.Info("Successfully subscribed to EMQX")

	<-quit
	client.Disconnect(250)
	logrus.Info("MQTT disconnected")
}
