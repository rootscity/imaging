package imaging

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	"fmt"
)

// Decode reads an image from r.
func Decode(r io.Reader) (image.Image, error) {
	fmt.Println("Decode before")
	img, _, err := image.Decode(r)
	fmt.Println("Decode after", err)
	if err != nil {
		return nil, err
	}
	return toNRGBA(img), nil
}


