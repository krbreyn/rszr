package main

import (
	"flag"
	"fmt"
	"image"
	"math"
	"os"
	"path"
	"strings"

	"github.com/krbreyn/rszr/lib"
)

func main() {
	flag.Usage = usage

	widthPtr := flag.Int("w", 0, "maximum width")
	heightPtr := flag.Int("h", 0, "maximum height")
	kbPtr := flag.Int("kb", 0, "maximum file size, in kb")
	percentPtr := flag.Int("%", 0, "percentage to reduce dimensions by")

	flag.Parse()
	args := flag.Args()

	if *widthPtr == 0 && *heightPtr == 0 && *kbPtr == 0 && *percentPtr == 0 {
		fmt.Println("no options provided")
		usageDie()
	}

	if len(args) < 1 {
		fmt.Println("no image(s) provided")
		usageDie()
	}

	for _, p := range args {
		ext := path.Ext(p)
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			fmt.Println("unsupported file type:", p)
			continue
		}
		stat, err := os.Stat(p)
		if err != nil {
			fmt.Println(err)
			continue
		}

		img, err := lib.OpenImage(p)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var new_name_suffix string

		if *percentPtr != 0 {
			new_name_suffix += fmt.Sprintf("percent_%d_", *percentPtr)

			img = handlePercent(img, percentPtr)
			fmt.Printf("reduced %s dimensions by %d percent\n", stat.Name(), *percentPtr)
		}

		if *widthPtr != 0 || *heightPtr != 0 {
			if *widthPtr != 0 {
				new_name_suffix += fmt.Sprintf("width_%d_", *widthPtr)
			}
			if *heightPtr != 0 {
				new_name_suffix += fmt.Sprintf("height_%d_", *heightPtr)
			}

			img = handleSize(img, widthPtr, heightPtr)
			fmt.Printf("shrunk %s to fit within %dx%d\n", stat.Name(), *widthPtr, *heightPtr)

		}

		if *kbPtr != 0 {
			new_name_suffix += fmt.Sprintf("kb_%d_", *kbPtr)

			size, err := lib.GetImageSizeInBytes(img, ext)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if size > *kbPtr*1024 {
				img, err = handleFilesize(img, size, kbPtr, ext)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
		}

		new_name := strings.TrimSuffix(stat.Name(), ext) + "_" + new_name_suffix + ext
		err = lib.SaveImage(img, new_name)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("saved %s as %s\n", stat.Name(), new_name)
		}
	}

	fmt.Println("All done!")
}

func handleFilesize(img image.Image, size int, kbPtr *int, ext string) (image.Image, error) {
	factor := float64(*kbPtr*1024) / float64(size)
	bounds := img.Bounds()
	newWidth := int(float64(bounds.Dx()) * factor)
	newHeight := int(float64(bounds.Dy()) * factor)
	img = lib.ResizeImageNearest(img, newWidth, newHeight)
	size, err := lib.GetImageSizeInBytes(img, ext)
	if err != nil {
		return img, err
	}
	if size > *kbPtr*1024 {
		for {
			newWidth -= 10
			newHeight -= 10
			img = lib.ResizeImageNearest(img, newWidth, newHeight)
			size, err := lib.GetImageSizeInBytes(img, ext)
			if err != nil {
				return img, err
			}
			if size <= *kbPtr*1024 {
				break
			}
		}
	}

	return img, nil
}

func usage() {
	fmt.Println("Usage of rszr:")
	fmt.Println("rszr [-w, -h, -kb, -w int] [image paths ...]")
	fmt.Println("-% int")
	fmt.Println("\tpercentage to reduce dimensions by")
	fmt.Println("-w int")
	fmt.Println("\tmaximum width")
	fmt.Println("-h int")
	fmt.Println("\tmaximum height")
	fmt.Println("-kb int")
	fmt.Println("\tmaximum file size, in kb")
}

func handleSize(img image.Image, widthPtr *int, heightPtr *int) image.Image {
	bounds := img.Bounds()
	newWidth := *widthPtr
	newHeight := *heightPtr
	if *widthPtr == 0 {
		newWidth = bounds.Dx()
	}
	if *heightPtr == 0 {
		newHeight = bounds.Dy()
	}
	if bounds.Dx() <= newWidth && bounds.Dy() <= newHeight {
		return img
	}
	img = lib.ResizeImageNearest(img, newWidth, newHeight)
	return img
}

func handlePercent(img image.Image, percentPtr *int) image.Image {
	bounds := img.Bounds()
	factor := (100.0 - float64(*percentPtr)) / 100.0

	newWidth := int(math.Round(float64(bounds.Dx()) * factor))
	newHeight := int(math.Round(float64(bounds.Dy()) * factor))

	fmt.Println(newWidth, newHeight)

	img = lib.ResizeImageNearest(img, newWidth, newHeight)
	return img
}

func usageDie() {
	flag.Usage()
	os.Exit(0)
}
