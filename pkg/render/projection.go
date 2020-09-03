package render

import (
	"github.com/willbeason/worldproc/pkg/geodesic"
	"math"
)

type Projection struct {
	Screen
	Pixels []geodesic.Angle
}

func Project(screen Screen, projector Projector) Projection {
	invWidth := 1.0 / float64(screen.Width + 1)
	invHeight := 1.0 / float64(screen.Height + 1)

	result := Projection{
		Screen: screen,
		Pixels: make([]geodesic.Angle, screen.Width*screen.Height),
	}

	for px := 0; px < screen.Width; px++ {
		x := (float64(px) + 0.5 - (float64(screen.Width) / 2.0))  * invWidth

		for py := 0; py < screen.Height; py++ {
			y :=(float64(py) - (float64(screen.Height) / 2.0)) * invHeight
			result.Pixels[py*screen.Width+px] = projector.Project(x, y)
		}
	}

	return result
}

type Projector interface {
	// Project transforms a pair of x,y coordinates representing a location in
	// a rectangle into an angle on a geodesic sphere.
	//
	// -1 < x,y < 1
	// x corresponds to left/right and y to up/down:
	// -1 is left/bottom
	// 0 is center
	// +1 is right/top
	Project(x, y float64) geodesic.Angle
}

type Equirectangular struct {}

func (Equirectangular) Project(x, y float64) geodesic.Angle {
	return geodesic.Angle{
		Theta: y * math.Pi,
		Phi: x * math.Pi * 2,
	}
}
