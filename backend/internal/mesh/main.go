package mesh

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Fkhalilullin/route-planner/internal/config"
	"github.com/Fkhalilullin/route-planner/internal/helpers"
	"github.com/Fkhalilullin/route-planner/internal/models"
	"github.com/Fkhalilullin/route-planner/internal/pather"

	"github.com/Fkhalilullin/route-planner/internal/services/openelevation"
)

func GetMesh(w http.ResponseWriter, r *http.Request) {
	var (
		resp Request
		err  error
	)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&resp)

	elevationService := openelevation.NewService()

	topLeftPoint := models.Point{
		Lat: resp.TopLeftPoint.Lat,
		Lon: resp.TopLeftPoint.Lon,
	}

	botRightPoint := models.Point{
		Lat: resp.BotRightPoint.Lat,
		Lon: resp.BotRightPoint.Lon,
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
	for lat := topLeftPoint.Lat; lat < botLeftPoint.Lat; lat += config.Step {
		var elevations []*pather.Coordinate
		for lon := topLeftPoint.Lon; lon < topRightPoint.Lon; lon += config.Step {
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

	var request = Response{
		Points: []Points{},
	}
	for _, m := range pather.Mesh {
		for _, mm := range m {
			request.Points = append(request.Points, Points{
				Lat:       mm.Point.Lat,
				Lon:       mm.Point.Lon,
				Elevation: mm.Value,
			})
		}
	}

	request.RowCount = len(pather.Mesh)
	request.ColumnCount = len(pather.Mesh[0])

	err = json.NewEncoder(w).Encode(request)
	if err != nil {
		log.Printf("[GET/Points] can't encode to json: %w", err)
		return
	}
}
