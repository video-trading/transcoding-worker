package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"video_transcoding_worker/internal/types"
)

type TranscodingClient struct {
	config types.TranscodingConfig
}

// NewTranscodingClient Creates a new transcoding client for submitting transcoding result and analyzing result
func NewTranscodingClient(config types.TranscodingConfig) *TranscodingClient {
	return &TranscodingClient{
		config: config,
	}
}

// SubmitFinishedResult SubmitTranscodingResult will submit the transcoding result to the transcoding service
func (t *TranscodingClient) SubmitFinishedResult(id string, transcodingInfo *types.TranscodingResult) error {
	transcodingInfo.Status = types.COMPLETED
	jsonValue, _ := json.Marshal(transcodingInfo)
	requestURL := fmt.Sprintf("%s/transcoding/%s", t.config.URL, id)

	// Send http request with JWT token
	client := &http.Client{}
	req, err := http.NewRequest("PATCH", requestURL, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.config.JWTToken))
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err := fmt.Errorf("get status %s with error %s", response.Status, content)
		return err
	}

	return nil
}

// SubmitAnalyzingResult Submit the analyzing result to the transcoding service
func (t *TranscodingClient) SubmitAnalyzingResult(analyzingResult *types.AnalyzingResult) error {
	jsonValue, _ := json.Marshal(analyzingResult)
	requestURL := fmt.Sprintf("%s/video/%s/analyzing/result", t.config.URL, analyzingResult.VideoId)

	// Send http request with JWT token
	client := &http.Client{}
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.config.JWTToken))
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		content, _ := io.ReadAll(response.Body)
		err := fmt.Errorf("get status %s with error %s", response.Status, content)
		return err
	}

	return nil
}

// SubmitFailedAnalyzingResult Submit failed analyzing result to the transcoding service
func (t *TranscodingClient) SubmitFailedAnalyzingResult(videoId string) error {
	requestURL := fmt.Sprintf("%s/video/%s/analyzing/failed", t.config.URL, videoId)

	// Send http request with JWT token
	client := &http.Client{}
	req, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.config.JWTToken))
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		content, _ := io.ReadAll(response.Body)
		err := fmt.Errorf("get status %s with error %s", response.Status, content)
		return err
	}

	return nil
}
