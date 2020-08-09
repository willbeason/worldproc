package planet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Save(seed int64, planet *Planet) {
	bytes, err := json.Marshal(planet)
	if err != nil {
		panic(err)
	}
	name := fmt.Sprintf("%d.json", seed)
	err = ioutil.WriteFile(filepath.Join("planets", name), bytes, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func Load(seed int64, size int) *Planet {
	p := &Planet{}
	name := fmt.Sprintf("%d.json", seed)
	bytes, err := ioutil.ReadFile(filepath.Join("planets", name))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, p)
	if err != nil {
		panic(err)
	}

	//if size > p.Size {
	//	panic(fmt.Sprintf("requested size %d larger than rendered size %d", size, p.Size))
	//}

	nFaces := 1
	for i := 0; i < size; i++ {
		nFaces *= 4
	}
	nFaces = 10*nFaces + 2

	p.Size = size
	if len(p.Heights) > 0 {
		p.Heights = p.Heights[:nFaces]
	}
	if len(p.Waters) > 0 {
		p.Waters = p.Waters[:nFaces]
	}
	if len(p.Flows) > 0 {
		p.Flows = p.Flows[:nFaces]
	}
	if len(p.Temperatures) > 0 {
		p.Temperatures = p.Temperatures[:nFaces]
	}

	return p
}
