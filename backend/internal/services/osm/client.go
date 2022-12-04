package osm

import (
	"context"
	"fmt"
	"github.com/Fkhalilullin/route-planner/internal/config"
	"github.com/Fkhalilullin/route-planner/internal/models"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmapi"
	"math"
)

type OpenStreetMapProvider interface {
	GetTypePoints(box Box) error
}

type service struct {
	openStreetMapProvider OpenStreetMapProvider
}

func NewService() *service {
	return &service{}
}

func (s *service) GetTypePoints(elevations models.Elevations, box Box) (models.Elevations, error) {
	var (
		pointsByID   = make(map[int64]models.Point)
		pointsType   []Type
		typeOfNature string
	)

	bounds := &osm.Bounds{
		MinLat: box.MinLat, MaxLat: box.MaxLat,
		MinLon: box.MinLon, MaxLon: box.MaxLon,
	}

	o, err := osmapi.Map(context.Background(), bounds) // fetch data from the osm api.
	if err != nil {
		return nil, fmt.Errorf("openstreetmap.GetTypePoints failed http GET: %w", err)
	}

	for _, n := range o.Nodes {
		pointsByID[int64(n.ID)] = models.Point{
			Lat: n.Lat,
			Lon: n.Lon,
		}
	}

	for _, w := range o.Ways {

		switch w.Tags.Find(config.KeyNatural) {
		case config.TypeWater:
			typeOfNature = config.TypeWater
		case config.TypeForest:
			typeOfNature = config.TypeForest
		default:
			typeOfNature = config.TypeLand
			continue
		}

		for _, n := range w.Nodes {
			v, ok := pointsByID[int64(n.ID)]
			if ok {
				pointsType = append(pointsType, Type{Lat: v.Lat, Lon: v.Lon, Value: typeOfNature})
			}
		}
	}

	var maxPoints []Type
	for _, p := range pointsType {
		var minDistance = math.MaxFloat64
		var maxPoint Type
		for _, e := range elevations {
			for _, ee := range e {
				distance := math.Sqrt(
					(p.Lat-ee.Point.Lat)*(p.Lat-ee.Point.Lat) +
						(p.Lon-ee.Point.Lon)*(p.Lon-ee.Point.Lon),
				)
				if distance < minDistance {
					minDistance = distance
					maxPoint = Type{
						Lat:   ee.Point.Lat,
						Lon:   ee.Point.Lon,
						Value: p.Value,
					}
				}
			}
		}
		maxPoints = append(maxPoints, maxPoint)
	}

	for i, e := range elevations {
		for j, ee := range e {
			for _, p := range maxPoints {
				if ee.Point.Lat == p.Lat && ee.Point.Lon == p.Lon {
					elevations[i][j].SetType(p.Value)
					continue
				}

				if GetPolyPoints(maxPoints, ee.Point.Lon, ee.Point.Lat) {
					elevations[i][j].SetType(p.Value)
				}
			}
		}
	}

	return elevations, nil
}

func GetPolyPoints(vertices []Type, lon float64, lat float64) bool {
	var (
		collision bool
		next      int
	)

	for current := 0; current < len(vertices); current++ {
		next = current + 1
		if next == len(vertices) {
			next = 0
		}

		vc := vertices[current]
		vn := vertices[next]

		if ((vc.Lat >= lat && vn.Lat < lat) || (vc.Lat < lat && vn.Lat >= lat)) &&
			(lon < (vn.Lon-vc.Lon)*(lat-vc.Lon)/(vn.Lat-vc.Lat)+vc.Lon) {
			collision = !collision
		}
	}

	return collision
}
