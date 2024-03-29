package clients

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	path2 "path"
	"strings"
	"video_transcoding_worker/internal/types"
)

type UploadDownloader struct {
	config types.UploadDownloaderConfig
}

// NewUploadDownloader Creates a new upload downloader
func NewUploadDownloader(config types.UploadDownloaderConfig) *UploadDownloader {
	return &UploadDownloader{
		config: config,
	}
}

// Init initializes the client
func (u *UploadDownloader) Init() {
	if _, err := os.Stat(u.config.DownloadPath); os.IsNotExist(err) {
		os.MkdirAll(u.config.DownloadPath, os.ModePerm)
	}
}

// Download the file using the signed url. Key is the file name like a/b/c.png is the key from the pre-signed url object
func (u *UploadDownloader) Download(downloadURL string) (string, error) {

	// Build fileName from fullPath
	fileURL, err := url.Parse(downloadURL)
	if err != nil {
		return "", err
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := path2.Join(u.config.DownloadPath, segments[len(segments)-1])

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(downloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d\n", fileName, size)
	return fileName, nil
}

// Upload the file using the signed url
func (u *UploadDownloader) Upload(uploadURL string, fileName string) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		err = fmt.Errorf("unable to readfile %s", fileName)
		return err
	}

	req, err := http.NewRequest("PUT", uploadURL, bytes.NewBuffer(file))
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s, with content: %s\n", res.Status, string(content))
		return err
	}
	return nil
}
