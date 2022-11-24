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
			distance := math.Sqrt(
				(p.Lat-e.Point.Lat)*(p.Lat-e.Point.Lat) +
					(p.Lon-e.Point.Lon)*(p.Lon-e.Point.Lon),
			)
			if distance < minDistance {
				minDistance = distance
				maxPoint = Type{
					Lat:   e.Point.Lat,
					Lon:   e.Point.Lon,
					Value: p.Value,
				}
			}
		}
		maxPoints = append(maxPoints, maxPoint)
	}

	for i, e := range elevations {
		for _, p := range maxPoints {
			if e.Point.Lat == p.Lat && e.Point.Lon == p.Lon {
				elevations[i].SetType(p.Value)
				continue
			}
		}
	}

	return elevations, nil
}

//max = max_integer
//max_value = -1
//for x in array:
//if abs(x_core - x) < max:
//max = abs(x_core - x)
//max_value = x
