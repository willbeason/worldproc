package water

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"math/rand"
	"testing"
)

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
