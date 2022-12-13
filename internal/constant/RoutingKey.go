package constant

type RoutingKey string

const (
	// TranscodeRoutingKey TranscodingExchange channel is used to send transcoding jobs to the transcoding service.
	TranscodeRoutingKey RoutingKey = "transcoding.*"
	// AnalyzeRoutingKey AnalyzingExchange channel is used to send analyzing jobs to the analyzing service.
	AnalyzeRoutingKey RoutingKey = "analyzing.*"
)
