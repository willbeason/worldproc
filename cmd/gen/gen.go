package main

import (
	"flag"
	"fmt"
	"github.com/willbeason/hydrology/pkg/climate"
	"github.com/willbeason/hydrology/pkg/geodesic"
	"github.com/willbeason/hydrology/pkg/noise"
	"github.com/willbeason/hydrology/pkg/planet"
	"github.com/willbeason/hydrology/pkg/render"
	"github.com/willbeason/hydrology/pkg/sun"
	"image"
	"math"
	"math/rand"
	"time"
)

var seed = flag.Int64("seed", time.Now().UnixNano(),
	"The seed of the planet to generate")

func main() {
	flag.Parse()
	rand.Seed(*seed)

	size := 9
	spheres := geodesic.New(size, false)

	sphere := spheres[size]
	p := loadOrCreate(*seed, size, sphere)

	climates := make([]climate.Climate, len(p.Heights))
	for i, w := range p.Waters {
		if w > 0.01 {
			climates[i].SpecificHeat = climate.OceanSpecificHeat
		} else if w > 0 {
			climates[i].SpecificHeat = climate.CoastSpecificHeat
		} else {
			climates[i].SpecificHeat = climate.DesertSpecificHeat
		}
		// Initialize to 0 Celsius.
		climates[i].Temperature = 273
	}

	screen := render.Screen{
		Width:  1920,
		Height: 960,
	}
	projection := render.Project(screen, render.Equirectangular{})

	light := &sun.Directional{}

	// Begin simulating every hour.
	idx := 0
	imax := 24
	seconds := 3600.0
	nDiffuse := 18
	for day := 0; day < 362; day++ {
		if day == 358 {
			// Every 10 minutes.
			// Begin two days before rendering begins.
			imax = 144
			seconds = 600.0
			nDiffuse = 3
		}
		if day == 360 {
			fmt.Println("Begin Render")
		}
		for i := 0; i < imax; i++ {
			// Ten minute intervals.
			t := float64(day) + float64(i) / float64(imax)
			fmt.Printf("t = %.03f", t)
			light.Set(t)

			fmt.Print(" ... heat")
			heat(climates, p, sphere, light, seconds)
			for k := 0; k < nDiffuse; k++ {
				fmt.Print(" ... diffuse")
				diffuseHeat(climates, sphere, seconds)
			}
			fmt.Println()

			if day >= 360 || i == 0 {
				// Heat up for a year before rendering.
				RenderTemperature(*seed, idx, projection, spheres, climates)
				idx++
			}
		}
	}
}

func heat(climates []climate.Climate, p *planet.Planet, sphere *geodesic.Geodesic, light *sun.Directional, seconds float64) {
	for i, c := range sphere.Centers {
		flux := climate.Flux * light.Sun.Dot(c)
		flux = math.Max(0.0, flux)
		//before := climates[i].Temperature
		height := p.Heights[i]
		if p.Waters[i] > 0.00 {
			height = 0
		}
		climates[i].Simulate(flux, height, seconds)
		//if i == 19284 {
		//}
		//if climates[i].Temperature > 273 + 40 {
		//	fmt.Println()
		//	fmt.Println(i)
		//	fmt.Println(c)
		//	fmt.Println(light.Sun)
		//	fmt.Println(light.SunAngle)
		//	fmt.Println(math.Acos(light.VisualIntensity(c)))
		//	fmt.Println(before)
		//	fmt.Println(climates[i].Temperature)
		//	fmt.Println("Flux", flux)
		//	panic("HERE")
		//}
	}
}

func diffuseHeat(climates []climate.Climate, sphere *geodesic.Geodesic, seconds float64) {
	// rates is the rate of energy exchange between adjacent climates.
	rates := make([]float64, len(climates))

	for i, c := range climates {
		ns := sphere.Faces[i].Neighbors
		rates[i] -= float64(len(ns)) * c.Temperature
		for _, n := range ns {
			rates[n] += c.Temperature
		}
	}

	for i, c := range climates {
		climates[i].Temperature += rates[i] * 3 * seconds / c.SpecificHeat
	}
}

func loadOrCreate(seed int64, size int, sphere *geodesic.Geodesic) *planet.Planet {
	p := planet.Load(seed, size)
	mutated := false
	if p == nil {
		p = &planet.Planet{}
		mutated = true
	}
	if len(p.Heights) == 0 {
		perlinNoise := noise.NewPerlinFractal(seed, 10, 30, 0.6)
		planet.AddTerrain(p, sphere, perlinNoise)
		mutated = true
	}
	if len(p.Waters) == 0 {
		planet.AddWater(p, 0.5, sphere)
		mutated = true
	}
	if mutated {
		planet.Save(seed, p)
	}
	return p
}

func RenderTemperature(seed int64, idx int, projection render.Projection, spheres []*geodesic.Geodesic, climates []climate.Climate) {
	img := renderClimate(projection, spheres, climates)
	render.WriteImage(img, fmt.Sprintf("renders/temperature/%d-%03d.png", seed, idx))
}

func renderImg(seed int64, name string, projection render.Projection, spheres []*geodesic.Geodesic, light sun.Light, p *planet.Planet) {
	img := planet.RenderTerrain(p, projection, spheres, light)
	render.WriteImage(img, fmt.Sprintf("renders/midsummer/%d-%s.png", seed, name))
}

func renderClimate(projection render.Projection, spheres []*geodesic.Geodesic, climates []climate.Climate) *image.RGBA {
	screen := projection.Screen
	img := image.NewRGBA(image.Rect(0, 0, screen.Width, screen.Height))

	pxTemperatures := make([]float64, screen.Width*screen.Height)

	sphere := spheres[len(spheres)-1]
	for x := 0; x < screen.Width; x++ {
		for y := 0; y < screen.Height; y++ {
			pidx := y*screen.Width + x
			angle := projection.Pixels[pidx]
			v := angle.Vector()
			idx := geodesic.Find(spheres, v)
			dist := math.Sqrt(geodesic.DistSq(v, sphere.Centers[idx]))

			pxT1 := climates[idx].Temperature

			// Linearly interpolate the cell's stats with the second-closest cell.
			idx2 := 0
			distSq2 := math.MaxFloat64
			for _, n := range sphere.Faces[idx].Neighbors {
				nDistSq2 := geodesic.DistSq(v, sphere.Centers[n])
				if nDistSq2 < distSq2 {
					idx2 = n
					distSq2 = nDistSq2
				}
			}
			dist2 := math.Sqrt(distSq2)
			pxT2 := climates[idx2].Temperature

			pxTemperatures[pidx] = render.Lerp(pxT1, pxT2, dist/(dist+dist2))
		}
	}

	screen.PaintTemperature(pxTemperatures, img)
	return img
}
