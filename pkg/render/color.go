package render

import (
	"image/color"
)

type ColorPoint struct {
	Threshold float64
	Color color.RGBA
}

type ColorScale struct {
	Thresholds []float64
	Colors []color.RGBA
}

func NewColorScale(ps []ColorPoint) *ColorScale {
	result := &ColorScale{
		Thresholds: make([]float64, len(ps)),
		Colors: make([]color.RGBA, len(ps)),
	}
	for i, p := range ps {
		result.Thresholds[i] = p.Threshold
		result.Colors[i] = p.Color
	}

	return result
}

func Lerp(a0, a1, w float64) float64 {
	return (1 - w)*a0 + w*a1
}

// lerpC linearly interpolates colors.
func lerpC(left, right color.RGBA, w float64) color.RGBA {
	return color.RGBA{
		R: uint8(Lerp(float64(left.R), float64(right.R), w)),
		G: uint8(Lerp(float64(left.G), float64(right.G), w)),
		B: uint8(Lerp(float64(left.B), float64(right.B), w)),
		A: uint8(Lerp(float64(left.A), float64(right.A), w)),
	}
}

func (s ColorScale) ColorAt(f float64) color.RGBA {
	if f < s.Thresholds[0] {
		// We're before the first threshold, so use the first color.
		return s.Colors[0]
	}

	for rightIdx := 1; rightIdx < len(s.Thresholds); rightIdx++ {
		rightT := s.Thresholds[rightIdx]
		if f < rightT {
			// Linearly interpolate between the threshold colors.
			leftT := s.Thresholds[rightIdx-1]
			p := (f - leftT) / (rightT - leftT)
			return lerpC(s.Colors[rightIdx-1], s.Colors[rightIdx], p)
		}
	}

	// We're above the highest threshold, so use the last color.
	return s.Colors[len(s.Colors)-1]
}
