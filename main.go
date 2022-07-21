package main

import (
	"fmt"
	"os"

	"video_transcoding_worker/internal/channels"
	"video_transcoding_worker/internal/clients"
	"video_transcoding_worker/internal/handlers"
	"video_transcoding_worker/internal/types"
)

func main() {
	var forever chan struct{}
	config := types.Config{
		UploadDownloaderConfig: types.UploadDownloaderConfig{
			DownloadPath: "download",
		},
		TranscodingConfig: types.TranscodingConfig{
			URL: "http://localhost:8080",
		},
		AnalyzerConfig: types.AnalyzerConfig{
			CoverPath: "cover",
		},
		ConverterConfig: types.ConverterConfig{
			OutputFolder: "coverted",
		},
	}

	converterClient := clients.NewConverter(config.ConverterConfig)
	uploadDownloader := clients.NewUploadDownloader(config.UploadDownloaderConfig)
	cleaner := clients.NewCleaner()
	transcodingClient := clients.NewTranscodingClient(config.TranscodingConfig)
	analyzer := clients.NewAnalyzer(config.AnalyzerConfig)

	go func() {
		transcodingConfig := types.MessageQueueConfig{
			MessageQueueURL: os.Getenv("message_queue"),
			Topic:           channels.Transcode,
		}
		fmt.Println("Setting up transcoding handler")
		transcodingJobHandler := handlers.NewTranscodingJobHandler(transcodingConfig, converterClient, uploadDownloader, cleaner, transcodingClient)
		transcodingJobHandler.Init()
		transcodingJobHandler.Run()
	}()

	go func() {
		analyzingConfig := types.MessageQueueConfig{
			MessageQueueURL: os.Getenv("message_queue"),
			Topic:           channels.Analyze,
		}

		fmt.Println("Setting up analyzing job handler")
		analyzingHandler := handlers.NewAnalyzingJobHandler(analyzingConfig, uploadDownloader, cleaner, analyzer, transcodingClient)
		analyzingHandler.Init()
		analyzingHandler.Run()
	}()
	<-forever
}
