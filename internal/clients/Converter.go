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
	config     types.ConverterConfig
}

// NewConverter Creates a new ffmpeg client
func NewConverter(config types.ConverterConfig) *Converter {
	return &Converter{
		config: config,
	}
}

// getOutputName will generate a random unique output filename
func (f *Converter) getOutputName() string {
	id := uuid.New()
	return fmt.Sprintf("%s.mp4", id.String())
}

// getArgs will generate a ffmpeg arguments based on the system
func (f *Converter) getArgs() (map[string]interface{}, error) {
	m := make(map[string]interface{})

	// check if it is running in test mode
	if runtime.GOOS == types.MacOS {
		m["c:v"] = "h264_videotoolbox"
		m["b:v"] = f.getByteRates(f.resolution)
	}

	m["vf"] = f.getResolution(f.resolution)
	m["crf"] = "0"
	m["c:a"] = "copy"
	return m, nil
}

// getResolution will generate a ffmpeg resolution based on the system
func (f *Converter) getResolution(resolution types.Resolution) string {
	return fmt.Sprintf("scale=-1:%s", resolution)
}

func (f *Converter) getByteRates(resolution types.Resolution) string {
	switch resolution {
	case types.Resolution2160p:
		return "12000k"
	case types.Resolution1080p:
		return "8000k"
	case types.Resolution720p:
		return "4000k"
	case types.Resolution480p:
		return "2000k"
	case types.Resolution360p:
		return "1000k"
	case types.Resolution240p:
		return "500k"
	default:
		return "8000k"
	}
}

// Convert input file to an output file and returns the output filename
func (f *Converter) Convert(filename string, resolution types.Resolution) (string, error) {
	f.resolution = resolution
	outputName := path.Join(f.config.OutputFolder, f.getOutputName())
	if _, err := os.Stat(f.config.OutputFolder); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(f.config.OutputFolder, os.ModePerm)
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
