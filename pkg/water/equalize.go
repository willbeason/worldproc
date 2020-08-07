package water

import (
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
	landVolume := 0.0
	for i < maxI {
		ih := l.IndexHeights[i]
		landVolume += ih.Height

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
		waterLevel = (l.WaterVolume + landVolume) / float64(maxI)
	}

	for j := range l.IndexHeights[:i] {
		l.IndexHeights[j].Water = waterLevel - l.IndexHeights[j].Height
	}
}

func Equalize(waters, heights []float64, sphere *geodesic.Geodesic) {
	var lakes []Lake

	// Visit nodes from lowest to highest.
	toVisit := make([]Ordinal, len(heights))
	for i, h := range heights {
		toVisit[i].Index = i
		toVisit[i].Height = h
	}
	sort.Slice(toVisit, func(i, j int) bool {
		return toVisit[i].Height < toVisit[j].Height
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

	sort.Slice(lakes, func(i, j int) bool {
		return len(lakes[i].IndexHeights) > len(lakes[j].IndexHeights)
	})
	for _, l := range lakes {
		l.Equalize()
		for _, ih := range l.IndexHeights {
			waters[ih.Index] += ih.Water
		}
	}
}

func visitEqualize(i int, waters, heights []float64, visited map[int]bool, sphere *geodesic.Geodesic) *Lake {
	if visited[i] {
		return nil
	}
	visited[i] = true

	l := &Lake{}

	var toVisit OrdinalList
	hi := heights[i]
	toVisit.Insert(Ordinal{Index: i, Height: hi})
	// We are guaranteed that every neighbor is the same level or higher.
	for _, n := range sphere.Faces[i].Neighbors {
		toVisit.Insert(Ordinal{Index: n, Height: hi})
	}

	for cell := toVisit.Pop(); cell != nil; cell = toVisit.Pop() {
		if i != cell.Index && visited[cell.Index] {
			continue
		}
		visited[cell.Index] = true

		// cell.Height actually records the minimum height of the water we may
		// take from this cell.
		hc := heights[cell.Index]
		w := math.Max(0.0, hc+waters[cell.Index]-cell.Height)
		w = math.Min(w, waters[cell.Index])
		waters[cell.Index] -= w
		l.Add(cell.Index, math.Max(hc, cell.Height), w)

		hi := math.Max(cell.Height, hc)
		for _, n := range sphere.Faces[cell.Index].Neighbors {
			if visited[n] {
				// Don't add already-visited cells.
				continue
			}
			if heights[n]+waters[n]-hi < 0.001 {
				// Don't add cells we can't possibly take water from.
				continue
			}
			toVisit.Insert(Ordinal{Index: n, Height: hi})
		}
	}

	return l
}
