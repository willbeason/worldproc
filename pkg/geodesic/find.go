package geodesic

import (
	"math"
)

// NaiveFind is the naive algorithm for determining the face containing a point.
// Searches every face on the sphere, so it's incredibly inefficient for large
// numbers of faces.
func NaiveFind(g *Geodesic, v Vector) int {
	start := 0
	minDistSq := math.MaxFloat64
	for i, center := range g.Centers {
		iDistSq := DistSq(center, v)
		if iDistSq < minDistSq {
			minDistSq = iDistSq
			start = i
		}
	}

	return start
}

// Find returns the *Node closest to v in the last Geodesic in gs.
// gs is a precomputed sequence of geodesics, each homomorphic to the previous one.
// Further, if face F is in gs[i] and gs[k], then gs[i].Center[F] = gs[j].Center[F].
func Find(gs []*Geodesic, v Vector) int {
	if len(gs) == 1 {
		return NaiveFind(gs[0], v)
	}

	gs0 := gs[0]

	start := 0
	minDistSq := math.MaxFloat64
	if v.Z > 0 {
		// We are in the northern hemisphere, so we can disregard all southern
		// faces.
		// Start from the north pole at index 0.
		minDistSq = DistSq(gs0.Centers[start], v)
		for i := 1; i <= 5; i++ {
			iDistSq := DistSq(gs0.Centers[i], v)
			if iDistSq < minDistSq {
				minDistSq = iDistSq
				start = i
			}
		}
	} else {
		// We are in the southern hemisphere, so we can disregard all northern
		// faces.
		// Start from the south pole at index 11.
		start = 11
		minDistSq = DistSq(gs0.Centers[start], v)
		for i := 6; i <= 10; i++ {
			iDistSq := DistSq(gs0.Centers[i], v)
			if iDistSq < minDistSq {
				minDistSq = iDistSq
				start = i
			}
		}
	}

	if len(gs) == 1 {
		return start
	}

	return find(gs[1:], v, minDistSq, start)
}

func DistSq(v1, v2 Vector) float64 {
	return (v1.X-v2.X)*(v1.X-v2.X) + (v1.Y-v2.Y)*(v1.Y-v2.Y) + (v1.Z-v2.Z)*(v1.Z-v2.Z)
}

// find returns the *Node closest to v, starting from the node at index start.
// start must be the face index closest to v in the previous geodesic.
func find(gs []*Geodesic, v Vector, minDistSq float64, start int) int {
	gs0 := gs[0]
	nextStart := start
	neighbors := gs0.Faces[start].Neighbors
	for _, n := range neighbors {
		iDistSq := DistSq(gs0.Centers[n], v)
		if iDistSq < minDistSq {
			minDistSq = iDistSq
			nextStart = n
		}
	}

	if len(gs) == 1 {
		// There's a bug related to correctly classifying the neighbor of a
		// pentagon face, so this corrects for that.
		if nextStart == start {
			return nextStart
		}
		return find(gs, v, minDistSq, nextStart)
	}

	return find(gs[1:], v, minDistSq, nextStart)
}
