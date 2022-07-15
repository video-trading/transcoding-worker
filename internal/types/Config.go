package types

type Config struct {
	UploadDownloaderConfig UploadDownloaderConfig
	TranscodingConfig      TranscodingConfig
}

type UploadDownloaderConfig struct {
}

type MessageQueueConfig struct {
	MessageQueueURL string
	Topic           string
}

type TranscodingConfig struct {
	URL string
}
