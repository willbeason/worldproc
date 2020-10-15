package diffeq

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"math"
	"testing"
)

func TestApproximate_Stability(t *testing.T) {
	x, y, v := 1.0, -1.0, 0.0

	// n is the number of steps each quarter-oscillation is broken into.
	n := 2
	h := math.Pi / (2 * math.Sqrt(2) * float64(n))
	for i := 0; i < n*4*5; i++ {
		// Ensure that over 5 oscillations, an accumulated energy offset of
		// less than 1%.
		x, y, v = RK4(h, x, y, v)
		energy := x*x + y*y + v*v
		if diff := cmp.Diff(2.0, energy, cmpopts.EquateApprox(0.01, 0.0)); diff != "" {
			t.Error(i, diff)
		}
	}
}

func TestApproximate(t *testing.T) {
	tcs := []struct {
		name   string
		method Approximation
		steps  int
		errX   float64
		errV   float64
	}{
		{
			name:   "Euler 10",
			method: Euler,
			steps:  10,
			errX:   0.037,
			errV:   0.0906,
		},
		{
			name:   "Euler 20",
			method: Euler,
			steps:  20,
			errX:   0.019,
			errV:   0.0445,
		},
		{
			name:   "Trapezoid 2",
			method: Trapezoid,
			steps:  2,
			errX:   0.1051,
			errV:   0.0701,
		},
		{
			name:   "Trapezoid 5",
			method: Trapezoid,
			steps:  5,
			errX:   0.0146,
			errV:   0.004400,
		},
		{
			name:   "Trapezoid 10",
			method: Trapezoid,
			steps:  10,
			errX:   0.00346,
			errV:   0.000545,
		},
		{
			name:   "Trapezoid 20",
			method: Trapezoid,
			steps:  20,
			errX:   0.00083,
			errV:   0.0000677,
		},
		{
			name:   "RK4 2",
			method: RK4,
			steps:  2,
			errX:   0.002118,
			errV:   0.000823,
		},
		{
			name:   "RK4 5",
			method: RK4,
			steps:  5,
			errX:   0.0001185,
			errV:   0.0000233,
		},
		{
			name:   "RK4 10",
			method: RK4,
			steps:  10,
			errX:   0.00001411,
			errV:   0.00000153,
		},
		{
			name:   "RK4 20",
			method: RK4,
			steps:  20,
			errX:   0.000001723,
			errV:   0.000000098,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			x, v := Approximate(tc.steps, tc.method)
			errX := math.Abs(x)
			errV := math.Abs(v - math.Sqrt(2))

			if diff := cmp.Diff(tc.errX, errX, cmpopts.EquateApprox(0.01, 0.0)); diff != "" {
				t.Error(diff)
			}
			if diff := cmp.Diff(tc.errV, errV, cmpopts.EquateApprox(0.01, 0.0)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func BenchmarkApproximate_Trapezoid(b *testing.B) {
	h := 0.01
	x, y, v := 1.0, -1.0, 0.0

	for i := 0; i < b.N; i++ {
		x, y, v = Trapezoid(h, x, y, v)
	}
	if x > 1E10 {
		b.Log(x, y, v)
	}
}

func BenchmarkApproximate_RK4(b *testing.B) {
	h := 0.01
	x, y, v := 1.0, -1.0, 0.0

	for i := 0; i < b.N; i++ {
		x, y, v = RK4(h, x, y, v)
	}
}
