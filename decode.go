package imaging

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"fmt"
	"bufio"
	"golang.org/x/image/webp"
)

// Decode reads an image from r.
func Decode(r io.Reader) (image.Image, error) {
	fmt.Println("Decode before")
	img, _, err := decode(r)
	fmt.Println("Decode after", err)
	if err != nil {
		return nil, err
	}
	return toNRGBA(img), nil
}

const (
	leHeader = "II\x2A\x00" // Header for little-endian files.
	beHeader = "MM\x00\x2A" // Header for big-endian files.
	pngHeader = "\x89PNG\r\n\x1a\n"
)

func init() {
	RegisterFormat("bmp", "BM????\x00\x00\x00\x00", bmp.Decode, bmp.DecodeConfig)
	RegisterFormat("tiff", leHeader, tiff.Decode, tiff.DecodeConfig)
	RegisterFormat("tiff", beHeader, tiff.Decode, tiff.DecodeConfig)
	RegisterFormat("webp", "RIFF????WEBPVP8", webp.Decode, webp.DecodeConfig)
	RegisterFormat("gif", "GIF8?a", gif.Decode, gif.DecodeConfig)
	RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)
	RegisterFormat("png", pngHeader, png.Decode, png.DecodeConfig)
}

// A format holds an image format's name, magic header and how to decode it.
type format struct {
	name, magic  string
	decode       func(io.Reader) (image.Image, error)
	decodeConfig func(io.Reader) (image.Config, error)
}

// Formats is the list of registered formats.
var formats []format

// RegisterFormat registers an image format for use by Decode.
// Name is the name of the format, like "jpeg" or "png".
// Magic is the magic prefix that identifies the format's encoding. The magic
// string can contain "?" wildcards that each match any one byte.
// Decode is the function that decodes the encoded image.
// DecodeConfig is the function that decodes just its configuration.
func RegisterFormat(name, magic string, decode func(io.Reader) (image.Image, error), decodeConfig func(io.Reader) (image.Config, error)) {
	formats = append(formats, format{name, magic, decode, decodeConfig})
}

// A reader is an io.Reader that can also peek ahead.
type reader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

// asReader converts an io.Reader to a reader.
func asReader(r io.Reader) reader {
	if rr, ok := r.(reader); ok {
		return rr
	}
	return bufio.NewReader(r)
}

// Match reports whether magic matches b. Magic may contain "?" wildcards.
func match(magic string, b []byte) bool {
	if len(magic) != len(b) {
		return false
	}
	for i, c := range b {
		if magic[i] != c && magic[i] != '?' {
			return false
		}
	}
	return true
}

// Sniff determines the format of r's data.
func sniff(r reader) format {
	for _, f := range formats {
		b, err := r.Peek(len(f.magic))
		if err == nil && match(f.magic, b) {
			return f
		}
	}
	return format{}
}

// Decode decodes an image that has been encoded in a registered format.
// The string returned is the format name used during format registration.
// Format registration is typically done by an init function in the codec-
// specific package.
func decode(r io.Reader) (image.Image, string, error) {
	rr := asReader(r)
	f := sniff(rr)
	if f.decode == nil {
		return nil, "", image.ErrFormat
	}
	m, err := f.decode(rr)
	return m, f.name, err
}
