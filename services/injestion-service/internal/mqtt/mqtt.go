package mqtt_broker

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func RunMqttClient(brokerURL, clientID, topic string, quit chan bool) {
	opts := mqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(clientID).
		SetDefaultPublishHandler(ProcessMessages).
		SetKeepAlive(60 * time.Second).
		SetPingTimeout(1 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("error connecting to MQTT server: %v", token.Error())
	}

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		client.Disconnect(250)
		log.Fatalf("error subscribing to MQTT topic: %v", token.Error())
	}

	<-quit
	client.Disconnect(250)
	logrus.Info("MQTT disconnected")
}

func ProcessMessages(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("recieved")
}
