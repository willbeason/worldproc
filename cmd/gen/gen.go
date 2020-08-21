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
	screen := render.Screen{
		Width:  1920,
		Height: 960,
	}
	projection := render.Project(screen, render.Equirectangular{})
	renderImg(*seed, "", projection, spheres, sun.Constant{}, p)

	if len(p.Climates) == 0 {
		fmt.Println("Initializing Climate")
		initializeClimate(p, sphere, spheres, projection)
		planet.Save(*seed, p)
	} else {
		fmt.Println("Loaded Climate")
	}

	imax := 144
	seconds := 600.0
	nDiffuse := 1
	nWind := 5
	light := &sun.Directional{}
	idx := 0
	for day := 0; day < 20; day++ {
		for i := 0; i < imax; i++ {
			t := float64(day) + float64(i) / float64(imax)
			fmt.Printf("t = %.03f", t)
			light.Set(t)

			fmt.Print(" ... heat")
			heat(p.Climates, p, sphere, light, seconds)
			for k := 0; k < nWind; k++ {
				fmt.Print(" ... wind")
				climate.Flow(p.Climates, sphere, 2.0)
			}
			for k := 0; k < nDiffuse; k++ {
				fmt.Print(" ... diffuse")
				diffuseHeat(p.Climates, sphere, seconds)
			}
			fmt.Println()

			// Heat up for a year before rendering.
			RenderClimate(*seed, idx, projection, spheres, p.Climates)
			idx++

			printAveragePressure(p.Climates)
		}
	}
}

func printAveragePressure(climates []climate.Climate) {
	total := 0.0
	totV := 0.0
	for _, c := range climates {
		total += c.Pressure()
		totV += c.AirVelocity.Length()
	}
	fmt.Printf("Mean Pressure: %.04f\n", total / float64(len(climates)))
	fmt.Printf("Mean Velocity: %.04f\n", totV / float64(len(climates)))
}

func initializeClimate(p *planet.Planet, sphere *geodesic.Geodesic, spheres []*geodesic.Geodesic, projection render.Projection) {
	light := &sun.Directional{}
	p.Climates = make([]climate.Climate, len(p.Heights))
	for i, w := range p.Waters {
		if w > 0.01 {
			p.Climates[i].LandSpecificHeat = climate.OceanSpecificHeat
		} else if w > 0 {
			p.Climates[i].LandSpecificHeat = climate.CoastSpecificHeat
		} else {
			p.Climates[i].LandSpecificHeat = climate.DesertSpecificHeat
		}
		p.Climates[i].Air = 1.0
		// Initialize to 0 Celsius.
		p.Climates[i].SetTemperature(climate.ZeroCelsius)
	}
	// Begin simulating every hour.
	idx := 0
	imax := 24
	seconds := 3600.0
	nDiffuse := 6
	for day := 0; day < 360; day++ {
		for i := 0; i < imax; i++ {
			t := float64(day) + float64(i) / float64(imax)
			fmt.Printf("t = %.03f", t)
			light.Set(t)

			fmt.Print(" ... heat")
			heat(p.Climates, p, sphere, light, seconds)
			for k := 0; k < nDiffuse; k++ {
				fmt.Print(" ... diffuse")
				diffuseHeat(p.Climates, sphere, seconds)
			}
			fmt.Println()

			if day >= 360 || i == 0 {
				// Heat up for a year before rendering.
				RenderClimate(*seed, idx, projection, spheres, p.Climates)
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
		//fmt.Println(climates[i])
		//fmt.Println(climates[i].AirTemperature(), climates[i].LandTemperature())
		climates[i].Simulate(flux, math.Asin(sphere.Centers[i].Z), height, seconds)
		//fmt.Println(climates[i])
		//fmt.Println(climates[i].AirTemperature(), climates[i].LandTemperature())
		//if climates[i].AirTemperature() > 50 + climate.ZeroCelsius {
		//	panic(i)
		//}
	}
}

func diffuseHeat(climates []climate.Climate, sphere *geodesic.Geodesic, seconds float64) {
	// rates is the rate of energy exchange between adjacent climates.
	landRates := make([]float64, len(climates))
	airRates := make([]float64, len(climates))

	for i, c := range climates {
		ns := sphere.Faces[i].Neighbors
		landRates[i] -= float64(len(ns)) * c.LandTemperature()
		airRates[i] -= float64(len(ns)) * c.AirTemperature()
		for _, n := range ns {
			landRates[n] += c.LandTemperature()
			airRates[n] += c.AirTemperature()
		}
	}

	for i := range climates {
		climates[i].LandEnergy += landRates[i] * 3 * seconds
		climates[i].AirEnergy += airRates[i] * 3 * seconds
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

func RenderClimate(seed int64, idx int, projection render.Projection, spheres []*geodesic.Geodesic, climates []climate.Climate) {
	img, img2, img3 := renderClimate(projection, spheres, climates)
	n := 10
	render.WriteImage(img, fmt.Sprintf("renders/wind-test-%d/temp-%d-%03d.png", n, seed, idx))
	render.WriteImage(img2, fmt.Sprintf("renders/wind-test-%d/wind-%d-%03d.png", n, seed, idx))
	render.WriteImage(img3, fmt.Sprintf("renders/wind-test-%d/pressure-%d-%03d.png", n, seed, idx))
}

func renderImg(seed int64, name string, projection render.Projection, spheres []*geodesic.Geodesic, light sun.Light, p *planet.Planet) {
	img := planet.RenderTerrain(p, projection, spheres, light)
	render.WriteImage(img, fmt.Sprintf("renders/%d-%s.png", seed, name))
}

func renderClimate(projection render.Projection, spheres []*geodesic.Geodesic, climates []climate.Climate) (*image.RGBA, *image.RGBA, *image.RGBA) {
	screen := projection.Screen

	pxTemperatures := make([]float64, screen.Width*screen.Height)
	pxAirVelocities := make([]float64, screen.Width*screen.Height)
	pxAirPressures := make([]float64, screen.Width*screen.Height)

	sphere := spheres[len(spheres)-1]
	for x := 0; x < screen.Width; x++ {
		for y := 0; y < screen.Height; y++ {
			pidx := y*screen.Width + x
			angle := projection.Pixels[pidx]
			v := angle.Vector()
			idx := geodesic.Find(spheres, v)
			dist := math.Sqrt(geodesic.DistSq(v, sphere.Centers[idx]))

			pxT1 := climates[idx].AirTemperature()
			pxA1 := climates[idx].AirVelocity.Length()
			pxP1 := climates[idx].Pressure()

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
			pxT2 := climates[idx2].AirTemperature()
			pxA2 := climates[idx2].AirVelocity.Length()
			pxP2 := climates[idx2].Pressure()

			pxTemperatures[pidx] = render.Lerp(pxT1, pxT2, dist/(dist+dist2))
			pxAirVelocities[pidx] = render.Lerp(pxA1, pxA2, dist/(dist+dist2))
			pxAirPressures[pidx] = render.Lerp(pxP1, pxP2, dist/(dist+dist2))
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, screen.Width, screen.Height))
	screen.PaintTemperature(pxTemperatures, img)

	img2 := image.NewRGBA(image.Rect(0, 0, screen.Width, screen.Height))
	screen.PaintAirVelocities(pxAirVelocities, img2)

	img3 := image.NewRGBA(image.Rect(0, 0, screen.Width, screen.Height))
	screen.PaintAirPressure(pxAirPressures, img3)

	return img, img2, img3
}
