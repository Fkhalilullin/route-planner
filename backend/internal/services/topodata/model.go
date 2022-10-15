package topodata

import "github.com/Fkhalilullin/route-planner/internal/route"

type Elevation struct {
	Value float64     `json:"value,omitempty"`
	Point route.Point `json:"point"`
}

type Elevations []Elevation

type elevationResponse struct {
	Results []struct {
		Dataset   string  `json:"dataset"`
		Elevation float64 `json:"elevation"`
		Location  struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"results"`
	Status string `json:"status"`
}
