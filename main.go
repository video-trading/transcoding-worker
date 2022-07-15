package main

import (
	"os"

	"video_transcoding_worker/internal/clients"
	"video_transcoding_worker/internal/handlers"
	"video_transcoding_worker/internal/types"
)

func main() {
	config := types.Config{
		UploadDownloaderConfig: types.UploadDownloaderConfig{},
		TranscodingConfig: types.TranscodingConfig{
			URL: "http://localhost:8080",
		},
	}

	converterClient := clients.NewConverter()
	uploadDownloader := clients.NewUploadDownloader(&config)
	cleaner := clients.NewCleaner()
	transcodingClient := clients.NewTranscodingClient(&config)

	transcodingConfig := types.MessageQueueConfig{
		MessageQueueURL: os.Getenv("message_queue"),
		Topic:           "transcodingWorker",
	}
	messageQueueReceiver := handlers.NewTranscodingJobHandler(transcodingConfig, converterClient, uploadDownloader, cleaner, transcodingClient)
	messageQueueReceiver.Init()
	messageQueueReceiver.Run()
}
