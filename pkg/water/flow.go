package water

import (
	"github.com/willbeason/worldproc/pkg/geodesic"
)

func Rain(amt float64, waters, heights, flow []float64, sphere *geodesic.Geodesic) {
	for i := range waters {
		rainFlow(amt, i, waters, heights, flow, sphere)
	}
}

func rainFlow(amt float64, idx int, waters, heights, flow []float64, sphere *geodesic.Geodesic) {
	flow[idx] += amt

	flowTo := idx
	flowToWh := waters[idx] + heights[idx]
	for _, n := range sphere.Faces[idx].Neighbors {
		nwh := waters[n] + heights[n]
		if nwh < flowToWh {
			flowTo = n
			flowToWh = nwh
		}
	}

	if flowTo != idx {
		rainFlow(amt, flowTo, waters, heights, flow, sphere)
		return
	}

	// All goes to this index.
	waters[idx] += amt
	return
}
