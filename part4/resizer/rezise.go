package resizer

import (
	"fmt"
	"image"
	"image/png"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

const (
	presetMaxSize = 2000
)

type Query struct {
	Type   string
	Preset *Preset
}

type Preset struct {
	Width  uint
	Height uint
}

func Resize(filename string, q *Query) (image.Image, error) {
	filename = outputDir + "/" + filename

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	var i image.Image
	switch q.Type {
	case "r":
		i = resize.Resize(q.Preset.Width, q.Preset.Height, img, resize.Lanczos2)
	case "t":
		i = resize.Thumbnail(q.Preset.Width, q.Preset.Height, img, resize.Lanczos2)
	}

	return i, nil
}

func GetQueryFromURL(u *url.URL) (*Query, error) {
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the query: %s", err)
	}

	if len(m) != 1 {
		return nil, fmt.Errorf("only one manipulation permitted")
	}

	for t, preset := range m {
		p, err := getPreset(preset[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get the preset: %s", err)
		}

		return &Query{
			Type:   t,
			Preset: p,
		}, nil
	}

	return nil, fmt.Errorf("failed to get the manipulation %s", m)
}

func getPreset(p string) (*Preset, error) {
	hash := strings.Split(p, "x")
	if len(hash) != 2 {
		return nil, fmt.Errorf("failed to get the proper size")
	}

	width, err := strconv.Atoi(hash[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get the proper width")
	}
	height, err := strconv.Atoi(hash[1])
	if err != nil {
		return nil, fmt.Errorf("failed to get the proper height")
	}

	if height > presetMaxSize {
		return nil, fmt.Errorf("invalid preset, too high, got %d, max %d", height, presetMaxSize)
	}
	if width > presetMaxSize {
		return nil, fmt.Errorf("invalid preset, too high, got %d, max %d", width, presetMaxSize)
	}

	return &Preset{
		Width:  uint(width),
		Height: uint(height),
	}, nil
}
