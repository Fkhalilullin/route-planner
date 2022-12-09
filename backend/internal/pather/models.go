package pather

import "github.com/Fkhalilullin/route-planner/internal/models"

type Coordinates [][]*Coordinate

type Coordinate struct {
	Value float64      `json:"value,omitempty"`
	Point models.Point `json:"point"`
	Type  string       `json:"type"`

	X int
	Y int
}

var Mesh Coordinates

type node struct {
	pather Pather
	cost   float64
	rank   float64
	parent *node
	open   bool
	closed bool
	index  int
}

type nodeMap map[Pather]*node

func (nm nodeMap) get(p Pather) *node {
	n, ok := nm[p]
	if !ok {
		n = &node{
			pather: p,
		}
		nm[p] = n
	}
	return n
}
