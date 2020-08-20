package climate

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/willbeason/hydrology/pkg/geodesic"
	"math"
	"testing"
)

func TestCoriolisDeflect(t *testing.T) {
	tcs := []struct {
		start geodesic.Vector
		want  geodesic.Vector
	}{
		// Zero vector.
		{
			start: geodesic.Vector{X: 0, Y: 0, Z: 0},
			want:  geodesic.Vector{X: 0, Y: 0, Z: 0},
		},
		// No deflection.
		{
			start: geodesic.Vector{X: 0, Y: 0, Z: 1},
			want:  geodesic.Vector{X: 0, Y: 0, Z: 1},
		},
		{
			start: geodesic.Vector{X: 0, Y: 0, Z: -1},
			want:  geodesic.Vector{X: 0, Y: 0, Z: -1},
		},
		// Unit deflection.
		{
			start: geodesic.Vector{X: 1, Y: 0, Z: 0},
			want:  geodesic.Vector{X: 1, Y: -1.454, Z: 0},
		},
		{
			start: geodesic.Vector{X: math.Sqrt(0.5), Y: math.Sqrt(0.5), Z: 0},
			want:  geodesic.Vector{X: 1.736, Y: -0.321, Z: 0},
		},
		{
			start: geodesic.Vector{X: 0, Y: 1, Z: 0},
			want:  geodesic.Vector{X: 1.454, Y: 1, Z: 0},
		},
		{
			start: geodesic.Vector{X: -math.Sqrt(0.5), Y: math.Sqrt(0.5), Z: 0},
			want:  geodesic.Vector{X: 0.321, Y: 1.736, Z: 0},
		},
		{
			start: geodesic.Vector{X: -1, Y: 0, Z: 0},
			want:  geodesic.Vector{X: -1, Y: 1.454, Z: 0},
		},
		{
			start: geodesic.Vector{X: -math.Sqrt(0.5), Y: -math.Sqrt(0.5), Z: 0},
			want:  geodesic.Vector{X: -1.736, Y: 0.321, Z: 0},
		},
		{
			start: geodesic.Vector{X: 0, Y: -1, Z: 0},
			want:  geodesic.Vector{X: -1.454, Y: -1, Z: 0},
		},
		{
			start: geodesic.Vector{X: math.Sqrt(0.5), Y: -math.Sqrt(0.5), Z: 0},
			want:  geodesic.Vector{X: -0.321, Y: -1.736, Z: 0},
		},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprint(tc.start), func(t *testing.T) {
			got := CoriolisDeflect(tc.start)

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateApprox(0, 0.001)); diff != "" {
				t.Error("direction", diff)
			}
		})
	}

	if diff := cmp.Diff(1.0, invCoriolisDrag, cmpopts.EquateApprox(0, 1E-7)); diff != "" {
		// Lets us neglect this term and avoid an unnecessary multiplication.
		t.Error("invCoriolisDrag is not nearly 1", diff)
	}
}
