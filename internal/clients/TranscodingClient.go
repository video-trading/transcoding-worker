package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"video_transcoding_worker/internal/types"
)

type TranscodingClient struct {
	config types.TranscodingConfig
}

func NewTranscodingClient(config types.TranscodingConfig) *TranscodingClient {
	return &TranscodingClient{
		config: config,
	}
}

func (t *TranscodingClient) SubmitFinishedResult(transcodingInfo *types.TranscodingInfo) error {
	transcodingInfo.Status = types.Uploaded
	jsonValue, _ := json.Marshal(transcodingInfo)
	requestURL := fmt.Sprintf("%s/video/transcoding/result", t.config.URL)
	response, err := http.Post(requestURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err := fmt.Errorf("get status %s with error %s", response.Status, content)
		return err
	}

	return nil
}

func (t *TranscodingClient) SubmitAnalyzingResult(analyzingResult *types.AnalyzingResult) error {
	jsonValue, _ := json.Marshal(analyzingResult)
	requestURL := fmt.Sprintf("%s/video/transcoding/analyzing", t.config.URL)
	response, err := http.Post(requestURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err := fmt.Errorf("get status %s with error %s", response.Status, content)
		return err
	}

	return nil
}
