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
	VideoId  string `json:"videoId"`
	Cover    string `json:"cover"`
	Source   string `json:"source"`
	FileName string `json:"fileName"`
}
