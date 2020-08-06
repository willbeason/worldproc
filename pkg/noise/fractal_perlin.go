package noise

import (
	"github.com/willbeason/hydrology/pkg/geodesic"
	"math/rand"
)

type PerlinFractal struct {
	Perlin
	Depth int

	Scale float64
	InvScale float64
	Offset geodesic.Vector
}

func NewPerlinFractal(dim int, depth int, scale float64) *PerlinFractal {
	return &PerlinFractal{
		Perlin:   *NewPerlin(dim),
		Depth:    depth,
		Scale:    scale,
		InvScale: 1.0 / scale,
		Offset:   geodesic.Vector{
			X: float64(dim) * rand.Float64(),
			Y: float64(dim) * rand.Float64(),
			Z: float64(dim) * rand.Float64(),
		},
	}
}

func (p *PerlinFractal) ValueAt(v geodesic.Vector) float64 {
	result := p.Perlin.ValueAt(v)

	cScale := 1.0
	for i := 0; i < p.Depth; i++ {
		cScale *= p.Scale

		v = v.Add(geodesic.Vector{X: 2, Y: 2, Z: 2})
		v = v.Scale(p.InvScale)

		result += p.Perlin.ValueAt(v) * cScale
	}

	return result
}
