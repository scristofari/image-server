package resizer

import (
	"fmt"
	"mime/multipart"

	uuid "github.com/satori/go.uuid"
)

const (
	mb = 1 << 20
)

var (
	UploadMaxSize = 5 * mb
)

// Uploadfile : ___
func Uploadfile(p Provider, image multipart.File) (string, error) {
	filename := uuid.NewV4().String() + ".png"

	err := p.Put(filename, image)

	if err != nil {
		return "", fmt.Errorf("failed to copy: %v", err)
	}

	return filename, nil
}
