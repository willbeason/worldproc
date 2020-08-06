package water

import (
	"github.com/willbeason/hydrology/pkg/geodesic"
	"math"
	"testing"
)

func TestRain(t *testing.T) {

	tcs := []struct{
		name string
		height []float64
		maxDiff float64
	} {
		{
			name: "length 2",
			height: make([]float64, 2),
			maxDiff: 0.0001,
		},
		{
			name: "length 10",
			height: make([]float64, 10),
			maxDiff: 0.00020,
		},
		{
			name: "length 100",
			height: make([]float64, 100),
			maxDiff: 0.03500,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			s := &geodesic.Geodesic{
				Faces: make([]geodesic.Node, len(tc.height)),
			}
			for i := range tc.height[1:] {
				s.Faces[i].Neighbors = append(s.Faces[i].Neighbors, i+1)
				s.Faces[i+1].Neighbors = append(s.Faces[i+1].Neighbors, i)
			}

			waters := make([]float64, len(tc.height))
			flow := make([]float64, len(tc.height))

			n := 10000
			for i := 0; i < n; i++ {
				//fmt.Println("i: ", i)
				rainFlow(0.001, 0, waters, tc.height, flow, s)
				//fmt.Println(waters)
			}

			avg := float64(n)*0.001 / float64(len(tc.height))
			offset := 0.0
			for _, w := range waters {
				offset += math.Abs(w - avg)
			}
			if offset > tc.maxDiff {
				//t.Error(waters)
				t.Errorf("%.04f > %.04f", offset, tc.maxDiff)
			}
		})
	}
}
