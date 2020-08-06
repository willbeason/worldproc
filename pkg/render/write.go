package render

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func WriteImage(img image.Image, file string) {
	out, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = out.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
