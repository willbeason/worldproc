package climate

import (
	"fmt"
	"github.com/willbeason/hydrology/pkg/geodesic"
	"math"
	"sort"
	"sync"
)

const (
	// w is the coriolis rotation vector magnitude in inverse minutes.
	// Points in the positive Z direction.
	w = 2 * math.Pi / 1440

	density    = 10.0
	Viscosity = 0.04
	LandDrag = 0.005

	invDensity = 1.0 / density

	twoW = 2 * w

	wSquared = w * w


	ThirdViscosity = Viscosity / 3.0

)

type airDelta struct {
	idx            int
	air            float64
	energy         float64
	n0, n1         int
	theta0, theta1 float64
}

type airNeighbor struct {
	idx        int
	toNeighbor geodesic.Vector
	theta      float64
}

func Flow(climates []Climate, sphere *geodesic.Geodesic, minutes float64) {
	fmt.Print("A")
	// Precalculate the pressure everywhere.
	velocities := make([]geodesic.Vector, len(climates))
	for i, c := range climates {
		velocities[i] = c.AirVelocity
	}

	pressures := make([]float64, len(climates))
	divUs := make([]float64, len(climates))
	for i, c := range climates {
		pressures[i] = c.Pressure()
		divUs[i] = Divergence(i, velocities, sphere)
	}

	fmt.Print("B")
	// Adjust acceleration in every cell.
	for i, c := range climates {
		center := sphere.Centers[i]
		p := PressureGradient(i, pressures, sphere)

		lenP := p.Length()
		pg := geodesic.Vector{}
		if lenP > 1E-10 {
			pg = p.Reject(center).Normalize().Scale(lenP)
		}

		laplacianU := LaplacianVelocity(i, climates, sphere)
		gradDivU := Gradient(i, divUs, sphere)
		a := AirAcceleration(pg, c.AirVelocity, center, laplacianU, gradDivU)
		//fmt.Println("Acceleration", i, a)

		dv := a.Scale(minutes)
		// Subtract out projection of node's face.
		dv = dv.Reject(center)
		climates[i].AirVelocity = climates[i].AirVelocity.Add(dv)
	}

	fmt.Print("C")
	wg := sync.WaitGroup{}
	deltas := make(chan airDelta)
	maxWorkers := 8
	for worker := 0; worker < maxWorkers; worker++ {
		wg.Add(1)
		start := worker * len(climates) / maxWorkers
		end := (worker + 1) * len(climates) / maxWorkers
		go func() {
			for i, c := range climates[start:end] {
				calculateDelta(start+i, c.AirVelocity, c.Air, c.AirEnergy, minutes, sphere, deltas)
			}
			wg.Done()
		}()
	}

	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	// Record air and energy transfers in every cell.
	go func() {
		for delta := range deltas {
			climates[delta.idx].Air -= delta.air
			climates[delta.idx].AirEnergy -= delta.energy
			invSum := 1.0 / (delta.theta0 + delta.theta1)
			climates[delta.n0].Air += delta.air * delta.theta1 * invSum
			climates[delta.n0].AirEnergy += delta.energy * delta.theta1 * invSum
			climates[delta.n1].Air += delta.air * delta.theta0 * invSum
			climates[delta.n1].AirEnergy += delta.energy * delta.theta0 * invSum
		}
		wg2.Done()
	}()
	wg.Wait()
	close(deltas)
	wg2.Wait()
}

func calculateDelta(i int, airVelocity geodesic.Vector, air, airEnergy float64, minutes float64, sphere *geodesic.Geodesic, out chan airDelta) {
	neighbors := sphere.Faces[i].Neighbors
	center := sphere.Centers[i]

	if airVelocity.Dot(airVelocity) < 0.00001 {
		return
	}
	norm := airVelocity.Normalize()

	aNeighbors := make([]airNeighbor, len(neighbors))
	for i, n := range neighbors {
		aNeighbors[i].idx = n
		aNeighbors[i].toNeighbor = sphere.Centers[n].Sub(center).Normalize()
		aNeighbors[i].theta = math.Acos(aNeighbors[i].toNeighbor.Dot(norm))
	}

	sort.Slice(aNeighbors, func(i, j int) bool {
		// Smallest angles first
		return aNeighbors[i].theta < aNeighbors[j].theta
	})

	n0 := aNeighbors[0].idx
	n1 := aNeighbors[1].idx

	theta0 := aNeighbors[0].theta
	theta1 := aNeighbors[1].theta

	// Don't let pressure get below 0.01.
	// Delta air is half of the velocity.
	outAir := math.Min(air-0.01, airVelocity.Length()*minutes) * 0.5
	outEnergy := airEnergy * outAir / air

	out <- airDelta{
		idx:    i,
		air:    outAir,
		energy: outEnergy,
		n0:     n0,
		n1:     n1,
		theta0: theta0,
		theta1: theta1,
	}
}

func LaplacianVelocity(idx int, climates []Climate, sphere *geodesic.Geodesic) geodesic.Vector {
	neighbors := sphere.Faces[idx].Neighbors

	v := climates[idx].AirVelocity
	result := geodesic.Vector{}
	for _, n := range neighbors {
		result = result.Add(climates[n].AirVelocity.Sub(v))
	}
	return result
}

func Divergence(idx int, vectors []geodesic.Vector, sphere *geodesic.Geodesic) float64 {
	neighbors := sphere.Faces[idx].Neighbors
	start := sphere.Centers[idx]

	result := 0.0
	for _, n := range neighbors {
		v := vectors[n]
		toN := sphere.Centers[n].Sub(start).Normalize()
		result += toN.Dot(v)
	}
	return result * 2.0 / float64(len(neighbors))
}

func PressureGradient(idx int, scalars []float64, sphere *geodesic.Geodesic) geodesic.Vector {
	neighbors := sphere.Faces[idx].Neighbors

	p := scalars[idx]
	start := sphere.Centers[idx]

	result := geodesic.Vector{}
	for _, n := range neighbors {
		if scalars[n] > p {
			continue
		}
		toN := sphere.Centers[n].Sub(start).Normalize()
		result = result.Add(toN.Scale(scalars[n] - p))
	}
	return result.Scale(2.0 / float64(len(neighbors)))
}

func Gradient(idx int, scalars []float64, sphere *geodesic.Geodesic) geodesic.Vector {
	neighbors := sphere.Faces[idx].Neighbors

	p := scalars[idx]
	start := sphere.Centers[idx]

	result := geodesic.Vector{}
	for _, n := range neighbors {
		toN := sphere.Centers[n].Sub(start).Normalize()
		result = result.Add(toN.Scale(scalars[n] - p))
	}
	return result.Scale(2.0 / float64(len(neighbors)))
}

// AirAcceleration returns the acceleration of air in inverse minutes squared.
//
// p is the pressure gradient.
// u is the velocity of the fluid.
// x is the position of the fluid.
func AirAcceleration(p geodesic.Vector, u geodesic.Vector, x geodesic.Vector, laplacianU geodesic.Vector, gradDivU geodesic.Vector) geodesic.Vector {
	// Gradient Pressure.
	result := geodesic.Vector{
		X: -invDensity * p.X,
		Y: -invDensity * p.Y,
		Z: -invDensity * p.Z,
	}
	//fmt.Println("Gradient Pressure", result)

	// Coriolis Force.
	result.X += twoW * u.Y
	result.Y -= twoW * u.X

	// Centrifugal Force.
	result.X += wSquared * x.X
	result.Y += wSquared * x.Y

	// Viscosity 1.
	result.X += Viscosity * laplacianU.X
	result.Y += Viscosity * laplacianU.Y
	result.Z += Viscosity * laplacianU.Z

	// Viscosity 2.
	result.X += ThirdViscosity * gradDivU.X
	result.Y += ThirdViscosity * gradDivU.Y
	result.Z += ThirdViscosity * gradDivU.Z

	// Drag against land. Ensures we don't end up with infinite velocity relative to land.
	result.X -= LandDrag * u.X
	result.Y -= LandDrag * u.Y
	result.Z -= LandDrag * u.Z

	return result
}
