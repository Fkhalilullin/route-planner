package openelevation

import (
	"encoding/json"
	"fmt"
	"github.com/Fkhalilullin/route-planner/internal/config"
	"github.com/Fkhalilullin/route-planner/internal/pather"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type ElevationProvider interface {
	GetElevationPoints(coordinateList string) (pather.Coordinates, error)
}

type service struct {
	elevationProvider ElevationProvider
}

func NewService() *service {
	return &service{}
}

func (s *service) GetElevationPoints(coordinates pather.Coordinates) (pather.Coordinates, error) {
	var endpoint string
	endpoint = "https://api.open-elevation.com/api/v1/lookup"
	if config.UseLocalHost {
		endpoint = "http://localhost:80/api/v1/lookup"
	}

	requests := []ElevationRequest{}
	for _, c := range coordinates {
		req := ElevationRequest{
			Locations: []Locations{},
		}
		for _, cc := range c {
			req.Locations = append(req.Locations, Locations{
				Latitude:  cc.Point.Lat,
				Longitude: cc.Point.Lon,
			})
		}
		requests = append(requests, req)
	}

	log.Printf("Total requests: %d", len(requests))
	responses := []ElevationResponse{}
	for i, req := range requests {
		log.Printf("Starting get %d response", i+1)
		reqByte, err := json.Marshal(req)
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

		var resp ElevationResponse
		if err = json.Unmarshal(bodyRaw, &resp); err != nil {
			return nil, fmt.Errorf("opentopodata.GetElevationPoints failed encoding body: %w", err)
		}
		responses = append(responses, resp)
	}

	for i, c := range coordinates {
		for j, _ := range c {
			coordinates[i][j].Value = responses[i].Results[j].Elevation
		}
	}

	return coordinates, nil
}
