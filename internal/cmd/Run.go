package main

import (
	"os"

	"video_transcoding_worker/internal/client"
	"video_transcoding_worker/internal/types"
)

func main() {
	config := types.Config{
		MessageQueueConfig: types.MessageQueueConfig{
			MessageQueueURL: os.Getenv("message_queue"),
		},
		UploadDownloaderConfig: types.UploadDownloaderConfig{
			Region:    "sgp1",
			AccessKey: os.Getenv("access_key"),
			SecretKey: os.Getenv("secret_key"),
		},
	}
	converterClient := client.NewConverter()
	uploadDownloader := client.NewUploadDownloader(&config)
	cleaner := client.NewCleaner()

	messageQueueReceiver := client.NewMessageQueueReceiver(config, converterClient, uploadDownloader, cleaner)
	messageQueueReceiver.Init()
	messageQueueReceiver.Run()
}
