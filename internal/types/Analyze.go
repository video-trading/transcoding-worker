package types

type AnalyzingResult struct {
	VideoId   string     `json:"videoId"`
	Quality   Resolution `json:"quality"`
	FrameRate string     `json:"frameRate"`
	Length    float64    `json:"length"`
	Cover     string     `json:"cover"`
	FileName  string     `json:"fileName"`
}

type AnalyzingJob struct {
	// VideoId is the id of the video to be analyzed.
	VideoId string `json:"videoId"`
	// Video is the video to be analyzed.
	Video SignedUrl `json:"video"`
	// Thumbnail is the thumbnail for us to upload.
	Thumbnail SignedUrl `json:"thumbnail"`
}

// SignedUrl is a signed url for a file.
type SignedUrl struct {
	// Url is the signed url and granted PUT permission.
	Url string `json:"url"`
	// Key is the file name.
	Key string `json:"key"`
	// PreviewUrl is the signed url and granted GET permission.
	PreviewUrl string `json:"previewUrl"`
}
