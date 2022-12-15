package types

import "video_transcoding_worker/internal/constant"

type Config struct {
	UploadDownloaderConfig UploadDownloaderConfig
	TranscodingConfig      TranscodingConfig
	ConverterConfig        ConverterConfig
	AnalyzerConfig         AnalyzerConfig
}

type UploadDownloaderConfig struct {
	DownloadPath string
}

type MessageQueueConfig struct {
	MessageQueueURL string
	Exchange        constant.Exchange
	RoutingKey      constant.RoutingKey
}
type ConverterConfig struct {
	OutputFolder string
}

type TranscodingConfig struct {
	URL      string
	JWTToken string
}

type AnalyzerConfig struct {
	CoverPath string
}
