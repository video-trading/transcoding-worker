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
	conn, channel, queue := m.messageQueue.Init("transcoding")

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
			log.Println("Receiving a transcoding job")
			if m.handle(d, body) {
				log.Printf("Failed to handle message: %s", d.MessageId)
			}
		}
	}()

	log.Printf("Listening to the job request")
	<-forever
}

func (m *TranscodingJobHandler) handle(d amqp.Delivery, body []byte) bool {
	var videoPath string
	var convertedVideoPath string
	var uploadError error

	defer func() {
		log.Println("Cleaning up transcoding job")
		m.cleaner.Clean([]string{
			videoPath,
			convertedVideoPath,
		})
	}()
	var data types.TranscodingInfo
	err := json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("Cannot decode: %s", err)
	}
	// download file
	videoPath, err = m.uploadDownloader.Download(data.VideoUrl.PreviewUrl)
	if err != nil {
		log.Printf("Cannot download video: %s", err)
		return true
	}
	convertedVideoPath, convertionError := m.converter.Convert(videoPath, data.TargetQuality)
	if convertionError != nil {
		log.Printf("Cannot convert video: %s", err)
	}

	if convertionError == nil {
		uploadError = m.uploadDownloader.Upload(data.TranscodingUrl.Url, convertedVideoPath)
		if uploadError != nil {
			log.Printf("Cannot upload video: %s", err)
		}
	}

	// send transcoding result
	result := types.TranscodingResult{
		Quality: data.TargetQuality,
		Status:  types.COMPLETED,
	}

	if convertionError != nil || uploadError != nil {
		result.Status = types.FAILED
	}

	err = m.transcodingClient.SubmitFinishedResult(data.VideoId, &result)
	if err != nil {
		log.Printf("Cannot submit transcoding result: %s", err)
		return true
	}
	log.Printf("Successfully upload converted data with resolution %s", data.TargetQuality)
	return convertionError != nil || uploadError != nil
}
