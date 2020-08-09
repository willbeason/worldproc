package main

import (
	"fmt"
	"github.com/willbeason/hydrology/pkg/geodesic"
	"github.com/willbeason/hydrology/pkg/noise"
	"github.com/willbeason/hydrology/pkg/planet"
	"github.com/willbeason/hydrology/pkg/render"
	"github.com/willbeason/hydrology/pkg/water"
	"image"
	"math"
	"math/rand"
	"sort"
	"time"
)

func main() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	screen := render.Screen{
		Width:  2160,
		Height: 1080,
	}

	size := 6
	spheres := geodesic.New(size, false)

	depth := 30
	perlinNoise := noise.NewPerlinFractal(10, depth, 0.8)

	sphere := spheres[len(spheres)-1]

	// heights by cell
	p := &planet.Planet{
		Size:    size,
		Heights: make([]float64, len(sphere.Centers)),
		Waters:  make([]float64, len(sphere.Centers)),
		Flows:   make([]float64, len(sphere.Centers)),
	}

	for cell, pos := range sphere.Centers {
		p.Heights[cell] = perlinNoise.ValueAt(pos)
	}

	oceanWater := 0.0
	sortedHeights := make([]float64, len(p.Heights))
	copy(sortedHeights, p.Heights)
	sort.Float64s(sortedHeights)

	halfway := sortedHeights[len(sortedHeights)/2]
	for _, h := range p.Heights {
		if h < halfway {
			oceanWater += halfway - h
		}
	}
	avgWater := oceanWater / float64(len(sortedHeights)/2)
	fmt.Println(halfway, avgWater)

	projection := render.Project(screen, render.Equirectangular{})

	quanta := 0.01
	iters := int(avgWater / quanta)
	fmt.Println("Total Iters:", iters)

	renderImg(seed, projection, spheres, p.Heights, p.Waters, p.Flows, 0)
	for iter := 1; iter < iters; iter++ {
		fmt.Print(iter, "...", "Raining")
		water.Rain(quanta, p.Waters, p.Heights, p.Flows, sphere)

		if iter%5 == 0 {
			fmt.Print("...", "Equalizing")
			water.Equalize(p.Waters, p.Heights, sphere)
			renderImg(seed, projection, spheres, p.Heights, p.Waters, p.Flows, iter)
		}

		fmt.Println()
	}

	water.Equalize(p.Waters, p.Heights, sphere)
	renderImg(seed, projection, spheres, p.Heights, p.Waters, p.Flows, iters)


}

func renderImg(seed int64, projection render.Projection, spheres []*geodesic.Geodesic, heights, waters, flow []float64, id int) {
	screen := projection.Screen
	img := image.NewRGBA(image.Rect(0, 0, screen.Width, screen.Height))

	pxWaterHeights := make([]float64, screen.Width*screen.Height)
	pxLandHeights := make([]float64, screen.Width*screen.Height)

	sphere := spheres[len(spheres)-1]
	for x := 0; x < screen.Width; x++ {
		for y := 0; y < screen.Height; y++ {
			pidx := y*screen.Width + x
			angle := projection.Pixels[pidx]
			v := angle.Vector()
			idx := geodesic.Find(spheres, v)
			dist := math.Sqrt(geodesic.DistSq(v, sphere.Centers[idx]))

			pxW1 := waters[idx] + flow[idx]/2000.0
			pxH1 := heights[idx]

			// Linearly interpolate the cell's stats with the second-closest cell.
			idx2 := 0
			distSq2 := math.MaxFloat64
			for _, n := range sphere.Faces[idx].Neighbors {
				nDistSq2 := geodesic.DistSq(v, sphere.Centers[n])
				if nDistSq2 < distSq2 {
					idx2 = n
					distSq2 = nDistSq2
				}
			}
			dist2 := math.Sqrt(distSq2)
			pxW2 := waters[idx2] + flow[idx2]/2000.0
			pxH2 := heights[idx2]

			pxWaterHeights[pidx] = render.Lerp(pxW1, pxW2, dist/(dist+dist2))
			pxLandHeights[pidx] = render.Lerp(pxH1, pxH2, dist/(dist+dist2))
		}
	}

	screen.PaintLandWater(pxLandHeights, pxWaterHeights, img)

	render.WriteImage(img, fmt.Sprintf("renders/hydro-%d-%d.png", seed, id))
}
