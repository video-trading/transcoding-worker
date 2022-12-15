package types

type TranscodingStatus = string

/**
 * TranscodingStatus is the status of a transcoding job.
 */
const (
	PENDING    TranscodingStatus = "PENDING"
	PROCESSING TranscodingStatus = "PROCESSING"
	COMPLETED  TranscodingStatus = "COMPLETED"
	FAILED     TranscodingStatus = "FAILED"
)
