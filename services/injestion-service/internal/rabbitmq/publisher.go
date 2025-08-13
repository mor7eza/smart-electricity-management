package rabbitmq

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

const QUEUE_NAME = "loggers_data"

func Publisher(url string, messages chan []byte) {
	conn := connectWithRetries(url)
	defer conn.Close()

	channel := openChannel(conn)
	defer channel.Close()

	declareQueue(channel, QUEUE_NAME)

	for msg := range messages {
		if err := publishMessage(channel, QUEUE_NAME, msg); err != nil {
			logrus.Errorf("failed to publish message: %v", err)
		}
	}

	logrus.Warn("Publisher: messages channel closed, exiting")
}

func connectWithRetries(url string) *amqp.Connection {
	for {
		conn, err := amqp.Dial(url)
		if err == nil {
			logrus.Info("Connected to RabbitMQ")
			return conn
		}
		logrus.Errorf("failed to connect to RabbitMQ: %v", err)
		logrus.Info("Retrying in 10 seconds...")
		time.Sleep(10 * time.Second)
	}
}

func openChannel(conn *amqp.Connection) *amqp.Channel {
	channel, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("failed to open RabbitMQ channel: %v", err)
	}
	return channel
}

func declareQueue(ch *amqp.Channel, queueName string) {
	_, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		logrus.Fatalf("failed to declare queue %q: %v", queueName, err)
	}
}

func publishMessage(ch *amqp.Channel, queueName string, body []byte) error {
	return ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
