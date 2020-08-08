package main

import (
	"encoding/json"
	"fmt"
	"github.com/willbeason/hydrology/pkg/geodesic"
	"github.com/willbeason/hydrology/pkg/noise"
	"github.com/willbeason/hydrology/pkg/planet"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func main() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	spheres := geodesic.New(8, false)

	p := planet.Planet{}

	perlinNoise := noise.NewPerlinFractal(10, 30, 0.6)
	sphere := spheres[len(spheres)-1]

	p.Heights = make([]float64, len(sphere.Centers))
	for cell, pos := range sphere.Centers {
		p.Heights[cell] = perlinNoise.ValueAt(pos)
	}

	bytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filepath.Join("planets", fmt.Sprintf("planet-%d.json", seed)),
		bytes, os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(spheres))
}
