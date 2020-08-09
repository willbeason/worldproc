package main

import (
	"flag"
	"fmt"
	"github.com/willbeason/hydrology/pkg/geodesic"
	"github.com/willbeason/hydrology/pkg/noise"
	"github.com/willbeason/hydrology/pkg/planet"
	"github.com/willbeason/hydrology/pkg/render"
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

	p := loadOrCreate(*seed, size, spheres[size])

	screen := render.Screen{
		Width:  2160,
		Height: 1080,
	}
	projection := render.Project(screen, render.Equirectangular{})
	renderImg(*seed, "satellite", projection, spheres, p)
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

func renderImg(seed int64, name string, projection render.Projection, spheres []*geodesic.Geodesic, p *planet.Planet) {
	img := planet.RenderTerrain(p, projection, spheres)
	render.WriteImage(img, fmt.Sprintf("renders/%d-%s.png", seed, name))
}
