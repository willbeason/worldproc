package climate

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/willbeason/hydrology/pkg/geodesic"
	"math"
	"testing"
)

func TestGradient(t *testing.T) {
	tcs := []struct {
		name      string
		pressures []float64
		// want is the direction and magnitude of the pressure gradient.
		// Recall that this points towards higher pressure, and is the opposite
		// the direction of acceleration.
		want geodesic.Vector
	}{
		{
			name: "all zero",
			pressures: []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			want: geodesic.Vector{X: 0.0, Y: 0.0},
		},
		{
			name: "all ones",
			pressures: []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0},
			want: geodesic.Vector{X: 0.0, Y: 0.0},
		},
		{
			name: "symmetric",
			pressures: []float64{1.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0},
			want: geodesic.Vector{X: 0.0, Y: 0.0},
		},
		{
			name: "positive X",
			pressures: []float64{0.0,
				1.0, math.Cos(math.Pi/3), math.Cos(2*math.Pi/3),
				-1.0, math.Cos(4*math.Pi/3), math.Cos(5*math.Pi/3)},
			want: geodesic.Vector{X: 1.0, Y: 0.0},
		},
		{
			name: "positive Y",
			pressures: []float64{0.0,
				0.0, math.Sin(math.Pi/3), math.Sin(2*math.Pi/3),
				0.0, math.Sin(4*math.Pi/3), math.Sin(5*math.Pi/3)},
			want: geodesic.Vector{X: 0.0, Y: 1.0},
		},
	}

	sphere := &geodesic.Geodesic{
		Centers: []geodesic.Vector{
			{X: 0, Y: 0, Z: 0},
			geodesic.Angle{Theta: 0, Phi: 0 * math.Pi / 3}.Vector(),
			geodesic.Angle{Theta: 0, Phi: 1 * math.Pi / 3}.Vector(),
			geodesic.Angle{Theta: 0, Phi: 2 * math.Pi / 3}.Vector(),
			geodesic.Angle{Theta: 0, Phi: 3 * math.Pi / 3}.Vector(),
			geodesic.Angle{Theta: 0, Phi: 4 * math.Pi / 3}.Vector(),
			geodesic.Angle{Theta: 0, Phi: 5 * math.Pi / 3}.Vector(),
		},
		Faces: make([]geodesic.Node, 7),
		Edges: map[geodesic.Edge]int{},
	}
	sphere.Link(0, 1)
	sphere.Link(0, 2)
	sphere.Link(0, 3)
	sphere.Link(0, 4)
	sphere.Link(0, 5)
	sphere.Link(0, 6)

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := Gradient(0, tc.pressures, sphere)

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateApprox(0.0, 0.001)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestDivergence6(t *testing.T) {
	tcs := []struct{
		name string
		vectors []geodesic.Vector
		want float64
	} {
		{
			name: "divergence 2.0",
			vectors: []geodesic.Vector{
				{X: 0.0, Y: 0.0, Z: 0.0},
				geodesic.Angle{Theta: 0, Phi: 0 * math.Pi / 3}.Vector(),
				geodesic.Angle{Theta: 0, Phi: 1 * math.Pi / 3}.Vector(),
				geodesic.Angle{Theta: 0, Phi: 2 * math.Pi / 3}.Vector(),
				geodesic.Angle{Theta: 0, Phi: 3 * math.Pi / 3}.Vector(),
				geodesic.Angle{Theta: 0, Phi: 4 * math.Pi / 3}.Vector(),
				geodesic.Angle{Theta: 0, Phi: 5 * math.Pi / 3}.Vector(),
			},
			want: 2.0,
		},
	}

	sphere := &geodesic.Geodesic{
		Centers: []geodesic.Vector{
			{X: 0, Y: 0, Z: 0},
			geodesic.Angle{Theta: 0, Phi: 0 * math.Pi / 3}.Vector(),
			geodesic.Angle{Theta: 0, Phi: 1 * math.Pi / 3}.Vector(),
			geodesic.Angle{Theta: 0, Phi: 2 * math.Pi / 3}.Vector(),
			geodesic.Angle{Theta: 0, Phi: 3 * math.Pi / 3}.Vector(),
			geodesic.Angle{Theta: 0, Phi: 4 * math.Pi / 3}.Vector(),
			geodesic.Angle{Theta: 0, Phi: 5 * math.Pi / 3}.Vector(),
		},
		Faces: make([]geodesic.Node, 7),
		Edges: map[geodesic.Edge]int{},
	}
	sphere.Link(0, 1)
	sphere.Link(0, 2)
	sphere.Link(0, 3)
	sphere.Link(0, 4)
	sphere.Link(0, 5)
	sphere.Link(0, 6)


	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := Divergence(0, tc.vectors, sphere)

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateApprox(0.0, 0.001)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestDivergence5(t *testing.T) {
	tcs := []struct{
		name string
		vectors []geodesic.Vector
		want float64
	} {
		{
			name: "divergence 2.0",
			vectors: []geodesic.Vector{
				{X: 0.0, Y: 0.0, Z: 0.0},
				geodesic.Angle{Theta: 0, Phi: 0 * math.Pi / 5}.Vector(),
				geodesic.Angle{Theta: 0, Phi: 2 * math.Pi / 5}.Vector(),
				geodesic.Angle{Theta: 0, Phi: 4 * math.Pi / 5}.Vector(),
				geodesic.Angle{Theta: 0, Phi: 6 * math.Pi / 5}.Vector(),
				geodesic.Angle{Theta: 0, Phi: 8 * math.Pi / 5}.Vector(),
			},
			want: 2.0,
		},
	}

	sphere := &geodesic.Geodesic{
		Centers: []geodesic.Vector{
			{X: 0, Y: 0, Z: 0},
			geodesic.Angle{Theta: 0, Phi: 0 * math.Pi / 5}.Vector(),
			geodesic.Angle{Theta: 0, Phi: 2 * math.Pi / 5}.Vector(),
			geodesic.Angle{Theta: 0, Phi: 4 * math.Pi / 5}.Vector(),
			geodesic.Angle{Theta: 0, Phi: 6 * math.Pi / 5}.Vector(),
			geodesic.Angle{Theta: 0, Phi: 8 * math.Pi / 5}.Vector(),
		},
		Faces: make([]geodesic.Node, 6),
		Edges: map[geodesic.Edge]int{},
	}
	sphere.Link(0, 1)
	sphere.Link(0, 2)
	sphere.Link(0, 3)
	sphere.Link(0, 4)
	sphere.Link(0, 5)

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := Divergence(0, tc.vectors, sphere)

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateApprox(0.0, 0.001)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestFlow(t *testing.T) {
	tcs := []struct {
		name      string
		pressures []float64
		converge  int
	}{
		{
			name: "all zero",
			pressures: []float64{
				1.2,
				1.1,
				1.0,
				0.9,
				0.8,
				0.7,
				0.6,
				0.5,
				0.4,
				0.3,
				0.2,
				0.1,
			},
			converge: 250,
		},
		//{
		//	name: "problematic",
		//	pressures: []float64{
		//		0.7396,
		//		0.6609,
		//		0.6446,
		//		0.5908,
		//		0.6969,
		//		0.5249,
		//		0.5720,
		//		0.7002,
		//		0.5620,
		//		0.6613,
		//		0.6701,
		//		0.6465,
		//	},
		//	converge: 10,
		//},
	}

	sphere := geodesic.Dodecahedron()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			climates := make([]Climate, len(tc.pressures))
			average := 0.0
			totalEnergy := 0.0
			for i, p := range tc.pressures {
				climates[i].Air = p
				climates[i].SetTemperature(ZeroCelsius)
				average += p / float64(len(tc.pressures))
				totalEnergy += climates[i].AirEnergy
			}

			want := make([]float64, len(tc.pressures))
			for i := range want {
				want[i] = average
			}
			//lastDiff := average * float64(len(want))
			for i := 0; i < tc.converge; i++ {
				Flow(climates, sphere, 2.0)
				deltaAir := 0.0
				deltaEnergy := totalEnergy
				for _, c := range climates {
					deltaAir += c.Pressure() - average
					deltaEnergy -= c.AirEnergy
				}
				if diff := cmp.Diff(0.0, deltaEnergy, cmpopts.EquateApprox(0.0, 0.001)); diff != "" {
					jsn, _ := json.MarshalIndent(climates, "", "  ")
					t.Log(string(jsn))
					t.Fatal(diff)
				}

				got := make([]float64, len(tc.pressures))
				for i := range got {
					got[i] = climates[i].Pressure()
				}
				newDiff := sumDiff(want, got)
				//if newDiff > lastDiff + 0.02 {
				//	jsn, _ := json.MarshalIndent(climates, "", "  ")
				//	t.Log(string(jsn))
				//	t.Log(cmp.Diff(want, got, cmpopts.EquateApprox(0.0, 0.05)))
				//	t.Fatalf("%d: %.04f > %.04f", i, newDiff, lastDiff)
				//}
				t.Log(i, newDiff)
				//lastDiff = newDiff
			}

			got := make([]float64, len(tc.pressures))
			for i := range got {
				got[i] = climates[i].Pressure()
			}

			if diff := cmp.Diff(want, got, cmpopts.EquateApprox(0.0, 0.01)); diff != "" {
				jsn, _ := json.MarshalIndent(climates, "", "  ")
				t.Log(string(jsn))
				t.Log("sumDiff", sumDiff(want, got))

				t.Error(diff)
			}
		})
	}
}

func sumDiff(want, got []float64) float64 {
	result := 0.0
	for i, w := range want {
		result += math.Abs(w - got[i])
	}
	return result
}
