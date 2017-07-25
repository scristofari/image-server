package resizer

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

type Provider interface {
	Get(filename string) (io.Reader, error)
	Put(filename string, image multipart.File) error
}

type diskProvider struct {
	Provider
}

func (d *diskProvider) Get(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (d *diskProvider) Put(filename string, image multipart.File) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	_, err = io.Copy(f, image)
	return err
}

/**
type awsProvider struct {
	io.ReadWriter
}

func (a *awsProvider) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (a *awsProvider) Read(p []byte) (n int, err error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	ctx := context.Background()
	result, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String("my-bucket"),
		Key:    aws.String("my-key"),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			return 0, fmt.Errorf("failed to get the requested file, %s", err)
		}
		return 0, fmt.Errorf("failed to get the requested file, %s", err)
	}
	defer result.Body.Close()

	return 0, nil
}
*/
