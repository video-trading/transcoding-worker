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
	config *types.Config
}

func NewTranscodingClient(config *types.Config) *TranscodingClient {
	return &TranscodingClient{
		config: config,
	}
}

func (t *TranscodingClient) SubmitFinishedResult(transcodingInfo *types.TranscodingInfo) error {
	transcodingInfo.Status = types.Uploaded
	jsonValue, _ := json.Marshal(transcodingInfo)
	requestURL := fmt.Sprintf("%s/video/transcoding/result", t.config.TranscodingConfig.URL)
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
