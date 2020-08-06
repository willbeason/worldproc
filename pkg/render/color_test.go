package render

import (
	"github.com/google/go-cmp/cmp"
	"image/color"
	"testing"
)

func TestColorScale_ColorAt(t *testing.T) {
	cs := NewColorScale([]ColorPoint{
		{-1.0, color.RGBA{255, 0, 0, 255}},
		{0.0, color.RGBA{0, 255, 0, 255}},
		{1.0, color.RGBA{0, 0, 255, 255}},
	})

	tcs := []struct{
		p float64
		want color.RGBA
	} {
		{-2.0, color.RGBA{255, 0, 0, 255}},
		{-1.0, color.RGBA{255, 0, 0, 255}},
		{0.0, color.RGBA{0, 255, 0, 255}},
		{1.0, color.RGBA{0, 0, 255, 255}},
		{2.0, color.RGBA{0, 0, 255, 255}},
	}

	for _, tc := range tcs {
		got := cs.ColorAt(tc.p)

		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Error(tc.p, diff)
		}
	}
}

func TestColorScale_Grays(t *testing.T) {
	cs := NewColorScale([]ColorPoint{
		{0.35, color.RGBA{255, 255, 255, 255}},
		{0.39, color.RGBA{0, 0, 0, 255}},
	})

	tcs := []struct{
		p float64
		want color.RGBA
	} {
		{p: 0.0, want: color.RGBA{255, 255, 255, 255},},
		{p: 0.37, want: color.RGBA{127, 127, 127, 255},},
		{p: 1.0, want: color.RGBA{0, 0, 0, 255},},
	}

	for _, tc := range tcs {
		got := cs.ColorAt(tc.p)

		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Error(tc.p, diff)
		}
	}
}
