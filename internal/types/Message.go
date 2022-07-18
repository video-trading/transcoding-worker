package types

type TranscodingInfo struct {
	Id                  string            `json:"id"`
	Source              string            `json:"source"`
	OriginalVideoSource string            `json:"originalVideoSource"`
	VideoId             string            `json:"videoId"`
	Quality             Resolution        `json:"quality"`
	BucketName          string            `json:"bucketName"`
	FileName            string            `json:"fileName"`
	CreatedAt           string            `json:"CreatedAt"`
	Status              TranscodingStatus `json:"status"`
}
