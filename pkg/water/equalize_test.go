package water

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/willbeason/hydrology/pkg/geodesic"
	"math/rand"
	"testing"
)

func TestEqualize(t *testing.T) {
	tcs := []struct{
		name string
		heights []float64
		water []float64
		want []float64
	} {
		{
			name: "single tile",
			heights: []float64{0.0},
			water: []float64{1.0},
			want: []float64{1.0},
		},
		{
			name: "two even tiles",
			heights: []float64{0.0, 0.0},
			water: []float64{1.5, 0.5},
			want: []float64{1.0, 1.0},
		},
		{
			name: "incline",
			heights: []float64{0.0, 0.1, 0.2, 0.3},
			water: []float64{1.0, 1.0, 1.0, 1.0},
			want: []float64{1.15, 1.05, 0.95, 0.85},
		},
		{
			name: "decline",
			heights: []float64{0.3, 0.2, 0.1, 0.0},
			water: []float64{1.0, 1.0, 1.0, 1.0},
			want: []float64{0.85, 0.95, 1.05, 1.15},
		},
		{
			name: "unfallen water",
			heights: []float64{0.3, 0.2, 0.1, 0.0},
			water: []float64{1.0, 0.0, 0.0, 0.0},
			want: []float64{0.1, 0.2, 0.3, 0.4},
		},
		{
			name: "two adjacent lakes",
			heights: []float64{0.0, 0.5, 1.0, 0.5, 0.0},
			water: []float64{0.5, 0.5, 0.0, 0.5, 0.5},
			want: []float64{0.75, 0.25, 0.0, 0.25, 0.75},
		},
		{
			name: "two separate lakes",
			heights: []float64{0.0, 0.5, 1.0, 1.5, 1.0, 0.5, 0.0},
			water: []float64{0.5, 0.5, 0.0, 0.0, 0.0, 0.5, 0.5},
			want: []float64{0.75, 0.25, 0.0, 0.0, 0.0, 0.25, 0.75},
		},
		{
			name: "linked lakes",
			heights: []float64{0.0, 0.3, 0.5, 0.4, 1.0},
			water: []float64{0.5, 0.5, 0.5, 0.5, 0.0},
			want: []float64{0.8, 0.5, 0.3, 0.4, 0.0},
		},
		{
			name: "upper lake",
			heights: []float64{0.0, 0.3, 0.5, 0.4, 1.0},
			water: []float64{0.1, 0.1, 0.1, 0.2, 0.0},
			want: []float64{0.35, 0.05, 0.0, 0.1, 0.0},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			g := &geodesic.Geodesic{
				Faces:   make([]geodesic.Node, len(tc.heights)),
				Edges: map[geodesic.Edge]int{},
			}
			for i := 0; i < len(tc.heights) - 1; i++ {
				g.Link(i, i+1)
			}

			Equalize(tc.water, tc.heights, g)

			if diff := cmp.Diff(tc.want, tc.water, cmpopts.EquateApprox(0.0, 0.001)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestLake_Equalize(t *testing.T) {
	tcs := []struct{
		name string
		cells []IndexHeight
		want []IndexHeight
	} {
		{
			name: "basin",
			cells: []IndexHeight{
				{Height: 0.0, Water: 0.1},
				{Height: 1.0, Water: 0.1},
				{Height: 1.0, Water: 0.1},
				{Height: 1.0, Water: 0.1},
				{Height: 1.0, Water: 0.1},
				{Height: 1.0, Water: 0.1},
				{Height: 1.0, Water: 0.1},
			},
			want: []IndexHeight{
				{Height: 0.0, Water: 0.7},
				{Height: 1.0},
				{Height: 1.0},
				{Height: 1.0},
				{Height: 1.0},
				{Height: 1.0},
				{Height: 1.0},
			},
		},
		{
			name: "small basin",
			cells: []IndexHeight{
				{Height: 0.0, Water: 0.1},
				{Height: 0.5, Water: 0.1},
				{Height: 1.0, Water: 0.1},
				{Height: 1.0, Water: 0.1},
				{Height: 1.0, Water: 0.1},
				{Height: 1.0, Water: 0.1},
				{Height: 1.0, Water: 0.1},
			},
			want: []IndexHeight{
				{Height: 0.0, Water: 0.6},
				{Height: 0.5, Water: 0.1},
				{Height: 1.0},
				{Height: 1.0},
				{Height: 1.0},
				{Height: 1.0},
				{Height: 1.0},
			},
		},
		{
			name: "incline",
			cells: []IndexHeight{
				{Height: 0.0, Water: 0.1},
				{Height: 0.1, Water: 0.1},
				{Height: 0.2, Water: 0.1},
				{Height: 0.3, Water: 0.1},
				{Height: 0.4, Water: 0.1},
				{Height: 0.5, Water: 0.1},
				{Height: 0.6, Water: 0.1},
			},
			want: []IndexHeight{
				{Height: 0.0, Water: 0.325},
				{Height: 0.1, Water: 0.225},
				{Height: 0.2, Water: 0.125},
				{Height: 0.3, Water: 0.025},
				{Height: 0.4},
				{Height: 0.5},
				{Height: 0.6},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			l := &Lake{}
			rand.Shuffle(len(tc.cells), func(i, j int) {
				tc.cells[i], tc.cells[j] = tc.cells[j], tc.cells[i]
			})

			for i, c := range tc.cells {
				l.Add(i, c.Height, c.Water)
			}

			l.Equalize()

			if diff := cmp.Diff(tc.want, l.IndexHeights,
				cmpopts.EquateApprox(0.0, 0.001), cmpopts.IgnoreFields(IndexHeight{}, "Index")); diff != "" {
				t.Error(diff)
			}
		})
	}
}
