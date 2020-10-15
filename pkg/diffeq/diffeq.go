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
	k1 := x0 - y0
	// Euler's method on the velocity.
	v1 := v0 + h*k1
	// Trapezoid to estimate x and y.

	// Second estimate of f.
	k2 := k1 - h*(v0+v1)
	// Trapezoid on the two estimates of f.
	ho2 := h/2
	v2 := v0 + ho2*(k1+k2)
	x2 := x0 - ho2*(v0+v2)
	y2 := y0 + ho2*(v0+v2)

	return x2, y2, v2
}

func RK4(h, x0, y0, v0 float64) (float64, float64, float64) {
	// First estimate at h/2.
	k1 := x0 - y0
	ho2 := h/2
	v1 := v0 + ho2*k1

	// Second estimate at h/2.
	k2 := k1 - ho2*(v0+v1)
	v2 := v0 + ho2*k2
	wv2 := (2*v0+v1+v2)/4

	// Estimate at h.
	k3 := k1 - h*wv2
	v3 := v0 + h*k3
	ho6 := h/6
	wv3 := (v0 + 2*v1 + 2*v2 + v3)*ho6
	x3 := x0 - wv3
	y3 := y0 + wv3

	// Complete estimate.
	k4 := x3 - y3
	v := v0 + ho6*(k1+2*k2+2*k3+k4)
	wv4 := (v0 + 2*v1 + 2*v2 + (v3+v)/2)*ho6
	x4 := x0 - wv4
	y4 := y0 + wv4

	return x4, y4, v
}
