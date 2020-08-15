package render

import (
	"github.com/willbeason/hydrology/pkg/geodesic"
	"image"
	"image/color"
	"math"
	"sync"
)

type Screen struct {
	Width, Height int
}

var (
	deepWater   = color.RGBA{14, 31, 75, 255}
	sand        = color.RGBA{246, 223, 57, 255}
	brightGreen = color.RGBA{109, 150, 74, 255}
	green       = color.RGBA{67, 117, 63, 255}
	darkGreen   = color.RGBA{50, 97, 43, 255}
	stone       = color.RGBA{131, 115, 79, 255}
	snow        = color.RGBA{255, 255, 255, 255}
)

var landWaterCS = NewColorScale(
	[]ColorPoint{
		{-0.1, deepWater},
		{-0.01, color.RGBA{0, 32, 255, 255}},
		{0, color.RGBA{0, 128, 255, 255}},
		{0.02, color.RGBA{192, 192, 32, 255}},
		{0.1, color.RGBA{0, 192, 0, 255}},
		{0.15, color.RGBA{64, 128, 64, 255}},
		{0.2, color.RGBA{64, 96, 64, 255}},
		{0.5, color.RGBA{192, 192, 192, 255}},
		{1.0, color.RGBA{255, 255, 255, 255}},
	})

func (s Screen) PaintAndShade(heights []float64, lightAngles []geodesic.Angle, img *image.RGBA) {
	wg := sync.WaitGroup{}
	wg.Add(s.Width)

	for x := 0; x < s.Width; x++ {
		xl := x
		go func() {
			for y := 0; y < s.Height; y++ {
				idx := y*s.Width + xl
				h := heights[idx]

				h = h / 2.6
				if h > 1.0 || h < -1.0 {
					panic(h)
				}

				c := landWaterCS.ColorAt(h)
				c = s.shadow(c, heights, xl, y, idx, lightAngles[idx])

				img.Set(xl, y, c)
			}
			wg.Done()
		}()

	}
	wg.Wait()
}

func (s Screen) shadow(c color.RGBA, heights []float64, x, y, idx int, lightAngle geodesic.Angle) color.RGBA {
	if x <= 0 || x >= s.Width-1 || y <= 0 || y >= s.Height-1 {
		return c
	}

	hl := heights[idx-1]
	hr := heights[idx+1]
	hu := heights[idx+s.Width]
	hd := heights[idx-s.Width]
	m := 0.1
	if lightAngle.Theta > 0 {
		m = gradient(hl, hr, hu, hd).Dot(lightAngle.Vector())
		m = math.Max(0.1, m)
	}

	c.R = uint8(float64(c.R) * m)
	c.G = uint8(float64(c.G) * m)
	c.B = uint8(float64(c.B) * m)

	return c
}

var landCS = NewColorScale(
	[]ColorPoint{
		{0, sand},
		{0.05, brightGreen},
		{0.3, green},
		{0.5, darkGreen},
		{0.7, stone},
		{1.0, snow},
	})

func (s Screen) PaintLandWater(heights, waters, lights []float64, lightAngles []geodesic.Angle, img *image.RGBA) {
	for x := 0; x < s.Width; x++ {
		for y := 0; y < s.Height; y++ {
			idx := y*s.Width + x

			w := waters[idx]
			h := heights[idx]

			c := landCS.ColorAt(h)

			if w > 0.01 {
				c = deepWater
				c = lerpC(color.RGBA{A: 255}, c, lights[idx])
			} else {
				c = s.shadow(c, heights, x, y, idx, lightAngles[idx])
				if w > 0.0 {
					c = lerpC(c, deepWater, w/0.01)
				}
			}
			img.Set(x, y, c)
		}
	}
}

var temperatureCS = NewColorScale(
	[]ColorPoint{
		{223, color.RGBA{R: 255, G: 255, B: 255, A: 255}}, // -50 C
		{233, color.RGBA{R: 255, G: 0, B: 255, A: 255}}, // -40 C
		{243, color.RGBA{R: 128, G: 0, B: 255, A: 255}}, // -30 C
		{253, color.RGBA{R: 0, G: 0, B: 255, A: 255}}, // -20 C
		{263, color.RGBA{R: 0, G: 255, B: 255, A: 255}}, // -10 C
		{273, color.RGBA{R: 0, G: 255, B: 0, A: 255}}, // 0 C
		{283, color.RGBA{R: 128, G: 255, B: 0, A: 255}}, // 10 C
		{293, color.RGBA{R: 255, G: 255, B: 0, A: 255}}, // 20 C
		{303, color.RGBA{R: 255, G: 128, B: 0, A: 255}}, // 30 C
		{313, color.RGBA{R: 255, G: 0, B: 0, A: 255}}, // 40 C
		{323, color.RGBA{R: 128, G: 0, B: 0, A: 255}}, // 50 C
	})

func (s Screen) PaintTemperature(temperatures []float64, img *image.RGBA) {
	for x := 0; x < s.Width; x++ {
		for y := 0; y < s.Height; y++ {
			idx := y*s.Width + x

			t := temperatures[idx]

			c := temperatureCS.ColorAt(t)

			img.Set(x, y, c)
		}
	}
}

func (s Screen) Paint(heights []float64, cs *ColorScale, img *image.RGBA) {
	wg := sync.WaitGroup{}
	wg.Add(s.Width)

	for x := 0; x < s.Width; x++ {
		xl := x
		go func() {
			for y := 0; y < s.Height; y++ {
				idx := y*s.Width + xl
				h := heights[idx]

				c := cs.ColorAt(h)
				img.Set(xl, y, c)
			}
			wg.Done()
		}()

	}
	wg.Wait()
}

func gradient(l, r, u, d float64) geodesic.Vector {
	dx := l - r
	dy := d - u
	return geodesic.Vector{X: dx, Y: dy, Z: 0.2}.Normalize()
}
