package resizer

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	uuid "github.com/satori/go.uuid"
)

const (
	mb = 1 << 20
)

var (
	UploadMaxSize = 5 * mb
	outputDir     = "../../files"
)

// Uploadfile : ___
func Uploadfile(image multipart.File) (string, error) {
	uuid := uuid.NewV4().String()
	f, err := os.OpenFile(outputDir+"/"+uuid+".png", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	_, err = io.Copy(f, image)
	if err != nil {
		return "", fmt.Errorf("failed to copy: %v", err)
	}

	return uuid, nil
}
