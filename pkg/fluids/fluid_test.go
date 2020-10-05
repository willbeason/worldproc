package fluids

import (
	"fmt"
	"github.com/willbeason/worldproc/pkg/geodesic"
	"testing"
)

func TestFlow(t *testing.T) {
	tcs := []struct{
		name string
		start []float64
		want []float64
	} {
		{
			name: "Two nodes",
			start: []float64{1.0, -1.0},
			want: []float64{0.5, 0.5},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGraph(len(tc.start))
			copy(g.Nodes, tc.start)
			g.Link(0, 1, geodesic.Vector{1, 0, 0})

			energy := 0.0
			for _, n := range g.Nodes {
				energy += n * n
			}
			for _, e := range g.Edges {
				energy += e.Velocity.Length2()
			}
			fmt.Println("Energy: ", energy)

			for i := 0; i < 10; i++ {
				before := g.Nodes[0]
				Flow(g, 0.1)
				fmt.Println(i)
				fmt.Println(g.Nodes, (*g.Edges[0]).Velocity)
				energy := 0.0
				for _, n := range g.Nodes {
					energy += n * n
				}
				for _, e := range g.Edges {
					energy += e.Velocity.Length2()
				}
				fmt.Println("Energy: ", energy)
				after := g.Nodes[0]
				if after > before {
					//fmt.Println(after, ">", before)
					//break
				}
			}
			t.Fail()
		})
	}

}
