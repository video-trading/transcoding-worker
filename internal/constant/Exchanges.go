package constant

type Exchange string

const (
	// TranscodingExchange channel is used to send transcoding jobs to the transcoding service.
	TranscodingExchange Exchange = "video"

	// AnalyzingExchange channel is used to send analyzing jobs to the analyzing service.
	AnalyzingExchange Exchange = "video"
)
