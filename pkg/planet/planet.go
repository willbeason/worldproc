package planet

type Planet struct {
	Size int `json:"size"`

	Heights []float64 `json:"heights"`
	Waters []float64 `json:"waters,omitempty"`
	Flows []float64 `json:"flows,omitempty"`
	Temperatures []float64 `json:"temperatures,omitempty"`
}
