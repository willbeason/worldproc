package climate

import (
	"math"
)

// Flux is the solar flux at the equator at noon.
const Flux = 400

// SB is the Stefan-Boltzmann constant.
const SB = 5.670374419184429453970996731889231E-8

// WD is Wien's displacement constant.
const WD = 2.8977719E-3

const DefaultTemperature = 273.0

type Climate struct {
	SpecificHeat float64
	//Opacity float64

	// Latitude is the angle with the Equator, in radians.
	Latitude float64
}

func (t Climate) LowHigh(startNoon float64) (float64, float64) {
	temp := startNoon
	energy := t.SpecificHeat * temp
	cosLatitude := math.Cos(t.Latitude)

	lowest, highest := temp, temp
	for i := 0; i < 24; i++ {
		peakWavelength := WD / temp
		opacity := (11.75E-6 - peakWavelength)/(8E-6)
		opacity = math.Max(0.1, math.Min(opacity, 1.0))
		//fmt.Println(temp, opacity)

		sunAngle := float64(i) * math.Pi / 12
		incomingFlux := Flux * math.Cos(sunAngle)*cosLatitude * 3600
		if incomingFlux < 0 {
			incomingFlux = 0
		}
		outgoingFlux := 3600 * math.Pow(temp, 4) * opacity * SB

		flux := incomingFlux - outgoingFlux
		energy += flux
		temp = energy / t.SpecificHeat

		if temp < lowest {
			lowest = temp
		}
		if temp > highest {
			highest = temp
		}
	}

	diff := startNoon - temp
	if math.Abs(diff) < 0.001 {
		return lowest, highest
	}
	return t.LowHigh(temp)
}
