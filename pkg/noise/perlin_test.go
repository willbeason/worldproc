package noise

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/willbeason/worldproc/pkg/geodesic"
	"math"
	"testing"
)

func TestPerlin_ValueAt(t *testing.T) {
	p := Perlin{
		Dim:   10,
		Dim2:  100,
		Noise: make([]geodesic.Vector, 1000),
	}

	sqrt3 := math.Sqrt(3)
	p.Noise[999] = geodesic.Vector{
		X: sqrt3,
		Y: sqrt3,
		Z: sqrt3,
	}
	p.Noise[990] = geodesic.Vector{
		X: sqrt3,
		Y: sqrt3,
		Z: -sqrt3,
	}
	p.Noise[909] = geodesic.Vector{
		X: sqrt3,
		Y: -sqrt3,
		Z: sqrt3,
	}
	p.Noise[900] = geodesic.Vector{
		X: sqrt3,
		Y: -sqrt3,
		Z: -sqrt3,
	}
	p.Noise[99] = geodesic.Vector{
		X: -sqrt3,
		Y: sqrt3,
		Z: sqrt3,
	}
	p.Noise[90] = geodesic.Vector{
		X: -sqrt3,
		Y: sqrt3,
		Z: -sqrt3,
	}
	p.Noise[9] = geodesic.Vector{
		X: -sqrt3,
		Y: -sqrt3,
		Z: sqrt3,
	}
	p.Noise[0] = geodesic.Vector{
		X: -sqrt3,
		Y: -sqrt3,
		Z: -sqrt3,
	}

	want := 2.598

	got := p.ValueAt(geodesic.Vector{
		X: 9.5,
		Y: 9.5,
		Z: 9.5,
	})
	if diff := cmp.Diff(want, got, cmpopts.EquateApprox(0.0, 0.001)); diff != "" {
		t.Error(diff)
	}

	got = p.ValueAt(geodesic.Vector{
		X: -0.5,
		Y: -0.5,
		Z: -0.5,
	})
	if diff := cmp.Diff(want, got, cmpopts.EquateApprox(0.0, 0.001)); diff != "" {
		t.Error(diff)
	}
}
