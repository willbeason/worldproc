package planet

import (
	"fmt"
	"github.com/willbeason/worldproc/pkg/geodesic"
	"github.com/willbeason/worldproc/pkg/water"
	"sort"
)

const WaterQuanta = 0.01

// AddWater adds water to the Planet.
//
// coverage is the estimate of the planet's surface area to be covered with water.
func AddWater(p *Planet, coverage float64, sphere *geodesic.Geodesic) {
	p.Waters = make([]float64, len(p.Heights))
	p.Flows = make([]float64, len(p.Heights))

	sortedHeights := make([]float64, len(p.Heights))
	copy(sortedHeights, p.Heights)
	sort.Float64s(sortedHeights)

	seaWater := 0.0
	idx := int(float64(len(sortedHeights)) * coverage)
	seaLevel := sortedHeights[idx]
	for _, h := range sortedHeights {
		if h >= seaLevel {
			break
		}
		seaWater += seaLevel - h
	}
	avgWater := seaWater / float64(idx)

	iters := int(avgWater / WaterQuanta)
	for iter := 0; iter < iters; iter++ {
		fmt.Println(iter, "...", "Raining")
		water.Rain(WaterQuanta, p.Waters, p.Heights, p.Flows, sphere)
	}
	fmt.Println("... Equalizing")
	water.Equalize(p.Waters, p.Heights, sphere)
}
