package topodata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Fkhalilullin/route-planner/internal/models"
)

type ElevationProvider interface {
	GetElevationPoints(coordinateList string) (models.Elevations, error)
}

type service struct {
	elevationProvide ElevationProvider
}

func NewService() *service {
	return &service{}
}

const (
	endpoint                  = "https://api.opentopodata.org/v1"
	pathFormatElevationPoints = "/srtm90m?locations=%s"
)

func (s *service) GetElevationPoints(coordinateList string) (models.Elevations, error) {
	path := fmt.Sprintf(pathFormatElevationPoints, coordinateList)
	u := endpoint + path

	res, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("opentopodata.GetElevationPoints failed http GET: %s", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	bodyRaw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("opentopodata.GetElevationPoints failed reading body: %s", err)
	}

	var resp elevationResponse
	if err = json.Unmarshal(bodyRaw, &resp); err != nil {
		return nil, fmt.Errorf("opentopodata.GetElevationPoints failed encoding body: %s", err)
	}

	return ToElevations(resp), nil
}

func ToElevations(resp elevationResponse) models.Elevations {
	var elevations models.Elevations

	for _, r := range resp.Results {
		elevations = append(elevations, models.Elevation{
			Value: r.Elevation,
			Point: models.Point{
				Lat: r.Location.Lat,
				Lon: r.Location.Lng,
			},
		})

	}

	return elevations
}
