package pather

import (
	"container/heap"
	"math"

	"github.com/Fkhalilullin/route-planner/internal/config"
)

type Pather interface {
	PathNeighbors() []Pather
	PathNeighborCost(to Pather) float64
	PathEstimatedCost(to Pather) float64
}

func Path(from, to Pather) (path []Pather, distance float64, found bool) {
	nm := nodeMap{}
	nq := &priorityQueue{}
	heap.Init(nq)
	fromNode := nm.get(from)
	fromNode.open = true
	heap.Push(nq, fromNode)
	for {
		if nq.Len() == 0 {
			return
		}
		current := heap.Pop(nq).(*node)
		current.open = false
		current.closed = true

		if current.pather == nm.get(to).pather {
			p := []Pather{}
			curr := current
			for curr != nil {
				p = append(p, curr.pather)
				curr = curr.parent
			}
			return p, current.cost, true
		}

		for _, neighbor := range current.pather.PathNeighbors() {
			cost := current.cost + current.pather.PathNeighborCost(neighbor)
			neighborNode := nm.get(neighbor)
			if cost < neighborNode.cost {
				if neighborNode.open {
					heap.Remove(nq, neighborNode.index)
				}
				neighborNode.open = false
				neighborNode.closed = false
			}
			if !neighborNode.open && !neighborNode.closed {
				neighborNode.cost = cost
				neighborNode.open = true
				neighborNode.rank = cost + neighbor.PathEstimatedCost(to)
				neighborNode.parent = current
				heap.Push(nq, neighborNode)
			}
		}
	}
}

func (c *Coordinate) getNeighboringPoints() []*Coordinate {
	var bufElevations []*Coordinate

	for _, ee := range Mesh {
		if c.X+1 < len(Mesh) {
			bufElevations = append(bufElevations, Mesh[c.X+1][c.Y])
		}
		if c.X-1 >= 0 {
			bufElevations = append(bufElevations, Mesh[c.X-1][c.Y])
		}
		if c.Y+1 < len(ee) {
			bufElevations = append(bufElevations, Mesh[c.X][c.Y+1])
		}
		if c.Y-1 >= 0 {
			bufElevations = append(bufElevations, Mesh[c.X][c.Y-1])
		}

		if c.X-1 >= 0 && c.Y-1 >= 0 {
			bufElevations = append(bufElevations, Mesh[c.X-1][c.Y-1])
		}
		if c.X+1 < len(Mesh) && c.Y+1 < len(ee) {
			bufElevations = append(bufElevations, Mesh[c.X+1][c.Y+1])
		}

		if c.X+1 < len(Mesh) && c.Y-1 >= 0 {
			bufElevations = append(bufElevations, Mesh[c.X+1][c.Y-1])
		}
		if c.X-1 >= 0 && c.Y+1 < len(ee) {
			bufElevations = append(bufElevations, Mesh[c.X-1][c.Y+1])
		}

		break
	}

	return bufElevations
}

func (c *Coordinate) PathNeighbors() []Pather {
	var neighbors []Pather
	for _, n := range c.getNeighboringPoints() {
		neighbors = append(neighbors, n)
	}
	return neighbors
}

func (c *Coordinate) PathNeighborCost(to Pather) float64 {
	toT := to.(*Coordinate)

	diagonalMove := 1.0
	if toT.X != c.X && toT.Y != c.Y {
		diagonalMove = 1.42412
	}

	switch toT.Type {
	case config.TypeLand:
		return config.LandCost * diagonalMove
	case config.TypeForest:
		return config.ForestCost * diagonalMove
	case config.TypeWater:
		return config.WaterCost * diagonalMove
	}
	return 1
}

func (c *Coordinate) PathEstimatedCost(to Pather) float64 {
	toT := to.(*Coordinate)

	// Евклид
	absLat := (toT.Point.Lat - c.Point.Lat) * (toT.Point.Lat - c.Point.Lat)
	absLon := (toT.Point.Lon - c.Point.Lon) * (toT.Point.Lon - c.Point.Lon)
	absElevation := (toT.Value - c.Value) * (toT.Value - c.Value)

	// Чебушев
	//absLat := math.Abs(toT.Point.Lat - c.Point.Lat)
	//absLon := math.Abs(toT.Point.Lon - c.Point.Lon)

	return math.Sqrt(absLat + absLon + absElevation)
	//return math.Max(absLat, absLon)
}

func (c *Coordinate) SetType(newType string) {
	c.Type = newType
}
