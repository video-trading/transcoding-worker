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
	conn, channel, queue := m.messageQueue.Init("analyzing_job")

	m.conn = conn
	m.channel = channel
	m.queue = queue
}

func (m *AnalyzingJobHandler) Run() {
	msgs, err := m.channel.Consume(
		m.queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Println("Receiving analyzing job")
			body := d.Body
			if m.handle(body) {
				log.Printf("Failed to handle message: %s", d.MessageId)
				// request to retry
				err = d.Nack(false, true)
				if err != nil {
					log.Printf("Failed to reject message: %s", d.MessageId)
				}
			} else {
				err := d.Ack(false)
				if err != nil {
					log.Printf("Failed to ack message: %s", d.MessageId)
				}
			}
		}
	}()

	log.Printf("Listening to the job request")
	<-forever
}

// handle handles the message and will return true if there is an error
func (m *AnalyzingJobHandler) handle(body []byte) bool {
	var data types.AnalyzingJob

	// path for the cover
	var coverPath string
	// path for the video
	var videoPath string

	defer func() {
		log.Println("Cleaning up analyzing job")
		m.cleaner.Clean([]string{
			videoPath,
			coverPath,
		})
	}()

	err := json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("Cannot decode: %s", err)
		return true
	}
	// download file
	videoPath, err = m.uploadDownloader.Download(data.Video.PreviewUrl)
	if err != nil {
		log.Printf("Cannot download video: %s\n", err)
		return true
	}
	analyzingResult, err := m.analyzer.Analyze(videoPath, data.VideoId, data.Video.Key)
	if err != nil {
		log.Printf("Cannot analyze video: %s\n", err)
		return true
	}

	fmt.Printf("Analyzing result: %v\n", analyzingResult)
	coverPath = analyzingResult.Cover
	if err != nil {
		log.Printf("Cannot analyze video: %s\n", err)
		return true
	}

	err = m.uploadDownloader.Upload(data.Thumbnail.Url, analyzingResult.Cover)
	if err != nil {
		log.Printf("Cannot upload cover: %s\n", err)
		return true
	}

	err = m.transcodingClient.SubmitAnalyzingResult(analyzingResult)
	if err != nil {
		log.Printf("Cannot submit analyzing result: %s\n", err)
		return true
	}
	log.Printf("Finished analyzing job\n")
	return false
}
