package osm

import (
	"context"
	"fmt"
	"github.com/Fkhalilullin/route-planner/internal/models"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmapi"
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
	bounds := &osm.Bounds{
		//MinLon: 50.6874, MaxLon: 50.7074,
		//MinLat: 63.9391, MaxLat: 63.9454,
		MinLat: box.MinLat, MaxLat: box.MaxLat,
		MinLon: box.MinLon, MaxLon: box.MaxLon,
	}

	o, err := osmapi.Map(context.Background(), bounds) // fetch data from the osm api.
	if err != nil {
		return nil, fmt.Errorf("openstreetmap.GetTypePoints failed http GET: %w", err)
	}

	fmt.Println(o)

	return nil, nil
}
