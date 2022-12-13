package types

type TranscodingInfo struct {
	VideoId        string     `json:"videoId"`
	TargetQuality  Resolution `json:"targetQuality"`
	Status         string     `json:"status"`
	TranscodingUrl SignedUrl  `json:"TranscodingUrl"`
	VideoUrl       SignedUrl  `json:"VideoUrl"`
}

type TranscodingResult struct {
	Quality Resolution        `json:"quality"`
	Status  TranscodingStatus `json:"status"`
}
