package main

import (
	"fmt"
	"github.com/willbeason/worldproc/pkg/geodesic"
	"github.com/willbeason/worldproc/pkg/noise"
	"github.com/willbeason/worldproc/pkg/render"
	"github.com/willbeason/worldproc/pkg/water"
	"image"
	"image/color"
	"sort"
)

func main() {
	screen := render.Screen{
		Width:  1500,
		Height: 750,
	}

	projection := render.Project(screen, render.Equirectangular{})

	nSpheres := 7
	spheres := []*geodesic.Geodesic{geodesic.Dodecahedron()}
	for i := 0; i < nSpheres; i++ {
		spheres = append(spheres, geodesic.Chamfer(spheres[i]))
	}

	depth := 20
	perlinNoise := noise.NewPerlinFractal(10, depth, 0.8)


	sphere := spheres[nSpheres]

	// heights by cell
	heights := make([]float64, len(sphere.Centers))
	waters := make([]float64, len(sphere.Centers))

	//water := make([]float64, len(sphere.Centers))

	for cell, pos := range sphere.Centers {
		heights[cell] = perlinNoise.ValueAt(pos)
	}

	sortedHeights := make([]float64, len(sphere.Centers))
	copy(sortedHeights, heights)
	sort.Slice(sortedHeights, func(i, j int) bool {
		return sortedHeights[i] < sortedHeights[j]
	})
	seaLevel := sortedHeights[len(sphere.Centers)/2]
	oceanWater := 0.0
	for _, h := range sortedHeights {
		if h < seaLevel {
			oceanWater += seaLevel - h
		}
	}
	avgWater := oceanWater / float64(len(sphere.Centers) / 2)
	fmt.Println(avgWater)

	iters := 4000
	previous := make([]float64, len(waters))
	order := make([]water.Order, len(waters))

	for i := range waters {
		// Rain avgWater units evenly.
		waters[i] += avgWater
	}

	for iter := 0; iter < iters; iter++ {
		water.EqualizeIter(heights, waters, order, sphere)
		if iter % 100 == 98 {
			copy(previous, waters)
		}
		if iter % 100 == 99 {
			diffsq := 0.0
			for i, w := range waters {
				diffsq += (w - previous[i])*(w - previous[i])
			}
			fmt.Println("Iter:", iter, "Diff:", diffsq)
		}
	}

	sortedWaters := make([]float64, len(sphere.Centers))
	copy(sortedWaters, waters)
	sort.Float64s(sortedWaters)
	threshold := sortedWaters[int(float64(len(sortedWaters)) * 0.5)]
	threshold2 := sortedWaters[int(float64(len(sortedWaters)) * 0.3)]

	fmt.Println(sortedWaters[0], threshold2, threshold, sortedWaters[len(sortedWaters)-1])

	csw := render.NewColorScale(
		[]render.ColorPoint{
			{threshold2, color.RGBA{255,255, 255, 255}},
			{threshold, color.RGBA{0, 0, 0, 255}},
		})

	pxWaterHeigts := make([]float64, screen.Width*screen.Height)
	pxLandHeights := make([]float64, screen.Width*screen.Height)
	for x := 0; x < screen.Width; x++ {
		for y := 0; y < screen.Height; y++ {
			pidx := y*screen.Width + x
			angle := projection.Pixels[pidx]
			idx := geodesic.Find(spheres, angle.Vector())
			pxWaterHeigts[pidx] = waters[idx]
			pxLandHeights[pidx] = heights[idx] / 2.6
		}
	}
	img := image.NewRGBA(image.Rect(0, 0, screen.Width, screen.Height))
	screen.Paint(pxWaterHeigts, csw, img)

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

	img2 := image.NewRGBA(image.Rect(0, 0, screen.Width, screen.Height))
	screen.PaintAndShade(pxLandHeights, csl, img2)

	render.WriteImage(img, fmt.Sprintf("hydro-%d-e.png", iters))
	render.WriteImage(img2, fmt.Sprintf("hydro-%d-f.png", iters))
}
