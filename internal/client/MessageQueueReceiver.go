package client

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"video_transcoding_worker/internal/types"
)

type MessageQueueReceiver struct {
	messageQueue     *MessageQueue
	converter        *Converter
	uploadDownloader *UploadDownloader
	cleaner          *Cleaner
	conn             *amqp.Connection
	channel          *amqp.Channel
	queue            amqp.Queue
}

func NewMessageQueueReceiver(config types.Config, converter *Converter, uploadDownloader *UploadDownloader, cleaner *Cleaner) *MessageQueueReceiver {
	messageQueue := NewMessageQueue(config)

	return &MessageQueueReceiver{
		messageQueue:     messageQueue,
		converter:        converter,
		uploadDownloader: uploadDownloader,
		cleaner:          cleaner,
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (m *MessageQueueReceiver) Init() {
	m.uploadDownloader.Init()
	conn, channel, queue := m.messageQueue.Init()

	m.conn = conn
	m.channel = channel
	m.queue = queue
}

func (m *MessageQueueReceiver) Run() {
	defer m.conn.Close()
	defer m.channel.Close()

	msgs, err := m.channel.Consume(
		m.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			body := d.Body
			var data types.Message
			err := json.Unmarshal(body, &data)
			if err != nil {
				fmt.Printf("Cannot decode: %s", err)
			}
			// download file
			downloadPath := m.uploadDownloader.Download(data.BucketName, data.FileName)
			convertedPath, err := m.converter.Convert(downloadPath, data.Resolution)
			if err != nil {
				log.Printf("Cannot convert video: %s", err)
				return
			}
			err = m.uploadDownloader.Upload(convertedPath, data.BucketName)
			if err == nil {
				log.Printf("Successfully upload converted data with resolution %s", data.Resolution)
			}
			m.cleaner.Clean([]string{
				downloadPath,
				convertedPath,
			})
		}
	}()

	<-forever
}
