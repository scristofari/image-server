package resizer

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Provider interface {
	Get(filename string) (io.ReadCloser, error)
	Put(filename string, image multipart.File) error
}

type DiskProvider struct {
	Provider
}

func (d *DiskProvider) Get(filename string) (io.ReadCloser, error) {
	file, err := os.Open("../../files/" + filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (d *DiskProvider) Put(filename string, image multipart.File) error {
	f, err := os.OpenFile("../../files/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	_, err = io.Copy(f, image)
	return err
}

type AWSProvider struct {
	Provider
}

func (a *AWSProvider) Get(filename string) (io.ReadCloser, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	downloader := s3manager.NewDownloaderWithClient(svc)

	buffer := &aws.WriteAtBuffer{}
	_, err := downloader.Download(buffer, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}

	return aws.ReadSeekCloser(bytes.NewReader(buffer.Bytes())), nil
}

func (a *AWSProvider) Put(filename string, image multipart.File) error {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	uploader := s3manager.NewUploaderWithClient(svc)

	upParams := &s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(filename),
		Body:   image,
	}

	_, err := uploader.Upload(upParams)
	if err != nil {
		return fmt.Errorf("failed to upload object, %v\n", err)
	}

	return nil
}
