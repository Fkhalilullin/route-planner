package models

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Route []Point

type Elevation struct {
	Value             float64 `json:"value,omitempty"`
	Point             Point   `json:"point"`
	Type              string  `json:"type"`
	NeighboringPoints []Point `json:"neighboring_points"`
}

type Elevations []Elevation
