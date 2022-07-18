package clients

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"

	"video_transcoding_worker/internal/types"
)

type Converter struct {
	resolution types.Resolution
}

// NewConverter Create a new ffmpeg client
func NewConverter() *Converter {
	return &Converter{}
}

// getOutputName will generate a random unique output filename
func (f *Converter) getOutputName() string {
	id := uuid.New()
	return fmt.Sprintf("%s.mp4", id.String())
}

// getArgs will generate a ffmpeg arguments based on the system
func (f *Converter) getArgs() (map[string]interface{}, error) {
	m := make(map[string]interface{})

	if runtime.GOOS == types.MacOS {
		m["c:v"] = "h264_videotoolbox"
	} else {
		return m, fmt.Errorf("OS %s is not supported", runtime.GOOS)
	}

	m["vf"] = f.getResolution(f.resolution)
	m["crf"] = "0"
	m["c:a"] = "copy"
	return m, nil
}

func (f *Converter) getResolution(resolution types.Resolution) string {
	return fmt.Sprintf("scale=-1:%s", resolution)
}

// Convert input file to an output file and returns the output filename
func (f *Converter) Convert(filename string, resolution types.Resolution) (string, error) {
	f.resolution = resolution
	outputName := path.Join("converted", f.getOutputName())
	if _, err := os.Stat("converted"); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll("converted", os.ModePerm)
		if err != nil {
			log.Printf("Cannot create directory: %s", err)
		}
	}

	arguments, err := f.getArgs()

	if err != nil {
		log.Printf("OS is not supported")
	}

	err = ffmpeg.Input(filename).
		Output(outputName, arguments).
		OverWriteOutput().
		ErrorToStdOut().
		Run()

	return outputName, err
}
