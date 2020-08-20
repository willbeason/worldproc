package geodesic

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestDodecahedron(t *testing.T) {
	tcs := []struct{
		id int
		want []int
	}{
		{
			id: 0,
			want: []int{1, 5, 4, 3, 2},
		},
		{
			id: 1,
			want: []int{0, 2, 6, 10, 5},
		},
		{
			id: 2,
			want: []int{0, 3, 7, 6, 1},
		},
		{
			id: 3,
			want: []int{0, 4, 8, 7, 2},
		},
		{
			id: 4,
			want: []int{0, 5, 9, 8, 3},
		},
		{
			id: 5,
			want: []int{0, 1, 10, 9, 4},
		},
		{
			id: 6,
			want: []int{1, 2, 7, 11, 10},
		},
		{
			id: 7,
			want: []int{2, 3, 8, 11, 6},
		},
		{
			id: 8,
			want: []int{3, 4, 9, 11, 7},
		},
		{
			id: 9,
			want: []int{4, 5, 10, 11, 8},
		},
		{
			id: 10,
			want: []int{1, 6, 11, 9, 5},
		},
		{
			id: 11,
			want: []int{6, 7, 8, 9, 10},
		},
	}

	d := Dodecahedron()

	for _, tc := range tcs {
		if diff := cmp.Diff(tc.want, d.Faces[tc.id].Neighbors); diff != "" {
			t.Error(diff)
		}
	}
}

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
