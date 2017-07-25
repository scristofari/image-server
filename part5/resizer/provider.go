package resizer

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

type awsProvider struct {
	Provider
}

func (a *awsProvider) Get(filename string) (io.ReadCloser, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	ctx := context.Background()
	result, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String("image-server-filer"),
		Key:    aws.String(filename),
	})

	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			return nil, fmt.Errorf("failed to get the requested file, %s", err)
		}
		return nil, fmt.Errorf("failed to get the requested file, %s", err)
	}

	return result.Body, nil
}

func (a *awsProvider) Put(filename string, image multipart.File) error {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	ctx := context.Background()
	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String("image-server-filer"),
		Key:    aws.String(filename),
		Body:   image,
	})
	if err != nil {
		return fmt.Errorf("failed to upload object, %v\n", err)
	}

	return nil
}
