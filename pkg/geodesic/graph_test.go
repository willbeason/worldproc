package geodesic

import (
	"testing"
)

func TestGoldberg(t *testing.T) {
	tcs := []struct {
		name       string
		iterations int
		wantFaces  int
		wantEdges  int
	}{
		{
			name:       "m = 1",
			iterations: 0,
			wantFaces:  12,
			wantEdges: 30,
		},
		{
			name:       "m = 2",
			iterations: 1,
			wantFaces:  42,
			wantEdges:  120,
		},
		{
			name:       "m = 4",
			iterations: 2,
			wantFaces:  162,
			wantEdges:  480,
		},
		{
			name:       "m = 8",
			iterations: 3,
			wantFaces:  642,
			wantEdges:  1920,
		},
		{
			name:       "m = 16",
			iterations: 4,
			wantFaces:  2562,
			wantEdges:  7680,
		},
		{
			name:       "m = 32",
			iterations: 5,
			wantFaces:  10242,
			wantEdges:  30720,
		},
		{
			name:       "m = 64",
			iterations: 6,
			wantFaces:  40962,
			wantEdges:  122880,
		},
		{
			name:       "m = 128",
			iterations: 7,
			wantFaces:  163842,
			wantEdges:  491520,
		},
		{
			name:       "m = 256",
			iterations: 8,
			wantFaces:  655362,
			wantEdges:  1966080,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if testing.Short() && tc.iterations > 3 {
				t.Skip()
			}
			g := Dodecahedron()

			for i := 0; i < tc.iterations; i++ {
				g = Chamfer(g)
			}

			if len(g.Faces) != tc.wantFaces {
				t.Errorf("got len(g.Faces) = %d, want %d", len(g.Faces), tc.wantFaces)
			}

			if len(g.Edges)/2 != tc.wantEdges {
				t.Errorf("got len(g.Edges)/2 = %d, want %d", len(g.Edges)/2, tc.wantEdges)
			}
		})
	}
}
