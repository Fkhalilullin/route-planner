package models

import (
	"github.com/Fkhalilullin/route-planner/internal/config"
	"github.com/Fkhalilullin/route-planner/internal/route"
	"math"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Elevation struct {
	Value             float64      `json:"value,omitempty"`
	Point             Point        `json:"point"`
	Type              string       `json:"type"`
	NeighboringPoints []*Elevation `json:"neighboring_points"`
	X                 int          `json:"X"`
	Y                 int          `json:"Y"`
}

func (e *Elevation) getNeighboringPoints() []*Elevation {
	var bufElevations []*Elevation

	for i, ee := range Mesh {
		for j := range ee {
			if Mesh[i][j].X == e.X && Mesh[i][j].Y == e.Y {
				if i-1 >= 0 && j-1 >= 0 {
					bufElevations = append(bufElevations, Mesh[i-1][j-1])
				}
				if i+1 < len(Mesh) && j+1 < len(ee) {
					bufElevations = append(bufElevations, Mesh[i+1][j+1])
				}
				if i-1 >= 0 && j+1 < len(ee) {
					bufElevations = append(bufElevations, Mesh[i-1][j+1])
				}
				if i+1 < len(Mesh) && j-1 >= 0 {
					bufElevations = append(bufElevations, Mesh[i+1][j-1])
				}
				if i+1 < len(Mesh) {
					bufElevations = append(bufElevations, Mesh[i+1][j])
				}
				if i-1 >= 0 {
					bufElevations = append(bufElevations, Mesh[i-1][j])
				}
				if j+1 < len(ee) {
					bufElevations = append(bufElevations, Mesh[i][j+1])
				}
				if j-1 >= 0 {
					bufElevations = append(bufElevations, Mesh[i][j-1])
				}
				break
			}
		}
	}

	return bufElevations
}

func (e *Elevation) PathNeighbors() []route.Pather {
	var neighbors []route.Pather
	for _, n := range e.getNeighboringPoints() {
		neighbors = append(neighbors, n)
	}
	return neighbors
}

// PathNeighborCost returns the movement cost of the directly neighboring tile.
func (e *Elevation) PathNeighborCost(to route.Pather) float64 {
	toT := to.(*Elevation)

	switch toT.Type {
	case config.TypeLand:
		return 1
	case config.TypeForest:
		return 2
	case config.TypeWater:
		return 3
	}
	return 1
}

func (e *Elevation) PathEstimatedCost(to route.Pather) float64 {
	toT := to.(*Elevation)

	absX := (toT.X - e.X) * (toT.X - e.X)
	absY := (toT.Y - e.Y) * (toT.Y - e.Y)

	return math.Sqrt(float64(absX + absY))
}

func (e *Elevation) SetType(newType string) {
	e.Type = newType
}

type Elevations [][]*Elevation

var Mesh Elevations
