package water

import (
	"fmt"
	"github.com/willbeason/hydrology/pkg/geodesic"
	"math"
	"sort"
)

type IndexHeight struct {
	Index  int
	Height float64
	Water  float64
}

type Lake struct {
	// IndexHeights is the set of Geodesic indices this lake contains and how
	// much water they contain.
	IndexHeights []IndexHeight
	// WaterVolume is the total volume of water
	WaterVolume float64
}

func (l *Lake) Add(i int, h, w float64) {
	l.IndexHeights = append(l.IndexHeights, IndexHeight{
		Index:  i,
		Height: h,
	})
	l.WaterVolume += w
}

func (l *Lake) Merge(other *Lake) {
	if other == nil {
		return
	}
	l.IndexHeights = append(l.IndexHeights, other.IndexHeights...)
	l.WaterVolume += l.WaterVolume
}

func (l *Lake) Equalize() {
	if len(l.IndexHeights) < 2 {
		return
	}

	// Sort heights from lowest to highest.
	sort.Slice(l.IndexHeights, func(i, j int) bool {
		return l.IndexHeights[i].Height < l.IndexHeights[j].Height
	})

	// Now we're going to calculate the exact volume of water which would fill
	// every cell until we reach the volume of water we're looking for.

	i := 0
	volume := 0.0
	waterLevel := 0.0
	maxI := len(l.IndexHeights)
	for i < maxI {
		ih := l.IndexHeights[i]

		newVolume := volume + float64(i)*(ih.Height-waterLevel)
		if newVolume > l.WaterVolume {
			waterLevel += (l.WaterVolume - volume) / float64(i)
			break
		}
		volume = newVolume

		waterLevel = ih.Height
		i++
	}

	// Check for the case that everywhere ends up with water.
	if i == maxI {
		waterLevel = l.WaterVolume / float64(maxI)
	}

	for j := range l.IndexHeights[:i] {
		l.IndexHeights[j].Water = waterLevel - l.IndexHeights[j].Height
	}
}

type VisitOrdinal struct {
	Index int
	Height float64
}

func Equalize(waters, heights []float64, sphere *geodesic.Geodesic) {
	var lakes []Lake

	// Visit nodes from highest to lowest to allow for water flowing downhill.
	toVisit := make([]VisitOrdinal, len(heights))
	for i, h := range heights {
		toVisit[i].Index = i
		toVisit[i].Height = h
	}
	sort.Slice(toVisit, func(i, j int) bool {
		return toVisit[i].Height > toVisit[j].Height
	})

	visited := make(map[int]bool, len(sphere.Centers))
	for _, v := range toVisit {
		if visited[v.Index] {
			continue
		}
		newLake := visitEqualize(v.Index, waters, heights, visited, sphere)
		if newLake != nil {
			lakes = append(lakes, *newLake)
		}
	}

	fmt.Println(len(lakes))
	sort.Slice(lakes, func(i, j int) bool {
		return len(lakes[i].IndexHeights) > len(lakes[j].IndexHeights)
	})
	for _, l := range lakes {
		l.Equalize()
		for _, ih := range l.IndexHeights {
			waters[ih.Index] = ih.Water
		}
	}
}

func visitEqualize(i int, waters, heights []float64, visited map[int]bool, sphere *geodesic.Geodesic) *Lake {
	if visited[i] || waters[i] < 0.001 {
		return nil
	}
	visited[i] = true

	lake := &Lake{}
	lake.Add(i, heights[i], waters[i])

	var toVisit []int
	toVisit = append(toVisit, sphere.Faces[i].Neighbors...)
	nToVisit := len(toVisit)

	minHeight := heights[i]
	for idx := 0; idx < nToVisit; idx++ {
		n := toVisit[idx]
		if visited[n] || (waters[n] < 0.001 && heights[n] > minHeight) {
			continue
		}
		visited[n] = true

		minHeight = math.Min(minHeight, heights[n])

		lake.Add(n, heights[n], waters[n])
		toVisit = append(toVisit, sphere.Faces[n].Neighbors...)
		nToVisit += len(sphere.Faces[n].Neighbors)
	}
	if len(lake.IndexHeights) < 100 {
		return nil
	}

	return lake
}
