package image

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/h2non/bimg"
)

// Manager manages the uploaded images
type Manager struct {
	Directory           string
	MaxWidth, MaxHeight int
}

// Upload reads, optimizes and creates a file in a destination directory with a unique name
func (m Manager) Upload(r io.Reader) (key string, err error) {
	key = m.newKey()

	buf, err := m.process(r)
	if err != nil {
		return "", err
	}

	fpath := path.Join(m.Directory, key)
	if err := os.WriteFile(fpath, buf, 0666); err != nil {
		return "", fmt.Errorf("failed to write to file '%s': %w", fpath, err)
	}

	return
}

func (m Manager) newKey() string {
	return time.Now().Format("2006-01-02T15-04-05-99") + ".webp"
}

func (m Manager) process(r io.Reader) ([]byte, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	img := bimg.NewImage(buf)
	opt := bimg.Options{Type: bimg.WEBP}

	// resize if needed
	m.optimizeSize(img, &opt)

	buf, err = img.Process(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to process img: %w", err)
	}

	return buf, nil
}

func (m Manager) optimizeSize(img *bimg.Image, opt *bimg.Options) error {
	size, err := img.Size()
	if err != nil {
		return fmt.Errorf("failed to get image size: %w", err)
	}

	// in case the image already fits our constraints, do nothing
	if size.Height < m.MaxHeight && size.Width < m.MaxWidth {
		return nil
	}

	var (
		newWidth, newHeight int
	)

	// detect which side of the image has largest impact
	wRatio := size.Width / m.MaxWidth
	hRatio := size.Height / m.MaxHeight

	// i don't know why i do it like this, but i feel like this needs to be done
	if wRatio < hRatio {
		newWidth = min(m.MaxWidth, size.Width)
		newHeight = proportion(size.Height, size.Width, newWidth)
	} else {
		newHeight = min(m.MaxHeight, size.Height)
		newWidth = proportion(size.Width, size.Height, newHeight)
	}

	opt.Width = newWidth
	opt.Height = newHeight
	return nil
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func proportion(x1, y1, y2 int) (x2 int) {
	/*

		x1/x2 = y1/y2

		x2 = (x1 * y2) / y1

	*/

	x2 = int(float32(y2) * float32(x1) / float32(y1))
	return
}
