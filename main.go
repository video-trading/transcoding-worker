package main

import (
	"os"
	"video_transcoding_worker/worker"
)

func main() {
	worker.Setup(os.Getenv("ENDPOINT"), os.Getenv("JWT_TOKEN"), os.Getenv("MESSAGE_QUEUE"))
}
