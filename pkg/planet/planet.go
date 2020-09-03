package planet

import "github.com/willbeason/worldproc/pkg/climate"

type Planet struct {
	Size int `json:"size"`

	Heights []float64 `json:"heights"`
	Waters []float64 `json:"waters,omitempty"`
	Flows []float64 `json:"flows,omitempty"`
	Climates []climate.Climate `json:"temperatures,omitempty"`
}
