package rabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
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

func NewService() {}

func (rs *RabbitMQService) RunPublisher(rdb *redis.Client) {
	var (
		ctx = context.Background()
	)

	rs.connectWithRetries()
	defer rs.Conn.Close()

	rs.openChannel()
	defer rs.Channel.Close()

	rs.declareQueue()

	for {
		result, err := rdb.BRPop(ctx, 0*time.Second, "telemetry").Result()
		if err != nil {
			logrus.Errorf("error reading from Redis telemetry list: %v", err)
			time.Sleep(time.Second)
			continue
		}

		valByte := []byte(result[1])

		if err := publishMessage(channel, queueName, valByte); err != nil {
			logrus.Errorf("failed to publish message: %v", err)
		}
	}
}

func (rs *RabbitMQService) connectWithRetries() {
	var (
		url = viper.GetString("RABBITMQ_URL")
	)

	for {
		conn, err := amqp.Dial(url)
		if err == nil {
			logrus.Info("Connected to RabbitMQ")
			rs.Conn = conn
			return
		}
		logrus.Errorf("failed to connect to RabbitMQ: %v", err)
		logrus.Infof("Retrying in %v...", retryInterval)
		time.Sleep(retryInterval)
	}
}

func (rs *RabbitMQService) openChannel() {
	channel, err := rs.Conn.Channel()
	if err != nil {
		logrus.Fatalf("failed to open RabbitMQ channel: %v", err)
	}
	rs.Channel = channel
}

func (rs *RabbitMQService) declareQueue() {
	_, err := rs.Channel.QueueDeclare(
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

func (rs *RabbitMQService) publishMessage(body []byte) error {
	return rs.Channel.Publish(
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
