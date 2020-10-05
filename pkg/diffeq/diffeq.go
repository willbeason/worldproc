package diffeq

import (
	"math"
)

type Approximation func(h, x, y, v float64) (float64, float64, float64)

func Approximate(steps int, f Approximation) (float64, float64) {
	h := math.Pi / (2 * math.Sqrt(2) * float64(steps))

	x, y, v := 1.0, -1.0, 0.0

	for i := 0; i < steps; i++ {
		x, y, v = f(h, x, y, v)
	}

	return x, v
}

func Euler(h, x, y, v float64) (float64, float64, float64) {
	v2 := v + h*(x-y)
	x -= h * (v + v2) / 2
	y += h * (v + v2) / 2
	return x, y, v2
}

func Trapezoid(h, x0, y0, v0 float64) (float64, float64, float64) {
	// Initial estimate of f.
	f0 := x0 - y0
	// Euler's method on the velocity.
	v1 := v0 + h*f0
	// Trapezoid to estimate x and y.
	x1 := x0 - h*(v0+v1)/2
	y1 := y0 + h*(v0+v1)/2

	// Second estimate of f.
	f1 := x1 - y1
	// Trapezoid on the two estimates of f.
	v2 := v0 + h*(f0+f1)/2
	x2 := x0 - h*(v0+v2)/2
	y2 := y0 + h*(v0+v2)/2

	return x2, y2, v2
}

func RK4(h, x0, y0, v0 float64) (float64, float64, float64) {
	k1 := x0 - y0

	// First estimate of h/2.
	v1 := v0 + (h/2)*k1
	x1 := x0 - (h/4)*(v0+v1)
	y1 := y0 + (h/4)*(v0+v1)
	k2 := x1 - y1

	// Second estimate of h/2.
	v2 := v0 + (h/2)*k2
	x2 := x0 - (h/4)*(v0+v2)
	y2 := y0 + (h/4)*(v0+v2)
	k3 := x2 - y2

	// Estimate h.
	v3 := v0 + h*k3
	wv := v0+2*v1+2*v2+v3
	x3 := x0 - (h/6)*wv
	y3 := y0 + (h/6)*wv
	k4 := x3 - y3

	// Complete estimate.
	v := v0 + (h/6)*(k1+2*k2+2*k3+k4)
	x := x0 - (h/6)*wv
	y := y0 + (h/6)*wv
	return x, y, v
}
