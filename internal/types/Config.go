package types

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
	Topic           string
}
type ConverterConfig struct {
	OutputFolder string
}

type TranscodingConfig struct {
	URL string
}

type AnalyzerConfig struct {
	CoverPath string
}
