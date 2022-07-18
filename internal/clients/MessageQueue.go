package clients

import (
	"log"

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

func (c *MessageQueue) Init() (*amqp.Connection, *amqp.Channel, amqp.Queue) {
	conn, err := amqp.Dial(c.config.MessageQueueURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %s", err)
	}

	queue, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	err = channel.QueueBind(queue.Name, c.config.Topic, "amq.topic", false, nil)
	if err != nil {
		log.Fatalf("Failed to bind a queue: %s", err)
	}

	return conn, channel, queue
}
