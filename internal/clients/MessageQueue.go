package clients

import (
	"log"
	"video_transcoding_worker/internal/constant"

	amqp "github.com/rabbitmq/amqp091-go"

	"video_transcoding_worker/internal/types"
)

type MessageQueue struct {
	config types.MessageQueueConfig
}

// NewMessageQueue creates a message queue
func NewMessageQueue(config types.MessageQueueConfig) *MessageQueue {
	return &MessageQueue{
		config: config,
	}
}

func (c *MessageQueue) Init(name string) (*amqp.Connection, *amqp.Channel, amqp.Queue) {
	conn, err := amqp.Dial(c.config.MessageQueueURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	// connect to the exchange and routing key
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	// declare the exchange
	err = ch.ExchangeDeclare(
		string(c.config.Exchange),
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
	}

	// maximum retry limit is set to 3
	args := amqp.Table{
		"x-queue-type":     "quorum",
		"x-delivery-limit": constant.MaxRetry,
	}

	// declare the queue
	q, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,  // arguments
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	// bind the queue to the exchange
	err = ch.QueueBind(
		q.Name,
		string(c.config.RoutingKey),
		string(c.config.Exchange),
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Failed to bind a queue: %s", err)
	}

	return conn, ch, q
}
