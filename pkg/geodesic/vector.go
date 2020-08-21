package geodesic

import "math"

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (v Vector) Dot(v2 Vector) float64 {
	return v.X * v2.X + v.Y * v2.Y + v.Z * v2.Z
}

func (v Vector) Cross(v2 Vector) Vector {
	return Vector{
		X: v.Y*v2.Z - v.Z*v2.Y,
		Y: v.Z*v2.X - v.X*v2.Z,
		Z: v.X*v2.Y - v.Y*v2.X,
	}
}

func (v Vector) Add(v2 Vector) Vector {
	return Vector{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
		Z: v.Z + v2.Z,
	}
}

func (v Vector) Sub(v2 Vector) Vector {
	return v.Add(v2.Scale(-1))
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

func (v Vector) Reject(b Vector) Vector {
	return v.Sub(b.Scale(v.Dot(b) / b.Dot(b)))
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.Dot(v))
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
