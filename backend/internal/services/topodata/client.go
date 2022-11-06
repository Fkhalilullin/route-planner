package topodata

import (
	"encoding/json"
	"fmt"
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

func (s *service) GetElevationPoints(coordinateList string) (models.Elevations, error) {
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
