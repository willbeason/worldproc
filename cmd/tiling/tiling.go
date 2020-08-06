package main

import (
	"fmt"
	"github.com/willbeason/hydrology/pkg/geodesic"
	"github.com/willbeason/hydrology/pkg/render"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func findCell(spheres []*geodesic.Geodesic, v geodesic.Vector, weight float64) int {
	previous := geodesic.Find(spheres[:len(spheres)-1], v)
	next := geodesic.Find(spheres, v)

	if previous == next {
		return previous
	}

	previousV := spheres[len(spheres)-2].Centers[previous]
	nextV := spheres[len(spheres)-1].Centers[next]

	dPrevious := math.Sqrt(geodesic.DistSq(v, previousV))
	dNext := math.Sqrt(geodesic.DistSq(v, nextV))

	if dPrevious * weight < dNext {
		return previous
	}
	return next
}

func drawBorders(spheres []*geodesic.Geodesic, projection render.Projection, weight float64) []bool {
	cellId := make([]int, len(projection.Pixels))
	for i, px := range projection.Pixels {
		cellId[i] = findCell(spheres, px.Vector(), weight)
	}

	isBorder := make([]bool, len(projection.Pixels))
	for i, cell := range cellId {
		x := i / projection.Width
		y := i % projection.Width

		if y > 0 && cellId[i-1] != cell {
			isBorder[i] = true
			continue
		}
		if y < projection.Width - 1 && cellId[i+1] != cell {
			isBorder[i] = true
			continue
		}
		if x > 0 && cellId[i - projection.Width] != cell {
			isBorder[i] = true
			continue
		}
		if x < projection.Height - 1 && cellId[i + projection.Width] != cell {
			isBorder[i] = true
			continue
		}
	}
	return isBorder
}

func main() {
	screen := render.Screen{
		Width:  1500,
		Height: 750,
	}

	projection := render.Project(screen, render.Equirectangular{})

	spheres := []*geodesic.Geodesic{geodesic.Dodecahedron()}

	for frame := 0; frame < 120; frame++ {
		fmt.Printf("Frame %03d\n", frame)
		if frame % 24 == 0 {
			spheres = append(spheres, geodesic.Chamfer(spheres[frame / 24]))
		}

		weight := float64(frame % 24) / 24
		isBorder := drawBorders(spheres, projection, weight)

		img := image.NewRGBA(image.Rect(0, 0, screen.Width, screen.Height))
		for i, isBlack := range isBorder {
			x := i / projection.Width
			y := i % projection.Width
			if isBlack {
				img.Set(y, x, color.Black)
			} else {
				img.Set(y, x, color.White)
			}
		}

		out, err := os.Create(fmt.Sprintf("renders/tiling/frame-%03d.png", frame))
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
}
