package climate

import (
	"fmt"
	"github.com/willbeason/hydrology/pkg/geodesic"
	"math"
	"sort"
)


// k is the drag force constant.
const k = 0.0001

const twoOmegaOverK = 2 * w / k

const invCoriolisDrag = 1 / (1 + 4*w*w/k*k)

// pc is the pressure force constant
const pc = 1.0

const pcInvK = pc / k

type pressureNeighbor struct {
	toNeighbor   geodesic.Vector
	idx          int
	pressureDiff float64
}

type airTransfer struct {
	// air is the amount of air to transfer.
	air float64
	// energy is the amount of energy to transfer.
	energy float64
}

func Wind(climates []Climate, sphere *geodesic.Geodesic) {
	transfers := make([]airTransfer, len(climates))

	maxPressure := 0.0
	for i, c := range climates {
		// First, find the neighbor with the lowest pressure.
		// Wind will be attempting to flow in that direction.
		pressure := c.Pressure()
		maxPressure = math.Max(pressure, maxPressure)
		v := sphere.Centers[i]

		var neighbors []pressureNeighbor
		for _, n := range sphere.Faces[i].Neighbors {
			c := climates[n]
			np := c.Pressure()
			neighbors = append(neighbors, pressureNeighbor{
				toNeighbor:   v.Sub(sphere.Centers[n]).Normalize(),
				idx:          n,
				pressureDiff: pressure - np,
			})
		}

		sort.Slice(neighbors, func(i, j int) bool {
			return neighbors[i].pressureDiff > neighbors[j].pressureDiff
		})

		if neighbors[0].pressureDiff <= 0 {
			// No lower neighbors; we're at a local minimum.
			continue
		}
		// Distribute all air. This will be highly chaotic.
		deltaPressure := neighbors[0].pressureDiff + math.Max(0.0, neighbors[1].pressureDiff)
		deltaPressure *= 0.42
		// At most half of the available air.
		deltaAir := math.Min(deltaPressure*c.Air, 0.5*c.Air)

		pv := neighbors[0].toNeighbor.Scale(neighbors[0].pressureDiff)
		if neighbors[1].pressureDiff > 0 {
			pv = pv.Add(neighbors[1].toNeighbor.Scale(neighbors[1].pressureDiff))
		}
		// pv is the equilibrium vector of air.
		pv = CoriolisDeflect(pv)

		sort.Slice(neighbors, func(i, j int) bool {
			return neighbors[i].toNeighbor.Dot(pv) > neighbors[j].toNeighbor.Dot(pv)
		})

		//pvn := pv.Normalize()
		//theta0 := neighbors[0].toNeighbor.Dot(pvn)
		//theta1 := neighbors[1].toNeighbor.Dot(pvn)
		//invSum := 1.0 / (theta0 + theta1)
		//to0 := deltaAir * theta1 * invSum
		//to1 := deltaAir * theta0 * invSum

		if deltaAir > c.Air {
			panic(i)
		}
		transfers[i].air -= deltaAir
		deltaEnergy := c.AirEnergy * deltaAir / c.Air
		transfers[i].energy -= deltaEnergy

		transfers[neighbors[0].idx].air += deltaAir * 0.5
		transfers[neighbors[0].idx].energy += deltaEnergy * 0.5
		transfers[neighbors[1].idx].air += deltaAir * 0.5
		transfers[neighbors[1].idx].energy += deltaEnergy * 0.5
		//deltaI := -deltaEnergy + deltaEnergy * to0 * invSum + deltaEnergy * to1 * invSum
		//if deltaEnergy > 1E8 || transfers[neighbors[0].idx].energy > 1E8 || transfers[neighbors[1].idx].energy > 1E8 {
		//	fmt.Println()
		//	fmt.Println(theta0, theta1)
		//	fmt.Println("delta air", deltaAir)
		//	fmt.Println("delta energy", deltaEnergy)
		//	fmt.Println("delta energy 0", deltaEnergy * to0 * invSum)
		//	fmt.Println("delta energy 1", deltaEnergy * to1 * invSum)
		//	fmt.Println("transfer 0", transfers[neighbors[0].idx])
		//	fmt.Println("transfer 1", transfers[neighbors[1].idx])
		//	fmt.Println(c)
		//	fmt.Println(to0, to1, invSum)
		//	panic("HERE")
		//}
	}
	fmt.Printf(" pressure=%.04f ", maxPressure)

	maxAir := 0.0
	maxEnergy := 0.0
	sumEnergy := 0.0
	for i, c := range climates {
		t := transfers[i]
		maxAir = math.Max(t.air, maxAir)
		maxEnergy = math.Max(t.energy, maxEnergy)
		sumEnergy += t.energy
		if t.energy > 1E8 {
			panic(i)
		}
		c.AirEnergy += t.energy
		c.Air += t.air
		climates[i] = c
	}
	fmt.Printf(" transfer air=%.04f ", maxAir)
	fmt.Printf(" transfer energy=%.04f ", maxEnergy)
	if math.Abs(sumEnergy) > 1.0 {
		fmt.Println("energy delta", sumEnergy)
		panic("error")
	}
}

func CoriolisDeflect(dp geodesic.Vector) geodesic.Vector {
	return geodesic.Vector{
		X: dp.X + twoOmegaOverK*dp.Y,
		Y: dp.Y - twoOmegaOverK*dp.X,
		Z: dp.Z,
	}
}
