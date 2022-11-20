package topodata

import (
	"encoding/json"
	"fmt"
	"github.com/Fkhalilullin/route-planner/internal/config"
	"github.com/Fkhalilullin/route-planner/internal/services/osm"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Fkhalilullin/route-planner/internal/models"
)

type ElevationProvider interface {
	GetElevationPoints(coordinateList string) (models.Elevations, error)
}

type service struct {
	elevationProvider ElevationProvider
}

func NewService() *service {
	return &service{}
}

const endpoint = "https://api.opentopodata.org/v1/srtm90m"

type ElevationRequest struct {
	Locations string `json:"locations,omitempty"`
}

func (s *service) GetElevationPoints(coordinateList string, box osm.Box) (models.Elevations, error) {
	reqByte, err := json.Marshal(ElevationRequest{Locations: coordinateList})
	reader := strings.NewReader(string(reqByte))

	res, err := http.Post(endpoint, "application/json", reader)
	if err != nil {
		return nil, fmt.Errorf("opentopodata.GetElevationPoints failed http GET: %w", err)
	}
	defer res.Body.Close()

	bodyRaw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("opentopodata.GetElevationPoints failed reading body: %w", err)
	}

	var resp elevationResponse
	if err = json.Unmarshal(bodyRaw, &resp); err != nil {
		return nil, fmt.Errorf("opentopodata.GetElevationPoints failed encoding body: %w", err)
	}

	return ToElevations(resp, box), nil
}

func ToElevations(resp elevationResponse, box osm.Box) models.Elevations {
	var elevations models.Elevations

	for _, r := range resp.Results {
		elevations = append(elevations, models.Elevation{
			Value: r.Elevation,
			Point: models.Point{
				Lat: r.Location.Lat,
				Lon: r.Location.Lng,
			},
			Type: config.TypeLand,
			NeighboringPoints: getNeighboringPoints(models.Point{
				Lat: r.Location.Lat,
				Lon: r.Location.Lng,
			}, box),
		})
	}

	return elevations
}

func getNeighboringPoints(elevation models.Point, box osm.Box) []models.Point {
	var (
		topPoint   models.Point
		botPoint   models.Point
		leftPoint  models.Point
		rightPoint models.Point

		upperLeftPoint   models.Point
		upperRightPoint  models.Point
		bottomLeftPoint  models.Point
		bottomRightPoint models.Point

		point []models.Point
	)

	topPoint.Lat = elevation.Lat - config.Step
	topPoint.Lon = elevation.Lon
	if topPoint.Lat > box.MinLat {
		point = append(point, models.Point{
			Lat: topPoint.Lat,
			Lon: topPoint.Lon,
		})
	}

	botPoint.Lat = elevation.Lat + config.Step
	botPoint.Lon = elevation.Lon
	if botPoint.Lat < box.MaxLat {
		point = append(point, models.Point{
			Lat: botPoint.Lat,
			Lon: botPoint.Lon,
		})
	}

	leftPoint.Lat = elevation.Lat
	leftPoint.Lon = elevation.Lon - config.Step
	if leftPoint.Lon > box.MinLon {
		point = append(point, models.Point{
			Lat: leftPoint.Lat,
			Lon: leftPoint.Lon,
		})
	}

	rightPoint.Lat = elevation.Lat
	rightPoint.Lon = elevation.Lon + config.Step
	if rightPoint.Lon < box.MaxLon {
		point = append(point, models.Point{
			Lat: rightPoint.Lat,
			Lon: rightPoint.Lon,
		})
	}

	upperLeftPoint.Lat = elevation.Lat - config.Step
	upperLeftPoint.Lon = elevation.Lon - config.Step
	if upperLeftPoint.Lat > box.MinLat && upperLeftPoint.Lon > box.MinLon {
		point = append(point, models.Point{
			Lat: upperLeftPoint.Lat,
			Lon: upperLeftPoint.Lon,
		})
	}

	upperRightPoint.Lat = elevation.Lat - config.Step
	upperRightPoint.Lon = elevation.Lon + config.Step
	if upperRightPoint.Lat > box.MinLat && upperRightPoint.Lon < box.MaxLon {
		point = append(point, models.Point{
			Lat: upperRightPoint.Lat,
			Lon: upperRightPoint.Lon,
		})
	}

	bottomLeftPoint.Lat = elevation.Lat + config.Step
	bottomLeftPoint.Lon = elevation.Lon - config.Step
	if bottomLeftPoint.Lat < box.MaxLat && bottomLeftPoint.Lon > box.MinLon {
		point = append(point, models.Point{
			Lat: bottomLeftPoint.Lat,
			Lon: bottomLeftPoint.Lon,
		})
	}

	bottomRightPoint.Lat = elevation.Lat + config.Step
	bottomRightPoint.Lon = elevation.Lon + config.Step
	if bottomLeftPoint.Lat < box.MaxLat && bottomLeftPoint.Lon < box.MaxLon {
		point = append(point, models.Point{
			Lat: bottomLeftPoint.Lat,
			Lon: bottomLeftPoint.Lon,
		})
	}

	return point
}
