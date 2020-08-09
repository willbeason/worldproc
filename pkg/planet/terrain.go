package planet

import (
	"github.com/willbeason/hydrology/pkg/geodesic"
	"github.com/willbeason/hydrology/pkg/noise"
	"github.com/willbeason/hydrology/pkg/render"
	"image"
	"math"
)

func AddTerrain(p *Planet, sphere *geodesic.Geodesic, perlinNoise *noise.PerlinFractal) {
	p.Heights = make([]float64, len(sphere.Centers))
	for cell, pos := range sphere.Centers {
		p.Heights[cell] = perlinNoise.ValueAt(pos)
	}
}

func RenderTerrain(p *Planet, projection render.Projection, spheres []*geodesic.Geodesic) *image.RGBA {
	screen := projection.Screen
	img := image.NewRGBA(image.Rect(0, 0, screen.Width, screen.Height))

	pxWaterHeights := make([]float64, screen.Width*screen.Height)
	pxLandHeights := make([]float64, screen.Width*screen.Height)

	heights := p.Heights
	waters := p.Waters
	flow := p.Flows

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
	return img
}
