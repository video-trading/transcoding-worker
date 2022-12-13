package main

import (
	"fmt"
	"os"

	"video_transcoding_worker/internal/clients"
	"video_transcoding_worker/internal/constant"
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
			URL:      "http://localhost:3000",
			JWTToken: os.Getenv("JWT_TOKEN"),
		},
		AnalyzerConfig: types.AnalyzerConfig{
			CoverPath: "cover",
		},
		ConverterConfig: types.ConverterConfig{
			OutputFolder: "coverted",
		},
	}

	//converterClient := clients.NewConverter(config.ConverterConfig)
	uploadDownloader := clients.NewUploadDownloader(config.UploadDownloaderConfig)
	cleaner := clients.NewCleaner()
	transcodingClient := clients.NewTranscodingClient(config.TranscodingConfig)
	analyzer := clients.NewAnalyzer(config.AnalyzerConfig)

	//go func() {
	//	transcodingConfig := types.MessageQueueConfig{
	//		MessageQueueURL: os.Getenv("message_queue"),
	//		Exchange:        constant.TranscodingExchange,
	//		RoutingKey:      constant.TranscodeRoutingKey,
	//	}
	//	fmt.Println("Setting up transcoding handler")
	//	transcodingJobHandler := handlers.NewTranscodingJobHandler(transcodingConfig, converterClient, uploadDownloader, cleaner, transcodingClient)
	//	transcodingJobHandler.Init()
	//	transcodingJobHandler.Run()
	//}()

	go func() {
		analyzingConfig := types.MessageQueueConfig{
			MessageQueueURL: os.Getenv("MESSAGE_QUEUE"),
			Exchange:        constant.AnalyzingExchange,
			RoutingKey:      constant.AnalyzeRoutingKey,
		}

		fmt.Println("Setting up analyzing job handler")
		analyzingHandler := handlers.NewAnalyzingJobHandler(analyzingConfig, uploadDownloader, cleaner, analyzer, transcodingClient)
		analyzingHandler.Init()
		analyzingHandler.Run()
	}()
	<-forever
}
