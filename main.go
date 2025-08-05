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
	flag.Usage = func() {
		fmt.Println("Usage of rszr:")
	}

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
		// fmt.Printf("%s is %d kb big\n", stat.Name(), stat.Size()/1024)

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

			if size <= *kbPtr*1024 {
				fmt.Printf("size %d fits under %d\n", size, *kbPtr*1024)
			} else {
				fmt.Printf("size %d over %d\n", size, *kbPtr*1024)
			}
			fmt.Println("TODO")
		}

		new_name := strings.TrimSuffix(stat.Name(), ext) + "_" + new_name_suffix + ext
		fmt.Println(new_name)
	}

}

func handleSize(img image.Image, widthPtr *int, heightPtr *int) image.Image {
	bounds := img.Bounds()
	if *widthPtr == 0 {
		*widthPtr = bounds.Dx()
	}
	if *heightPtr == 0 {
		*heightPtr = bounds.Dy()
	}
	if bounds.Dx() <= *widthPtr && bounds.Dy() <= *heightPtr {
		return img
	}
	img = lib.ResizeImageNearest(img, *widthPtr, *heightPtr)
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
