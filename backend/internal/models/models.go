package models

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Route []Point

type Elevation struct {
	Value float64 `json:"value,omitempty"`
	Point Point   `json:"point"`
}

type Elevations []Elevation
