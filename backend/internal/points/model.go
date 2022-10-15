package points

import "github.com/Fkhalilullin/route-planner/internal/route"

type Elevation struct {
	Value int64       `json:"value,omitempty"`
	Point route.Point `json:"point"`
}

type Elevations []Elevation
