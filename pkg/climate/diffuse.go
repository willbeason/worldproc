package climate

import "github.com/willbeason/hydrology/pkg/geodesic"

func DiffuseAir(climates []Climate, sphere *geodesic.Geodesic) {
	airRates := make([]float64, len(climates))

	for i, c := range climates {
		ns := sphere.Faces[i].Neighbors
		airRates[i] -= float64(len(ns)) * c.Pressure()
		for _, n := range ns {
			airRates[n] += c.Pressure()
		}
	}

	for i, r := range airRates {
		climates[i].Air += r * 0.005
	}
}
