package types

type TranscodingStatus = string

const (
	Pending             TranscodingStatus = "pending"
	Transcoding         TranscodingStatus = "transcoding"
	TranscodingFinished TranscodingStatus = "transcodeFinished"
	PendingUpload       TranscodingStatus = "pendingUpload"
	Uploaded            TranscodingStatus = "uploaded"
	Success             TranscodingStatus = "success"
)
