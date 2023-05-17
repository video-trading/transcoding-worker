package handlers

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"video_transcoding_worker/internal/clients"
	"video_transcoding_worker/internal/constant"
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
			// check retry count
			log.Println("Receiving analyzing job")
			body := d.Body
			failed, videoId := m.handle(body)
			if failed {
				log.Printf("Failed to handle message: %s", d.MessageId)
				count := d.Headers["x-delivery-count"]
				countInt, ok := count.(int64)
				if ok {
					if countInt == constant.MaxRetry {
						//TODO: notify the server that the job is failed
						log.Printf("Max retry reached for message: %s", d.MessageId)
						err := m.transcodingClient.SubmitFailedAnalyzingResult(videoId)
						if err != nil {
							log.Printf("Failed to submit failed analyzing result: %s", err)
							return
						}
					}
				}
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

// handle handles the message and will return true if there is an error and return the video id
func (m *AnalyzingJobHandler) handle(body []byte) (bool, string) {
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
		return true, ""
	}
	// download file
	videoPath, err = m.uploadDownloader.Download(data.Video.PreviewUrl)
	if err != nil {
		log.Printf("Cannot download video: %s\n", err)
		return true, data.VideoId
	}
	analyzingResult, err := m.analyzer.Analyze(videoPath, data.VideoId, data.Video.Key)
	if err != nil {
		log.Printf("Cannot analyze video: %s\n", err)
		return true, data.VideoId
	}

	fmt.Printf("Analyzing result: %v\n", analyzingResult)
	coverPath = analyzingResult.Cover
	if err != nil {
		log.Printf("Cannot analyze video: %s\n", err)
		return true, data.VideoId
	}

	err = m.uploadDownloader.Upload(data.Thumbnail.Url, analyzingResult.Cover)
	if err != nil {
		log.Printf("Cannot upload cover: %s\n", err)
		return true, data.VideoId
	}

	err = m.transcodingClient.SubmitAnalyzingResult(analyzingResult)
	if err != nil {
		log.Printf("Cannot submit analyzing result: %s\n", err)
		return true, data.VideoId
	}
	log.Printf("Finished analyzing job\n")
	return false, data.VideoId
}
