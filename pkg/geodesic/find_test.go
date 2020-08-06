package geodesic

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"math"
	"testing"
)

func TestFind_Dodecahedron(t *testing.T) {
	tcs := []struct {
		name     string
		v        Vector
		want     int
		wantDist float64
	}{
		{
			name:     "north pole is north pole",
			v:        Vector{X: 0, Y: 0, Z: 1},
			want:     0,
			wantDist: 0.0,
		},
		{
			name:     "south pole is south pole",
			v:        Vector{X: 0, Y: 0, Z: -1},
			want:     11,
			wantDist: 0.0,
		},
		{
			name:     "equator 0 degrees",
			v:        Vector{X: 1, Y: 0, Z: 0},
			want:     7,
			wantDist: 0.546,
		},
		{
			name:     "equator 90 degrees",
			v:        Vector{X: 0, Y: 1, Z: 0},
			want:     1,
			wantDist: 0.459,
		},
		{
			name:     "equator 180 degrees",
			v:        Vector{X: -1, Y: 0, Z: 0},
			want:     5,
			wantDist: 0.546,
		},
		{
			name:     "equator 270 degrees",
			v:        Vector{X: 0, Y: -1, Z: 0},
			want:     8,
			wantDist: 0.459,
		},
	}

	gs := []*Geodesic{Dodecahedron()}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := Find(gs, tc.v)

			if got != tc.want {
				t.Errorf("got Find(Dodecahedron, %v) = %d, want %d", tc.v, got, tc.want)
			}

			gotDist := math.Sqrt(DistSq(gs[len(gs)-1].Centers[got], tc.v))

			if diff := cmp.Diff(tc.wantDist, gotDist, cmpopts.EquateApprox(0.0, 0.001)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestFind_M2(t *testing.T) {
	tcs := []struct {
		name     string
		v        Vector
		want     int
		wantDist float64
	}{
		{
			name:     "north pole is north pole",
			v:        Vector{X: 0, Y: 0, Z: 1},
			want:     0,
			wantDist: 0.0,
		},
		{
			name:     "south pole is south pole",
			v:        Vector{X: 0, Y: 0, Z: -1},
			want:     11,
			wantDist: 0.0,
		},
		{
			name:     "equator 0 degrees",
			v:        Vector{X: 1, Y: 0, Z: 0},
			want:     25,
			wantDist: 0.0,
		},
		{
			name:     "equator 90 degrees",
			v:        Vector{X: 0, Y: 1, Z: 0},
			want:     23,
			wantDist: 0.312,
		},
		{
			name:     "equator 180 degrees",
			v:        Vector{X: -1, Y: 0, Z: 0},
			want:     30,
			wantDist: 0.0,
		},
		{
			name:     "equator 270 degrees",
			v:        Vector{X: 0, Y: -1, Z: 0},
			want:     28,
			wantDist: 0.312,
		},
		{
			name: "0.64,0.6,0.48",
			v: Vector{X: 0.64, Y: 0.6, Z: 0.48},
			want: 17,
			wantDist: 0.171,
		},
	}

	gs := make([]*Geodesic, 2)
	gs[0] = Dodecahedron()
	gs[1] = Chamfer(gs[0])

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := Find(gs, tc.v)

			if got != tc.want {
				t.Errorf("got Find(Dodecahedron, %v) = %d, want %d", tc.v, got, tc.want)
			}

			gotDist := math.Sqrt(DistSq(gs[len(gs)-1].Centers[got], tc.v))

			if diff := cmp.Diff(tc.wantDist, gotDist, cmpopts.EquateApprox(0.0, 0.001)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestFind_M3(t *testing.T) {
	tcs := []struct {
		name     string
		v        Vector
		want     int
		wantDist float64
	}{
		{
			name:     "north pole is north pole",
			v:        Vector{X: 0, Y: 0, Z: 1},
			want:     0,
			wantDist: 0.0,
		},
		{
			name:     "south pole is south pole",
			v:        Vector{X: 0, Y: 0, Z: -1},
			want:     11,
			wantDist: 0.0,
		},
		{
			name:     "equator 0 degrees",
			v:        Vector{X: 1, Y: 0, Z: 0},
			want:     25,
			wantDist: 0.0,
		},
		{
			name:     "equator 90 degrees",
			v:        Vector{X: 0, Y: 1, Z: 0},
			want:     113,
			wantDist: 0.0,
		},
		{
			name:     "equator 180 degrees",
			v:        Vector{X: -1, Y: 0, Z: 0},
			want:     30,
			wantDist: 0.0,
		},
		{
			name:     "equator 270 degrees",
			v:        Vector{X: 0, Y: -1, Z: 0},
			want:     90,
			wantDist: 0.0,
		},
		{
			name: "0.64,0.6,0.48",
			v: Vector{X: 0.64, Y: 0.6, Z: 0.48},
			want: 72,
			wantDist: 0.119,
		},
	}

	gs := make([]*Geodesic, 3)
	gs[0] = Dodecahedron()
	gs[1] = Chamfer(gs[0])
	gs[2] = Chamfer(gs[1])

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := Find(gs, tc.v)

			if got != tc.want {
				t.Errorf("got Find(Dodecahedron, %v) = %d, want %d", tc.v, got, tc.want)
			}

			gotDist := math.Sqrt(DistSq(gs[len(gs)-1].Centers[got], tc.v))

			if diff := cmp.Diff(tc.wantDist, gotDist, cmpopts.EquateApprox(0.0, 0.001)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
