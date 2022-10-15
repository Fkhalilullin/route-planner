package topodata

import (
	"encoding/json"
	"fmt"
	"github.com/Fkhalilullin/route-planner/internal/route"
	"io/ioutil"
	"net/http"
)

type ElevationProvider interface {
	GetElevationPoints(lat float64, lon float64) (Elevations, error)
}

type service struct {
	elevationProvide ElevationProvider
}

func NewService() *service {
	return &service{}
}

const (
	endpoint                  = "https://api.opentopodata.org/v1"
	pathFormatElevationPoints = "/srtm90m?locations=%f,%f"
)

func (s *service) GetElevationPoints(lat float64, lon float64) (Elevation, error) {
	path := fmt.Sprintf(pathFormatElevationPoints, lat, lon)
	u := endpoint + path

	res, err := http.Get(u)
	if err != nil {
		return Elevation{}, fmt.Errorf("opentopodata.GetElevationPoints failed http GET: %s", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	bodyRaw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Elevation{}, fmt.Errorf("opentopodata.GetElevationPoints failed reading body: %s", err)
	}

	var resp elevationResponse
	if err = json.Unmarshal(bodyRaw, &resp); err != nil {
		return Elevation{}, fmt.Errorf("opentopodata.GetElevationPoints failed encoding body: %s", err)
	}

	return ToElevation(resp), nil
}

func ToElevation(resp elevationResponse) Elevation {
	var (
		value float64
		point route.Point
	)

	for _, r := range resp.Results {
		value = r.Elevation
		point = route.Point{
			Lat: r.Location.Lat,
			Lon: r.Location.Lng,
		}
	}

	return Elevation{
		Value: value,
		Point: point,
	}
}
