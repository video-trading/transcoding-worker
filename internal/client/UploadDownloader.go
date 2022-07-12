package client

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"video_transcoding_worker/internal/types"
)

type UploadDownloader struct {
	config     *types.Config
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
}

func NewUploadDownloader(config *types.Config) *UploadDownloader {
	return &UploadDownloader{
		config: config,
	}
}

func (u *UploadDownloader) Init() {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(u.config.UploadDownloaderConfig.Region),
		Credentials: credentials.NewStaticCredentials(
			u.config.UploadDownloaderConfig.AccessKey,
			u.config.UploadDownloaderConfig.SecretKey,
			"",
		),
		Endpoint: aws.String(fmt.Sprintf("%s.digitaloceanspaces.com", u.config.UploadDownloaderConfig.Region)),
	})

	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)

	u.downloader = downloader
	u.uploader = uploader
}

func (u *UploadDownloader) Download(bucket string, fileName string) string {
	baseName := path.Base(fileName)
	downloadPath := path.Join("download", baseName)
	if _, err := os.Stat("download"); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll("download", os.ModePerm)
		if err != nil {
			log.Printf("Cannot create directory: %s", err)
		}
	}

	out, err := os.Create(downloadPath)
	if err != nil {
		log.Printf("Cannot create file: %s", err)
	}

	defer out.Close()

	_, err = u.downloader.DownloadWithContext(aws.BackgroundContext(), out, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})

	if err != nil {
		log.Printf("Cannot download file %s: %s", fileName, err)
	}

	log.Printf("Download finished for file %s", fileName)
	return downloadPath
}

func (u *UploadDownloader) Upload(filename string, bucket string) error {
	uploadName := path.Join("converted", path.Base(filename))
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Cannot open converted file: %s", err)
		return err
	}

	defer file.Close()

	_, err = u.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(uploadName),
		Body:   file,
	})

	if err != nil {
		log.Printf("Cannot upload file: %s", err)
		return err
	}
	log.Printf("Successully uploaded file")

	return nil
}
