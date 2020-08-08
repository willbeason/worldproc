package geodesic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const spheresDir = "spheres"

func filename(size int) string {
	return filepath.Join(spheresDir, fmt.Sprintf("sphere-%02d.json", size))
}

func read(size int) *Geodesic {
	file := filename(size)

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil

		}
		panic(err)
	}

	result := &Geodesic{}
	err = json.Unmarshal(bytes, result)
	if err != nil {
		panic(err)
	}
	return result
}

func write(size int, sphere *Geodesic) {
	bytes, err := json.Marshal(sphere)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename(size), bytes, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

// New returns the generated sequence of geodesic spheres.
//
// fromScratch is whether to generate these manually.
// If false, attempts to read from disk.
func New(size int, fromScratch bool) []*Geodesic {
	result := make([]*Geodesic, size+1)

	for i := 0; i < size+1; i++ {
		var sphereI *Geodesic
		if !fromScratch {
			sphereI = read(i)
		}
		if sphereI == nil {
			fmt.Println("Generating Sphere", i)
			if i == 0 {
				sphereI = Dodecahedron()
			} else {
				sphereI = Chamfer(result[i-1])
			}
			if !fromScratch {
				write(i, sphereI)
			}
		} else {
			fmt.Println("Read Sphere", i)
		}
		result[i] = sphereI
	}

	return result
}
