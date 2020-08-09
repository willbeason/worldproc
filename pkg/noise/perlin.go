package noise

import (
	"github.com/willbeason/hydrology/pkg/geodesic"
	"math"
	"math/rand"
)

type Perlin struct {
	Dim int
	Dim2 int
	Noise []geodesic.Vector

}

func NewPerlin(r *rand.Rand, dim int) *Perlin {
	dim2 := dim*dim
	result := &Perlin{
		Dim: dim,
		Dim2: dim2,
		Noise: make([]geodesic.Vector, dim*dim*dim),
	}

	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			for k := 0; k < dim; k++ {
				angle := geodesic.Angle{
					Theta: 2.0*math.Pi*r.Float64(),
					Phi:   math.Acos(2.0*r.Float64() - 1.0),
				}
				result.Noise[i*dim2 + j*dim + k] = angle.Vector()
			}
		}
	}
	return result
}

func modf(f float64) (int, float64) {
	q, r := math.Modf(f)
	return int(q), r
}

func (p *Perlin) noiseAt(x, y, z int) geodesic.Vector {
	idx := x*p.Dim2 + y*p.Dim + z
	return p.Noise[idx]
}

func (p *Perlin) diffDotNoiseAt(xr, yr, zr float64, x, y, z int) float64 {
	noiseXYZ := p.noiseAt(x, y, z)

	return xr*noiseXYZ.X + yr*noiseXYZ.Y + zr*noiseXYZ.Z
}

func lerp(a0, a1, w, wc float64) float64 {
	return wc*a0 + w*a1
}

func (p *Perlin) ValueAt(v geodesic.Vector) float64 {
	//v = v.Add(geodesic.Vector{0.5, 0.5, 0.5})

	x0, xr := modf(v.X)
	if xr < 0 {
		x0--
		xr += 1
	}
	y0, yr := modf(v.Y)
	if yr < 0 {
		y0--
		yr += 1
	}
	z0, zr := modf(v.Z)
	if zr < 0 {
		z0--
		zr += 1
	}

	x0 %= p.Dim
	if x0 < 0 {
		x0 += p.Dim
	}
	y0 %= p.Dim
	if y0 < 0 {
		y0 += p.Dim
	}
	z0 %= p.Dim
	if z0 < 0 {
		z0 += p.Dim
	}

	x1 := (x0 + 1) % p.Dim
	y1 := (y0 + 1) % p.Dim
	z1 := (z0 + 1) % p.Dim

	xc := 1-xr
	yc := 1-yr
	zc := 1-zr

	noise000 := p.diffDotNoiseAt(xr, yr, zr, x0, y0, z0)
	noise001 := p.diffDotNoiseAt(xr, yr, -zc, x0, y0, z1)
	noise010 := p.diffDotNoiseAt(xr, -yc, zr, x0, y1, z0)
	noise011 := p.diffDotNoiseAt(xr, -yc, -zc, x0, y1, z1)
	noise100 := p.diffDotNoiseAt(-xc, yr, zr, x1, y0, z0)
	noise101 := p.diffDotNoiseAt(-xc, yr, -zc, x1, y0, z1)
	noise110 := p.diffDotNoiseAt(-xc, -yc, zr, x1, y1, z0)
	noise111 := p.diffDotNoiseAt(-xc, -yc, -zc, x1, y1, z1)

	// Linearly interpolate noise.
	noise00 := lerp(noise000, noise001, zr, zc)
	noise01 := lerp(noise010, noise011, zr, zc)
	noise10 := lerp(noise100, noise101, zr, zc)
	noise11 := lerp(noise110, noise111, zr, zc)

	noise0 := lerp(noise00, noise01, yr, yc)
	noise1 := lerp(noise10, noise11, yr, yc)

	return lerp(noise0, noise1, xr, xc)
}
