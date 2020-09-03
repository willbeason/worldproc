package main

import (
	"fmt"
	"github.com/willbeason/worldproc/pkg/geodesic"
	"github.com/willbeason/worldproc/pkg/render"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	screen := render.Screen{
		Width:  1500,
		Height: 750,
	}

	projection := render.Project(screen, render.Equirectangular{})

	img := image.NewRGBA(image.Rect(0, 0, screen.Width, screen.Height))

	spheres := []*geodesic.Geodesic{geodesic.Dodecahedron()}
	for i := 0; i < 8; i++ {
		spheres = append(spheres, geodesic.Chamfer(spheres[i]))
	}

	for x := 0; x < screen.Width; x++ {
		//fmt.Println(x)
		for y := 0; y < screen.Height; y++ {
			angle := projection.Pixels[y*screen.Width+x]
			idx := geodesic.Find(spheres, angle.Vector())
			//idx2 := geodesic.NaiveFind(spheres[len(spheres)-1], angle.Vector())
			border := false

			if y > 0 {
				if idx != geodesic.Find(spheres, projection.Pixels[(y-1)*screen.Width+x].Vector()) {
					border = true
				}
			}
			if !border && y < screen.Height - 1 {
				if idx != geodesic.Find(spheres, projection.Pixels[(y+1)*screen.Width+x].Vector()) {
					border = true
				}
			}
			if !border && x > 0 {
				if idx != geodesic.Find(spheres, projection.Pixels[y*screen.Width+(x-1)].Vector()) {
					border = true
				}
			}
			if !border && x < screen.Width - 1 {
				if idx != geodesic.Find(spheres, projection.Pixels[(y*screen.Width)+(x+1)].Vector()) {
					border = true
				}
			}

			if border {
				img.Set(x, y, color.Gray{Y: 64})
			} else {
				img.Set(x, y, color.Gray{Y: 192})
			}
		}
	}


	out, err := os.Create("file-08.png")
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


