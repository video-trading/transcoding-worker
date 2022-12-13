package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"video_transcoding_worker/internal/clients"
	"video_transcoding_worker/internal/types"
)

type TranscodingJobHandler struct {
	messageQueue      *clients.MessageQueue
	converter         *clients.Converter
	uploadDownloader  *clients.UploadDownloader
	cleaner           *clients.Cleaner
	transcodingClient *clients.TranscodingClient
	conn              *amqp.Connection
	channel           *amqp.Channel
	queue             amqp.Queue
}

func NewTranscodingJobHandler(config types.MessageQueueConfig, converter *clients.Converter,
	uploadDownloader *clients.UploadDownloader, cleaner *clients.Cleaner, transcodingClient *clients.TranscodingClient) *TranscodingJobHandler {
	messageQueue := clients.NewMessageQueue(config)

	return &TranscodingJobHandler{
		messageQueue:      messageQueue,
		converter:         converter,
		uploadDownloader:  uploadDownloader,
		cleaner:           cleaner,
		transcodingClient: transcodingClient,
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (m *TranscodingJobHandler) Init() {
	m.uploadDownloader.Init()
	conn, channel, queue := m.messageQueue.Init()

	m.conn = conn
	m.channel = channel
	m.queue = queue
}

func (m *TranscodingJobHandler) Run() {
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
			log.Printf("Receiving message: %s", d.MessageId)
			var data types.TranscodingInfo
			err := json.Unmarshal(body, &data)
			if err != nil {
				fmt.Printf("Cannot decode: %s", err)
			}
			// download file
			//TODO: Fix this
			downloadPath, err := m.uploadDownloader.Download(data.OriginalVideoSource)
			if err != nil {
				log.Printf("Cannot download video: %s", err)
				return
			}

			convertedPath, err := m.converter.Convert(downloadPath, data.Quality)
			if err != nil {
				log.Printf("Cannot convert video: %s", err)
				return
			}

			err = m.uploadDownloader.Upload(data.Source, convertedPath)
			if err == nil {
				err := m.transcodingClient.SubmitFinishedResult(&data)
				if err != nil {
					log.Printf("Cannot submit transcoding result: %s", err)
					return
				}
				log.Printf("Successfully upload converted data with resolution %s", data.Quality)
			} else {
				log.Printf("Cannot upload converted data: %s", err)
			}
		}
	}()

	log.Printf("Listening to the job request")
	<-forever
}
