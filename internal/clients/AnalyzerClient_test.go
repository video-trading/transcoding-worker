package clients

import (
	"gopkg.in/vansante/go-ffprobe.v2"
	"reflect"
	"testing"
	"video_transcoding_worker/internal/types"
)

func TestAnalyzer_GetVideoResolution(t *testing.T) {
	type fields struct {
		config types.AnalyzerConfig
	}
	type args struct {
		height int
	}
	type Test struct {
		name    string
		fields  fields
		args    args
		want    types.Resolution
		wantErr bool
	}

	var tests = []Test{
		Test{
			"Test 1",
			fields{
				config: types.AnalyzerConfig{
					CoverPath: "cover",
				},
			},
			args{
				height: 1070,
			},
			types.Resolution720p,
			false,
		},
		Test{
			"Test 1",
			fields{
				config: types.AnalyzerConfig{
					CoverPath: "cover",
				},
			},
			args{
				height: 1080,
			},
			types.Resolution1080p,
			false,
		},
		Test{
			"Test 1",
			fields{
				config: types.AnalyzerConfig{
					CoverPath: "cover",
				},
			},
			args{
				height: 2170,
			},
			types.Resolution2160p,
			false,
		},
		Test{
			"Test 1",
			fields{
				config: types.AnalyzerConfig{
					CoverPath: "cover",
				},
			},
			args{
				height: 120,
			},
			types.ResolutionUnknown,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Analyzer{
				config: tt.fields.config,
			}
			got, err := a.GetVideoResolution(tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVideoResolution() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetVideoResolution() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnalyzer_getStream(t *testing.T) {
	type fields struct {
		config types.AnalyzerConfig
	}
	type args struct {
		streams []*ffprobe.Stream
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ffprobe.Stream
		wantErr bool
	}{
		{
			"Should be able to get video stream",
			fields{
				config: types.AnalyzerConfig{
					CoverPath: "cover",
				},
			},
			args{
				streams: []*ffprobe.Stream{
					{
						Index:     0,
						CodecType: "audio",
					},
					{
						Index:     1,
						CodecType: "video",
					},
				},
			},
			&ffprobe.Stream{
				Index:     1,
				CodecType: "video",
			},
			false,
		},
		{
			"Should be able to get video stream",
			fields{
				config: types.AnalyzerConfig{
					CoverPath: "cover",
				},
			},
			args{
				streams: []*ffprobe.Stream{
					{
						Index:     0,
						CodecType: "video",
					},
				},
			},
			&ffprobe.Stream{
				Index:     0,
				CodecType: "video",
			},
			false,
		},
		{
			"Should not be able to get video stream",
			fields{
				config: types.AnalyzerConfig{
					CoverPath: "cover",
				},
			},
			args{
				streams: []*ffprobe.Stream{
					{
						Index:     1,
						CodecType: "audio",
					},
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Analyzer{
				config: tt.fields.config,
			}
			got, err := a.getStream(tt.args.streams)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getStream() got = %v, want %v", got, tt.want)
			}
		})
	}
}
