package clients

import (
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
