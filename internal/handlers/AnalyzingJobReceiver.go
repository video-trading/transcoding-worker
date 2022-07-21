package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"video_transcoding_worker/internal/clients"
	"video_transcoding_worker/internal/types"
)

type AnalyzingJobHandler struct {
	messageQueue      *clients.MessageQueue
	uploadDownloader  *clients.UploadDownloader
	cleaner           *clients.Cleaner
	transcodingClient *clients.TranscodingClient
	analyzer          *clients.Analyzer
	conn              *amqp.Connection
	channel           *amqp.Channel
	queue             amqp.Queue
}

func NewAnalyzingJobHandler(config types.MessageQueueConfig,
	uploadDownloader *clients.UploadDownloader, cleaner *clients.Cleaner,
	analyzer *clients.Analyzer, transcodingClient *clients.TranscodingClient) *AnalyzingJobHandler {
	messageQueue := clients.NewMessageQueue(config)

	return &AnalyzingJobHandler{
		messageQueue:      messageQueue,
		uploadDownloader:  uploadDownloader,
		cleaner:           cleaner,
		transcodingClient: transcodingClient,
		analyzer:          analyzer,
	}
}

func (m *AnalyzingJobHandler) Init() {
	m.uploadDownloader.Init()
	conn, channel, queue := m.messageQueue.Init()

	m.conn = conn
	m.channel = channel
	m.queue = queue
}

func (m *AnalyzingJobHandler) Run() {
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
			var data types.AnalyzingJob
			err := json.Unmarshal(body, &data)
			if err != nil {
				fmt.Printf("Cannot decode: %s", err)
			}
			// download file
			downloadPath, err := m.uploadDownloader.Download(data.Source)
			if err != nil {
				log.Printf("Cannot download video: %s\n", err)
				return
			}
			analyzingResult, err := m.analyzer.Analyze(downloadPath, data.VideoId, data.FileName)
			fmt.Printf("Analyzing result: %v\n", analyzingResult)
			if err != nil {
				log.Printf("Cannot analyze video: %s\n", err)
				return
			}

			err = m.uploadDownloader.Upload(data.Cover, analyzingResult.Cover)
			if err != nil {
				log.Printf("Cannot upload cover: %s\n", err)
				return
			}

			err = m.transcodingClient.SubmitAnalyzingResult(analyzingResult)
			if err != nil {
				log.Printf("Cannot submit analyzing result: %s\n", err)
				return
			}
			log.Printf("Finished analyzing job\n")
			//m.cleaner.Clean([]string{
			//	downloadPath,
			//	analyzingResult.Cover,
			//})
		}
	}()

	log.Printf("Listening to the job request")
	<-forever
}
