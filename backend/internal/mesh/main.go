package mesh

import (
	"encoding/json"
	"log"
	"math"
	"net/http"

	"github.com/Fkhalilullin/route-planner/internal/config"
	"github.com/Fkhalilullin/route-planner/internal/helpers"
	"github.com/Fkhalilullin/route-planner/internal/models"
	"github.com/Fkhalilullin/route-planner/internal/pather"

	"github.com/Fkhalilullin/route-planner/internal/services/openelevation"
)

func GetMesh(w http.ResponseWriter, r *http.Request) {
	var (
		req Request
		err error
	)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&req)

	elevationService := openelevation.NewService()

	topLeftPoint := models.Point{
		Lat: req.TopLeftPoint.Lat,
		Lon: req.TopLeftPoint.Lon,
	}

	botRightPoint := models.Point{
		Lat: req.BotRightPoint.Lat,
		Lon: req.BotRightPoint.Lon,
	}

	topRightPoint := models.Point{
		Lat: topLeftPoint.Lat,
		Lon: botRightPoint.Lon,
	}

	botLeftPoint := models.Point{
		Lat: botRightPoint.Lat,
		Lon: topLeftPoint.Lon,
	}

	pather.Mesh = pather.Coordinates{}
	for lat := topLeftPoint.Lat; lat <= botLeftPoint.Lat + config.Step; lat += config.Step {
		var elevations []*pather.Coordinate
		for lon := topLeftPoint.Lon; lon <= topRightPoint.Lon + config.Step; lon += config.Step {
			elevations = append(elevations, &pather.Coordinate{
				Value: 0,
				Point: models.Point{
					Lat: helpers.RoundFloat(lat, 6),
					Lon: helpers.RoundFloat(lon, 6),
				},
				Type: config.TypeLand,
			})
		}
		pather.Mesh = append(pather.Mesh, elevations)
	}

	log.Println("Get elevation...")
	pather.Mesh, err = elevationService.GetElevationPoints(pather.Mesh)
	if err != nil {
		log.Printf("[GET/Points] can't get elevaion: %w", err)
		return
	}

	var resp = Response{
		Points: []Points{},
	}
	for _, m := range pather.Mesh {
		for _, mm := range m {
			resp.Points = append(resp.Points, Points{
				Lat:       mm.Point.Lat,
				Lon:       mm.Point.Lon,
				Elevation: mm.Value,
			})
		}
	}

	resp.MinElevation = getMinElevation(resp)
	resp.MaxElevation = getMaxElevation(resp)
	resp.RowCount = len(pather.Mesh)
	resp.ColumnCount = len(pather.Mesh[0])

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("[GET/Points] can't encode to json: %w", err)
		return
	}
}

func getMinElevation(response Response) float64 {
	minElevation := math.MaxFloat64
	for _, p := range response.Points {
		if p.Elevation < minElevation {
			minElevation = p.Elevation
		}
	}
	return minElevation
}

func getMaxElevation(response Response) float64 {
	maxElevation := float64(-100000)
	for _, p := range response.Points {
		if p.Elevation > maxElevation {
			maxElevation = p.Elevation
		}
	}
	return maxElevation
}
