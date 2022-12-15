package clients

import (
	"strings"
	"testing"
	"video_transcoding_worker/internal/types"
)

func TestConverter_getResolution(t *testing.T) {
	type fields struct {
		resolution types.Resolution
		config     types.ConverterConfig
	}
	type args struct {
		resolution types.Resolution
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"Get resolution 1080p",
			fields{
				resolution: types.Resolution1080p,
				config: types.ConverterConfig{
					OutputFolder: "/tmp",
				},
			},
			args{
				resolution: types.Resolution1080p,
			},
			"scale=-1:1080p",
		},
		{
			"Get resolution 720p",
			fields{
				resolution: types.Resolution720p,
				config: types.ConverterConfig{
					OutputFolder: "/tmp",
				},
			},
			args{
				resolution: types.Resolution720p,
			},
			"scale=-1:720p",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Converter{
				resolution: tt.fields.resolution,
				config:     tt.fields.config,
			}
			if got := f.getResolution(tt.args.resolution); got != tt.want {
				t.Errorf("getResolution() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_getArgs(t *testing.T) {
	t.Run("Get resolution 1080p", func(t *testing.T) {
		f := &Converter{
			resolution: types.Resolution1080p,
			config: types.ConverterConfig{
				OutputFolder: "/tmp",
			},
		}
		got, err := f.getArgs()
		if err != nil {
			t.Errorf("getArgs() error = %v", err)
			return
		}

		if got["vf"] != "scale=-1:1080p" {
			t.Errorf("getArgs() vf = %v, want %v", got["vf"], "scale=-1:1080p")
		}
	})
}

func TestConverter_getOutputName(t *testing.T) {
	type fields struct {
		resolution types.Resolution
		config     types.ConverterConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"Get output name",
			fields{
				resolution: types.Resolution1080p,
				config: types.ConverterConfig{
					OutputFolder: "/tmp",
				},
			},
			".mp4",
		},
		{
			"Get output name",
			fields{
				resolution: types.Resolution1440p,
				config: types.ConverterConfig{
					OutputFolder: "/tmp",
				},
			},
			".mp4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Converter{
				resolution: tt.fields.resolution,
				config:     tt.fields.config,
			}
			if got := f.getOutputName(); !strings.Contains(got, tt.want) {
				t.Errorf("getOutputName() = %v, want %v", got, tt.want)
			}
		})
	}
}
