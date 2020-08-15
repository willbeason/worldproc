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

const DesertSpecificHeat = 200000
const CoastSpecificHeat =  500000
const OceanSpecificHeat = 1000000

type Climate struct {
	// SpecificHeat is the energy, in Joules/m^2 required to heat this area by 1 K.
	SpecificHeat float64

	// Temperature is how hot, in Kelvin, this climate is.
	Temperature float64
}

func (t *Climate) Simulate(flux float64, altitude float64, seconds float64) {
	incoming := flux * seconds

	peakWavelength := WD / t.Temperature
	offset := -(peakWavelength - 1E-5) / (1E-5)
	//fmt.Println(peakWavelength, offset)
	opacity := 0.26 + offset + altitude / 2.0
	//fmt.Println(peakWavelength)
	opacity = math.Max(0.0, math.Min(0.9, opacity))
	outgoing := seconds * math.Pow(t.Temperature, 4) * opacity * SB

	t.Temperature += (incoming - outgoing) / t.SpecificHeat
}

func LowHigh(specificHeat, latitude, startNoon float64) (float64, float64) {
	temp := startNoon
	cosLatitude := math.Cos(latitude)

	c := &Climate{
		SpecificHeat: specificHeat,
		Temperature:  startNoon,
	}

	lowest, highest := temp, temp
	for i := 0; i < 144; i++ {
		sunAngle := float64(i) * math.Pi / 72
		flux := Flux * math.Cos(sunAngle) * cosLatitude
		flux = math.Max(0.0, flux)
		c.Simulate(flux, 0.0, 600)
		temp = c.Temperature

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
	return LowHigh(specificHeat, latitude, temp)
}
