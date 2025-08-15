package rabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RabbitMQService struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

const (
	queueName     = "loggers_data"
	retryInterval = 10 * time.Second
)

func NewService() *RabbitMQService {
	url := viper.GetString("RABBITMQ_URL")

	for {
		rs, err := connect(url)
		if err != nil {
			logrus.Errorf("failed to connect to RabbitMQ: %v", err)
			logrus.Infof("Retrying in %v...", retryInterval)
			time.Sleep(retryInterval)
			continue
		}

		logrus.Info("Connected to RabbitMQ")
		return rs
	}
}

func connect(url string) (*RabbitMQService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	err = declareQueue(ch)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQService{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func declareQueue(ch *amqp.Channel) error {
	_, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)

	return err
}

func (rs *RabbitMQService) PublishMessage(body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return rs.Channel.PublishWithContext(
		ctx,
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
