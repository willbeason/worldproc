package fluids

import (
	"fmt"
	"github.com/willbeason/worldproc/pkg/geodesic"
)

type Edge struct {
	FromIdx, ToIdx int

	// From and To are the nodes this Edge connects.
	From, To *float64

	// Direction is the unit vector direction perpendicular to the edge from From
	// and to To.
	Direction geodesic.Vector

	// Velocity is the velocity of fluid at the center of this Edge.
	Velocity geodesic.Vector
}

type Graph struct {
	Nodes []float64
	Edges []*Edge
}

func NewGraph(n int) *Graph {
	return &Graph{
		Nodes: make([]float64, n),
		Edges: nil,
	}
}

func (g *Graph) Link(a, b int, direction geodesic.Vector) {
	g.Edges = append(g.Edges, &Edge{
		FromIdx:   a,
		ToIdx:     b,
		From:      &g.Nodes[a],
		To:        &g.Nodes[b],
		Direction: direction,
		Velocity:  geodesic.Vector{},
	})
}

func Flow(g *Graph, t float64) {
	deltas := make([]float64, len(g.Nodes))
	for _, e := range g.Edges {
		//negativeGradP1 := e.Direction.Scale(*e.From - *e.To)
		//previousFrom := e.Velocity.Dot(e.Direction) * t
		//negativeGradP2 := e.Direction.Scale(*e.From - *e.To - previousFrom)
		//
		//deltaV1 := negativeGradP1
		//deltav2 := negativeGradP2.Sub(negativeGradP1).Scale(0.25 * t)
		//e.Velocity = e.Velocity.Add(deltaV1.Add(deltav2).Scale(t))
		//
		//nextFrom := e.Velocity.Dot(e.Direction) * t
		//deltas[e.FromIdx] -= nextFrom
		//deltas[e.ToIdx] += nextFrom

		previousFrom := e.Velocity.Dot(e.Direction) * t
		negGrad := e.Direction.Scale(*e.From - *e.To)
		deltaV := negGrad
		e.Velocity = e.Velocity.Add(deltaV.Scale(t))

		nextFrom := e.Velocity.Dot(e.Direction) * t
		from := (previousFrom + nextFrom) / 2.0
		*e.From -= from
		*e.To += from
		fmt.Println("New V", e.Velocity)
	}

	for i, d := range deltas {
		g.Nodes[i] += d
	}
}
