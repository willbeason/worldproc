package geodesic

import (
	"math"
)

type Node struct {
	Neighbors []int `json:"Neighbors"`
}

type Edge struct {
	L int `json:"L"`
	R int `json:"R"`
}

// Geodesic represents a geodesic sphere.
type Geodesic struct {
	// Centers is a list of vectors representing the center of every face of the
	// geodesic sphere.
	Centers []Vector `json:"centers"`
	// Faces is a list of the neighboring faces of each
	Faces []Node `json:"nodes"`
	Edges map[Edge]int `json:"-"`
}

const sin_atan0_5 = 0.447213595
const cos_atan0_5 = 0.894427191

// Dodecahedron creates 12-sided polyhedron.
func Dodecahedron() *Geodesic {
	// Generate faces with centers corresponding to the points of an icosahedron.
	g := &Geodesic{
		Centers: []Vector{
			// North Pole
			{X: 0.0, Y: 0.0, Z: 1.0},
			// Northern
			{X: math.Sin(0*math.Pi/5)*cos_atan0_5, Y: math.Cos(0*math.Pi/5)*cos_atan0_5, Z: sin_atan0_5},
			{X: math.Sin(2*math.Pi/5)*cos_atan0_5, Y: math.Cos(2*math.Pi/5)*cos_atan0_5, Z: sin_atan0_5},
			{X: math.Sin(4*math.Pi/5)*cos_atan0_5, Y: math.Cos(4*math.Pi/5)*cos_atan0_5, Z: sin_atan0_5},
			{X: math.Sin(6*math.Pi/5)*cos_atan0_5, Y: math.Cos(6*math.Pi/5)*cos_atan0_5, Z: sin_atan0_5},
			{X: math.Sin(8*math.Pi/5)*cos_atan0_5, Y: math.Cos(8*math.Pi/5)*cos_atan0_5, Z: sin_atan0_5},
			// Southern
			{X: math.Sin(1*math.Pi/5)*cos_atan0_5, Y: math.Cos(1*math.Pi/5)*cos_atan0_5, Z: -sin_atan0_5},
			{X: math.Sin(3*math.Pi/5)*cos_atan0_5, Y: math.Cos(3*math.Pi/5)*cos_atan0_5, Z: -sin_atan0_5},
			{X: math.Sin(5*math.Pi/5)*cos_atan0_5, Y: math.Cos(5*math.Pi/5)*cos_atan0_5, Z: -sin_atan0_5},
			{X: math.Sin(7*math.Pi/5)*cos_atan0_5, Y: math.Cos(7*math.Pi/5)*cos_atan0_5, Z: -sin_atan0_5},
			{X: math.Sin(9*math.Pi/5)*cos_atan0_5, Y: math.Cos(9*math.Pi/5)*cos_atan0_5, Z: -sin_atan0_5},
			// South Pole
			{X: 0.0, Y: 0.0, Z: -1.0},
		},
		Faces: make([]Node, 12),
		Edges: map[Edge]int{},
	}

	// North Pole
	g.Link(0, 1)
	g.Link(0, 2)
	g.Link(0, 3)
	g.Link(0, 4)
	g.Link(0, 5)

	// North of Equator
	g.Link(1, 2)
	g.Link(2, 3)
	g.Link(3, 4)
	g.Link(4, 5)
	g.Link(5, 1)

	// Equator
	g.Link(1, 10)
	g.Link(1, 6)

	g.Link(2, 6)
	g.Link(2, 7)

	g.Link(3, 7)
	g.Link(3, 8)

	g.Link(4, 8)
	g.Link(4, 9)

	g.Link(5, 9)
	g.Link(5, 10)

	// South of Equator
	g.Link(6, 7)
	g.Link(7, 8)
	g.Link(8, 9)
	g.Link(9, 10)
	g.Link(10, 6)

	// South Pole
	g.Link(11, 6)
	g.Link(11, 7)
	g.Link(11, 8)
	g.Link(11, 9)
	g.Link(11, 10)

	return g
}

func (g *Geodesic) Link(i, j int) {
	if _, found := g.Edges[Edge{L: i, R: j}]; found {
		return
	}

	// Link the two nodes both ways.
	g.Faces[i].Neighbors = append(g.Faces[i].Neighbors, j)
	g.Faces[j].Neighbors = append(g.Faces[j].Neighbors, i)

	id := len(g.Edges) / 2
	g.Edges[Edge{L: i, R: j}] = id
	g.Edges[Edge{L: j, R: i}] = id
}

func bisect(n1, n2 Vector) Vector {
	x := n1.X + n2.X
	y := n1.Y + n2.Y
	z := n1.Z + n2.Z

	invLength := 1.0 / math.Sqrt(x*x + y*y + z*z)
	return Vector{
		X: x * invLength,
		Y: y * invLength,
		Z: z * invLength,
	}
}

// Chamfer replaces all edges in the Geodesic with a hexagon.
func Chamfer(g *Geodesic) *Geodesic {
	nFaces := len(g.Faces)
	nEdges := len(g.Edges)/2
	result := &Geodesic{
		Centers: make([]Vector, nFaces+nEdges),
		Faces: make([]Node, nFaces+nEdges),
		Edges: map[Edge]int{},
	}
	copy(result.Centers, g.Centers)

	for faceIdx, face := range g.Faces {
		// For each Node.
		for n1id, n1 := range face.Neighbors {
			// For each neighbor, create a new face.

			// Link it to the edge separating it from its neighbor.
			idIJ := nFaces + g.Edges[Edge{L: faceIdx, R: n1}]
			result.Link(faceIdx, idIJ)

			// The new face's center bisects the two other faces' centers.
			result.Centers[idIJ] = bisect(g.Centers[faceIdx], g.Centers[n1])

			// Link the new neighbor to its neighbors.
			for n2id, n2 := range face.Neighbors {
				if n1id == n2id {
					// We don't allow faces to be self-adjacent.
					continue
				}
				if _, found := g.Edges[Edge{L: n1, R: n2}]; found {
					// The three faces [face/n1/n2] were all adjacent, so the
					// faces (formerly edges) are adjacent to each other.
					idJK := nFaces + g.Edges[Edge{L: n1, R: n2}]
					result.Link(idIJ, idJK)
				}
			}
		}
	}

	return result
}
