package sun

import (
	"github.com/willbeason/hydrology/pkg/geodesic"
	"math"
)

type Light interface {
	Intensity(vector geodesic.Vector) float64
	AltitudeAzimuth(angle geodesic.Angle) geodesic.Angle
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

func (s *Directional) AltitudeAzimuth(a geodesic.Angle) geodesic.Angle {
	v := a.Vector()
	w := a.Theta - s.SunAngle.Theta

	altitude := math.Asin(s.Sun.Dot(v))
	declination := s.SunAngle.Theta
	//fmt.Println("declination =", declination)
	//fmt.Println("w =", w, "=", a.Theta, "-", s.SunAngle.Theta)
	//fmt.Println("altitude =", altitude)
	azimuth := math.Asin(-math.Cos(declination)*math.Sin(w)/math.Cos(altitude))
	//fmt.Println("azimuth =", azimuth)

	//if rand.Float64() < 0.0001 {
	//	fmt.Println(a.Phi, s.SunAngle.Phi)
	//}

	dPhi := a.Phi - s.SunAngle.Phi
	isPM := dPhi > 0 && math.Abs(dPhi) < math.Pi || (s.SunAngle.Phi > 0 && math.Abs(dPhi) > math.Pi)
	if isPM {
		azimuth = math.Pi - azimuth
		//isPM = a.Phi > s.SunAngle.Phi
		//if a.Phi > 0 || (a.Phi + s.SunAngle.Phi ){
		//}
		//isPM = true
		//isPM = a.Phi < math.Pi || a.Phi < s.SunAngle.Phi
	}

	return geodesic.Angle{
		// Theta = Pi/2 corresponds with sun directly overhead.
		Theta: altitude,
		// Phi = 0 corresponds with azimuth to North.
		Phi: azimuth,
	}
}

// Set sets the planet's date, in days since spring equinox year 0.
func (s *Directional) Set(date float64) {
	// Assume 360 days, 23.5 degree axial tilt and a circular orbit.
	// Start at the spring equinox for the northern hemisphere.
	eclipticLatitude := -23.5 * math.Sin(date * math.Pi / 180) * math.Pi / 180
	// Start at noon on the prime meridian.
	eclipticLongitude := (0.5-math.Mod(date+0.5, 1.0)) * 2 * math.Pi

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
