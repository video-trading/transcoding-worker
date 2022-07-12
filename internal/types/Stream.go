package types

type Stream interface {
	Output(outputName string, arguments map[string]interface{}) *Stream
	OverWriteOutput() *Stream
	ErrorToStdOut() *Stream
	Run() *Stream
}
