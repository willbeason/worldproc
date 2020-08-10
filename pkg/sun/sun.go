package sun

import (
	"github.com/willbeason/hydrology/pkg/geodesic"
	"math"
)

type Light interface {
	Intensity(vector geodesic.Vector) float64
	DeclinationAzimuth(angle geodesic.Angle) geodesic.Angle
}

type Constant struct {}

func (s Constant) Intensity(_ geodesic.Vector) float64 {
	return 1.0
}

type Directional struct {
	// Sun is the directional vector from the planet's core to the Sun.
	Sun geodesic.Vector
	SunAngle geodesic.Angle
}

func (s *Directional) DeclinationAzimuth(a geodesic.Angle) geodesic.Angle {
	v := a.Vector()
	L := s.SunAngle.Phi - a.Phi
	azimuth := math.Atan(math.Sin(L) / (math.Cos(a.Theta)*math.Tan(math.Pi/2 - s.SunAngle.Theta) - math.Sin(a.Theta)*math.Cos(L)))
	if a.Phi < s.SunAngle.Phi {
		//azimuth = -azimuth
	}

	return geodesic.Angle{
		// Theta = 0 corresponds with sun directly overhead.
		Theta: math.Pi/2 - math.Acos(s.Sun.Dot(v)),
		// Phi = 0 corresponds with azimuth to North.
		Phi:   azimuth,
	}
}

// Set sets the planet's date, in days since spring equinox year 0.
func (s *Directional) Set(date float64) {
	// Assume 360 days, 23.5 degree axial tilt and a circular orbit.
	// Start at the spring equinox for the northern hemisphere.
	eclipticLatitude := -23.5 * math.Sin(date * math.Pi / 180) * math.Pi / 180
	// Start at noon on the prime meridian.
	eclipticLongitude := -math.Mod(date, 1.0) * 2 * math.Pi

	s.SunAngle = geodesic.Angle{
		Phi: eclipticLongitude,
		Theta: eclipticLatitude,
	}
	s.Sun = s.SunAngle.Vector()
}

func (s *Directional) Intensity(v geodesic.Vector) float64 {
	dot := s.Sun.Dot(v)
	dot *= 2
	return math.Max(math.Min(1.0, dot), 0.1)
}
