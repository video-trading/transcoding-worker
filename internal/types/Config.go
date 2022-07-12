package types

type Config struct {
	MessageQueueConfig     MessageQueueConfig
	UploadDownloaderConfig UploadDownloaderConfig
}

type UploadDownloaderConfig struct {
	Region    string
	AccessKey string
	SecretKey string
}

type MessageQueueConfig struct {
	MessageQueueURL string
}
