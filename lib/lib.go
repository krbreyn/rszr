package lib

import (
	"bytes"
	"fmt"
	"image"
	"os"

	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	// _ "golang.org/x/image/webp"

	"golang.org/x/image/draw"
)

func GetImageSizeInBytes(img image.Image, format string) (int, error) {
	var buf bytes.Buffer

	switch format {
	case ".jpeg", ".jpg":
		err := jpeg.Encode(&buf, img, nil)
		if err != nil {
			return 0, err
		}
	case ".png":
		err := png.Encode(&buf, img)
		if err != nil {
			return 0, err
		}
	default:
		return 0, fmt.Errorf("unsupported format: %s", format)
	}

	return buf.Len(), nil
}

func OpenImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(file)
	file.Close()
	if err != nil {
		return nil, err
	}
	return img, nil
}

func ResizeImageNearest(src image.Image, maxWidth, maxHeight int) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	if srcWidth <= maxWidth && srcHeight <= maxHeight {
		return src
	}
	scaleX := float64(maxWidth) / float64(srcWidth)
	scaleY := float64(maxHeight) / float64(srcHeight)

	scale := min(scaleY, scaleX)

	newWidth := int(float64(srcWidth) * scale)
	newHeight := int(float64(srcHeight) * scale)
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	return dst
}
