package types

type Message struct {
	FileName   string `json:"file_name"`
	BucketName string `json:"bucket_name"`
	Resolution string `json:"resolution"`
}
