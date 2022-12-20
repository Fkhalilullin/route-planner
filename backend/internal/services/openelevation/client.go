package openelevation

import (
	"encoding/json"
	"fmt"
	"github.com/Fkhalilullin/route-planner/internal/config"
	"github.com/Fkhalilullin/route-planner/internal/pather"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
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

type Request struct {
	request ElevationRequest
	id      int
}

type Response struct {
	response ElevationResponse
	id       int
}

//const endpoint = "https://api.open-elevation.com/api/v1/lookup"
const endpoint = "http://0.0.0.0:80/api/v1/lookup"

func (s *service) GetElevationPoints(coordinates pather.Coordinates) (pather.Coordinates, error) {

	requests := getRequests(coordinates)
	requestChan := make(chan Request, 1)
	responseChan := make(chan Response, 1)
	go func() {
		defer close(requestChan)
		for _, r := range requests {
			requestChan <- r
		}
	}()

	log.Printf("Total requests: %d", len(requests))
	var wg sync.WaitGroup
	for i := 0; i < config.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range requestChan {
				response, _ := work(r)
				responseChan <- response
			}
		}()
	}

	responses := []Response{}
	go func() {
		defer close(responseChan)
		for r := range responseChan {
			responses = append(responses, r)
		}
	}()

	wg.Wait()

	sort.Slice(responses, func(i, j int) bool {
		return responses[i].id < responses[j].id
	})

	for i, c := range coordinates {
		for j, _ := range c {
			coordinates[i][j].Value = responses[i].response.Results[j].Elevation
		}
	}

	return coordinates, nil
}

func getRequests(coordinates pather.Coordinates) []Request {
	requests := []Request{}
	for i, c := range coordinates {
		req := ElevationRequest{
			Locations: []Locations{},
		}
		for _, cc := range c {
			req.Locations = append(req.Locations, Locations{
				Latitude:  cc.Point.Lat,
				Longitude: cc.Point.Lon,
			})
		}
		requests = append(requests, Request{
			request: req,
			id:      i,
		})
	}

	return requests
}

func work(request Request) (Response, error) {
	log.Printf("Starting get %d id response", request.id)
	reqByte, err := json.Marshal(request.request)
	reader := strings.NewReader(string(reqByte))

	res, err := http.Post(endpoint, "application/json", reader)
	if err != nil {
		return Response{}, fmt.Errorf("opentopodata.GetElevationPoints failed http GET: %w", err)
	}
	defer res.Body.Close()

	bodyRaw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Response{}, fmt.Errorf("opentopodata.GetElevationPoints failed reading body: %w", err)
	}

	var resp ElevationResponse
	if err = json.Unmarshal(bodyRaw, &resp); err != nil {
		return Response{}, fmt.Errorf("opentopodata.GetElevationPoints failed encoding body: %w", err)
	}

	return Response{
		response: resp,
		id:       request.id,
	}, nil

}
