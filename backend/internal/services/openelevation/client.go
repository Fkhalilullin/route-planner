package openelevation

import (
	"encoding/json"
	"fmt"
	"github.com/Fkhalilullin/route-planner/internal/pather"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
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

const endpoint = "https://api.open-elevation.com/api/v1/lookup"

func (s *service) GetElevationPoints(coordinates pather.Coordinates) (pather.Coordinates, error) {

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

	responses := []ElevationResponse{}
	for _, req := range requests {
		log.Print("responses: ", req)
		reqByte, err := json.Marshal(req)
		reader := strings.NewReader(string(reqByte))

		time.Sleep(time.Second / 2)
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
		log.Print("responses: ", responses)
	}

	for i, c := range coordinates {
		for j, _ := range c {
			coordinates[i][j].Value = responses[i].Results[j].Elevation
		}
	}

	return coordinates, nil
}
