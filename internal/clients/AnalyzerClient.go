package clients

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/vansante/go-ffprobe.v2"

	"video_transcoding_worker/internal/types"
)

type Analyzer struct {
	config types.AnalyzerConfig
}

// NewAnalyzer Creates a new ffmpeg client
func NewAnalyzer(config types.AnalyzerConfig) *Analyzer {
	return &Analyzer{
		config: config,
	}
}

// getArgs will generate a ffmpeg arguments based on the system
func (a *Analyzer) getArgs() map[string]interface{} {
	m := make(map[string]interface{})
	m["ss"] = "00:00:00"
	m["frames:v"] = "1"
	return m
}

// screenshot will generate a screenshot for the video
func (a *Analyzer) screenshot(filename string) (string, error) {
	fmt.Printf("Taking screenshot for file: %s\n", filename)
	id := uuid.New()
	cover := fmt.Sprintf("%s.png", id)
	args := a.getArgs()

	if _, err := os.Stat(a.config.CoverPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll("cover", os.ModePerm)
		if err != nil {
			log.Printf("Cannot create directory: %s", err)
		}
	}
	cover = path.Join(a.config.CoverPath, cover)

	err := ffmpeg.Input(filename).Output(cover, args).ErrorToStdOut().Run()
	if err != nil {
		return "", err
	}
	return cover, nil
}

func (a *Analyzer) GetVideoResolution(height int) (types.Resolution, error) {
	if height >= 144 && height < 240 {
		return types.Resolution144p, nil
	} else if height >= 240 && height < 360 {
		return types.Resolution240p, nil
	} else if height >= 360 && height < 480 {
		return types.Resolution360p, nil
	} else if height >= 480 && height < 720 {
		return types.Resolution480p, nil
	} else if height >= 720 && height < 1080 {
		return types.Resolution720p, nil
	} else if height >= 1080 && height < 1440 {
		return types.Resolution1080p, nil
	} else if height >= 1440 && height < 2160 {
		return types.Resolution1440p, nil
	} else if height >= 2160 {
		return types.Resolution2160p, nil
	} else {
		return types.ResolutionUnknown, fmt.Errorf("unknown resolution")
	}
}

// Analyze will analyze the video and generate a AnalyzingResult
func (a *Analyzer) Analyze(filename string, videoId string, uploadFileName string) (*types.AnalyzingResult, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, filename)

	if err != nil {
		return nil, err
	}

	cover, err := a.screenshot(filename)
	if err != nil {
		return nil, err
	}

	if len(data.Streams) < 1 {
		return nil, fmt.Errorf("video's stream size is less than 0")
	}

	stream := data.Streams[0]
	resolution, err := a.GetVideoResolution(stream.Height)
	if err != nil {
		return nil, err
	}

	result := types.AnalyzingResult{
		VideoId:   videoId,
		Length:    data.Format.DurationSeconds,
		Quality:   resolution,
		FrameRate: stream.AvgFrameRate,
		FileName:  uploadFileName,
		Cover:     cover,
	}

	return &result, nil
}
