package pather

import (
	"container/heap"
	"github.com/Fkhalilullin/route-planner/internal/config"
	"math"
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

	for i, ee := range Mesh {
		for j := range ee {
			if Mesh[i][j].X == c.X && Mesh[i][j].Y == c.Y {
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

func (c *Coordinate) PathNeighbors() []Pather {
	var neighbors []Pather
	for _, n := range c.getNeighboringPoints() {
		neighbors = append(neighbors, n)
	}
	return neighbors
}

func (c *Coordinate) PathNeighborCost(to Pather) float64 {
	toT := to.(*Coordinate)

	switch toT.Type {
	case config.TypeLand:
		return 1
	case config.TypeForest:
		return 3
	case config.TypeWater:
		return 2
	}
	return 1
}

func (c *Coordinate) PathEstimatedCost(to Pather) float64 {
	toT := to.(*Coordinate)

	//absLat := (toT.Point.Lat - c.Point.Lat) * (toT.Point.Lat - c.Point.Lat)
	//absLon := (toT.Point.Lon - c.Point.Lon) * (toT.Point.Lon - c.Point.Lon)
	absLat := (toT.X - c.X) * (toT.X - c.X)
	absLon := (toT.Y - c.Y) * (toT.Y - c.Y)
	absElevation := (toT.Value - c.Value) * (toT.Value - c.Value)

	return math.Sqrt(float64(absLat+absLon) + absElevation)
}

func (c *Coordinate) SetType(newType string) {
	c.Type = newType
}
