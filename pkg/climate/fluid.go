package climate

import (
	"fmt"
	"github.com/willbeason/hydrology/pkg/geodesic"
	"math"
	"sort"
)

const (
	// w is the coriolis rotation vector magnitude in inverse minutes.
	// Points in the positive Z direction.
	w = 2 * math.Pi / 1440

	density    = 50.0
	invDensity = 1.0 / density

	twoW = 2 * w

	wSquared = w * w

	Viscosity = 0.10

	ThirdViscosity = Viscosity / 3.0

	LandDrag = 0.001
)

type airDelta struct {
	air    float64
	energy float64
}

type airNeighbor struct {
	idx        int
	toNeighbor geodesic.Vector
	theta      float64
}

func Flow(climates []Climate, sphere *geodesic.Geodesic, minutes float64) {
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

	// Adjust acceleration in every cell.
	for i, c := range climates {
		center := sphere.Centers[i]
		p := Gradient(i, pressures, sphere)

		lenP := math.Sqrt(p.Dot(p))
		pg := geodesic.Vector{}
		if lenP > 0.0001 {
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
		//if i == 0 {
		//	fmt.Println("p", p)
		//	fmt.Println("a", a)
		//	fmt.Println("delta v", dv)
		//	fmt.Println("v", climates[i].AirVelocity)
		//}
	}

	// Record air and energy transfers in every cell.
	deltaAir := make([]airDelta, len(climates))
	for i, c := range climates {
		neighbors := sphere.Faces[i].Neighbors
		center := sphere.Centers[i]

		if c.AirVelocity.Dot(c.AirVelocity) < 0.00001 {
			continue
		}
		norm := c.AirVelocity.Normalize()

		aNeighbors := make([]airNeighbor, len(neighbors))
		for i, n := range neighbors {
			aNeighbors[i].idx = n
			aNeighbors[i].toNeighbor = sphere.Centers[n].Sub(center).Normalize()
			aNeighbors[i].theta = math.Acos(aNeighbors[i].toNeighbor.Dot(norm))
			//if math.IsNaN(aNeighbors[i].theta) {
				//fmt.Println(aNeighbors[i])
				//fmt.Println(norm)
				//panic("H")
			//}
		}
		//fmt.Println(aNeighbors)

		sort.Slice(aNeighbors, func(i, j int) bool {
			// Smallest angles first
			return aNeighbors[i].theta < aNeighbors[j].theta
		})

		n0 := aNeighbors[0].idx
		n1 := aNeighbors[1].idx

		theta0 := aNeighbors[0].theta
		theta1 := aNeighbors[1].theta
		if theta0 > math.Pi/2 {
			continue
		}

		invSum := 1.0 / (theta0 + theta1)

		outAir := math.Min(c.Air-0.01, c.AirVelocity.Length() * minutes)
		outEnergy := c.AirEnergy * outAir / c.Air
		deltaAir[i].air -= outAir
		deltaAir[i].energy -= outEnergy

		deltaAir[n0].air += theta1 * outAir * invSum
		deltaAir[n0].energy += theta1 * outEnergy * invSum
		deltaAir[n1].air += theta0 * outAir * invSum
		// Don't let pressure get below 0.01.
		deltaAir[n1].energy += theta0 * outEnergy * invSum
		//if i == 0 {
		//	fmt.Println(i)
		//	fmt.Println(norm)
		//	fmt.Println(n0, n1)
		//	fmt.Println(theta0, theta1, theta0+theta1, invSum)
		//	fmt.Println("out air", outAir)
		//	fmt.Println("delta 0", theta1*outAir*invSum)
		//	fmt.Println("delta 1", theta0*outAir*invSum)
		//}
	}

	sumDeltaAir := 0.0
	maxDeltaAir := 0.0
	maxVelocity := 0.0
	for i, d := range deltaAir {
		climates[i].Air += d.air
		climates[i].AirEnergy += d.energy
		sumDeltaAir += math.Abs(d.air)
		maxDeltaAir = math.Max(maxDeltaAir, d.air)
		maxVelocity = math.Max(maxVelocity, climates[i].AirVelocity.Length())
		//climates[i].AirVelocity = geodesic.Vector{}
	}
	fmt.Printf(" deltaMax(%.00f, %.04f, %.04f) ", sumDeltaAir, maxDeltaAir, maxVelocity)
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
