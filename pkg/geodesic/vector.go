package geodesic

import "math"

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (v Vector) Add(v2 Vector) Vector {
	return Vector{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
		Z: v.Z + v2.Z,
	}
}

func (v Vector) Scale(c float64) Vector {
	return Vector{
		X: v.X * c,
		Y: v.Y * c,
		Z: v.Z * c,
	}
}

func (v Vector) Normalize() Vector {
	lInv := 1.0 / math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	return Vector{
		X: v.X * lInv,
		Y: v.Y* lInv,
		Z: v.Z * lInv,
	}
}

type Angle struct {
	// Theta is the angle with the equator.
	Theta float64
	// Phi is the angle with the Prime Meridian.
	Phi float64
}

func (a Angle) Vector() Vector {
	return Vector{
		X: math.Cos(a.Theta)*math.Cos(a.Phi),
		Y: math.Cos(a.Theta)*math.Sin(a.Phi),
		Z: math.Sin(a.Theta),
	}
}
