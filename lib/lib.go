package lib

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"path"

	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	// _ "golang.org/x/image/webp"

	"golang.org/x/image/draw"
)

func EncodeImage(img image.Image, format string) (bytes.Buffer, error) {
	var buf bytes.Buffer

	switch format {
	case ".jpeg", ".jpg":
		err := jpeg.Encode(&buf, img, nil)
		if err != nil {
			return buf, err
		}
	case ".png":
		err := png.Encode(&buf, img)
		if err != nil {
			return buf, err
		}
	default:
		return buf, fmt.Errorf("unsupported format: %s", format)
	}

	return buf, nil
}

func GetImageSizeInBytes(img image.Image, format string) (int, error) {
	buf, err := EncodeImage(img, format)
	if err != nil {
		return 0, err
	}

	return buf.Len(), nil
}

func SaveImage(img image.Image, filename string) error {
	buf, err := EncodeImage(img, path.Ext(filename))
	if err != nil {
		return err
	}
	dst := path.Join(".", filename)
	err = os.WriteFile(dst, buf.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
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
