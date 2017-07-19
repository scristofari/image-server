package resizer

import (
	"io"
	"mime/multipart"
	"os"

	uuid "github.com/satori/go.uuid"
)

var (
	outputDir = "../../files"
)

// Uploadfile : ___
func Uploadfile(image multipart.File) (string, error) {
	uuid := uuid.NewV4().String()
	f, err := os.OpenFile(outputDir+"/"+uuid+".png", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	io.Copy(f, image)

	return uuid, nil
}
