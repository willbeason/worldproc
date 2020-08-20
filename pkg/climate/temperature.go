package climate

import (
	"math"
)

// Flux is the solar flux at the equator at noon.
const Flux = 400

// SB is the Stefan-Boltzmann constant.
const SB = 5.670374419184429453970996731889231E-8

// WD is Wien's displacement constant.
//const WD = 2.8977719E-3

const ZeroCelsius = 273.15

const DefaultAir = 1.0

const DesertSpecificHeat = 100000
const CoastSpecificHeat = 400000
const OceanSpecificHeat = 1400000

const AirSpecificHeat = 100000

type Climate struct {
	// LandSpecificHeat is the energy, in J/m^2 required to heat only the land
	// by 1 K.
	LandSpecificHeat float64

	// LandEnergy is the energy held by land.
	LandEnergy float64

	// Air is the proportion of air.
	// 1.0 is the mean across the planet.
	Air float64

	// AirEnergy is the energy held by the air.
	AirEnergy float64
}

func (t *Climate) LandTemperature() float64 {
	return t.LandEnergy / t.LandSpecificHeat
}

func (t *Climate) AirTemperature() float64 {
	return t.AirEnergy / (t.Air * AirSpecificHeat)
}

func (t *Climate) SetTemperature(kelvin float64) {
	t.LandEnergy = t.LandSpecificHeat * kelvin
	t.AirEnergy = t.Air * AirSpecificHeat * kelvin
}

func (t *Climate) Pressure() float64 {
	return t.Air * t.AirTemperature() / ZeroCelsius
}

func (t *Climate) Simulate(flux float64, latitude float64, altitude float64, seconds float64) {
	incoming := flux * seconds

	// Land absorbs sunlight and cools down, but not air.
	// opacity must be _at least_ 0.5 at the poles
	dLatitude := - math.Cos(latitude) * 0.23
	opacity := 0.5 + (altitude / 3.0) + dLatitude
	opacity = math.Max(0.0, math.Min(1.0, opacity))
	outgoing := seconds * math.Pow(t.LandTemperature(), 4) * opacity * SB

	t.LandEnergy += incoming - outgoing

	// Move towards equilibrium between land/air.
	invSpecificHeat := 1.0 / (t.Air * AirSpecificHeat + t.LandSpecificHeat)
	totalEnergy := t.AirEnergy + t.LandEnergy

	// deltaAirEnergy is the delta to AirEnergy that brings the system to equilibrium.
	deltaAirEnergy := totalEnergy * (t.Air * AirSpecificHeat) * invSpecificHeat - t.AirEnergy

	t.AirEnergy += 0.2 * deltaAirEnergy
	t.LandEnergy -= 0.2 * deltaAirEnergy
}

func yearMax(landSpecificHeat, startTemp float64, latitude float64, maxAngle float64) (float64, float64) {
	max := startTemp

	c := &Climate{
		LandSpecificHeat: landSpecificHeat,
		Air: DefaultAir,
	}
	c.SetTemperature(startTemp)
	for day := 0; day < 360; day++ {
		for hour := 0; hour < 24; hour ++ {
			declination := maxAngle * math.Sin((float64(day) + float64(hour) / 24) * math.Pi / 180)

			flux := Flux * math.Sin(declination)
			flux = math.Max(0.0, flux)

			c.Simulate(flux, latitude, 0.0, 3600)
			temp := c.LandTemperature()
			max = math.Max(temp, max)
		}
	}
	return max, c.LandTemperature()
}

func PoleEquilibrium(specificHeat, maxAngle float64) float64 {
	i := 0

	low := 0.0
	lowMax, _ := yearMax(specificHeat, low, math.Pi / 2, maxAngle)

	high := 2 * ZeroCelsius
	highMax, _ := yearMax(specificHeat, high, math.Pi / 2, maxAngle)

	for math.Abs(highMax - lowMax) > 0.001 {
		mid := (low + high) / 2.0
		max, end := yearMax(specificHeat, mid, math.Pi / 2, maxAngle)

		if end < mid {
			high = mid
			highMax = max
		} else {
			low = mid
			lowMax = max
		}
		i++
	}
	return highMax
}

func LowHigh(specificHeat, latitude, startNoon float64) (float64, float64) {
	temp := startNoon
	cosLatitude := math.Cos(latitude)

	c := &Climate{
		LandSpecificHeat: specificHeat,
		Air: DefaultAir,
	}
	c.SetTemperature(startNoon)

	lowest, highest := temp, temp
	for i := 0; i < 144; i++ {
		sunAngle := float64(i) * math.Pi / 72
		flux := Flux * math.Cos(sunAngle) * cosLatitude
		flux = math.Max(0.0, flux)
		c.Simulate(flux, latitude, 0.0, 600)
		temp = c.AirTemperature()

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
