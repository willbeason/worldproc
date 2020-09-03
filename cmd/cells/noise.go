package main

import (
	"fmt"
	"github.com/willbeason/worldproc/pkg/geodesic"
	"github.com/willbeason/worldproc/pkg/noise"
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

	nSpheres := 9
	spheres := []*geodesic.Geodesic{geodesic.Dodecahedron()}
	for i := 0; i < nSpheres; i++ {
		spheres = append(spheres, geodesic.Chamfer(spheres[i]))
	}

	depth := 20
	perlinNoise := noise.NewPerlinFractal(10, depth, 0.8)

	heights := make([]float64, screen.Width*screen.Height)

	sphere := spheres[nSpheres]
	for x := 0; x < screen.Width; x++ {
		//fmt.Println(x)
		for y := 0; y < screen.Height; y++ {
			pidx := y*screen.Width + x
			angle := projection.Pixels[pidx]
			idx := geodesic.Find(spheres, angle.Vector())
			pos := sphere.Centers[idx]

			heights[pidx] = perlinNoise.ValueAt(pos)
		}
	}

	csl := render.NewColorScale(
		[]render.ColorPoint{
			{-0.1, color.RGBA{0,0, 255, 255}},
			{-0.01, color.RGBA{0, 32, 255, 255}},
			{0, color.RGBA{0, 128, 255, 255}},
			{0.02, color.RGBA{192, 192, 32, 255}},
			{0.1, color.RGBA{0, 192, 0, 255}},
			{0.15, color.RGBA{64, 128, 64, 255}},
			{0.2, color.RGBA{64, 96, 64, 255}},
			{0.5, color.RGBA{192, 192, 192, 255}},
			{1.0, color.RGBA{255, 255, 255, 255}},
		})

	img := image.NewRGBA(image.Rect(0, 0, screen.Width, screen.Height))
	screen.PaintAndShade(heights, csl, img)

	out, err := os.Create(fmt.Sprintf("map2-%d-%d.png", nSpheres, depth))
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
