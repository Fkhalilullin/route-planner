package points

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Fkhalilullin/route-planner/internal/config"
	"github.com/Fkhalilullin/route-planner/internal/models"
	"github.com/Fkhalilullin/route-planner/internal/services/osm"
	"github.com/Fkhalilullin/route-planner/internal/services/topodata"
)

func GetPoints(w http.ResponseWriter, r *http.Request) {
	var (
		elevations     models.Elevations
		coordinateList string
	)

	elevationService := topodata.NewService()
	osmService := osm.NewService()

	topLeftPoint, err := getTopLeftPoint(r)
	if err != nil {
		log.Printf("[GET/Points] can't parse begin point: %w\n", err)
		return
	}

	botRightPoint, err := getBotRightPoint(r)
	if err != nil {
		log.Printf("[GET/Route] can't parse end point: %w\n", err)
		return
	}

	topRightPoint := models.Point{
		Lat: topLeftPoint.Lat,
		Lon: botRightPoint.Lon,
	}

	botLeftPoint := models.Point{
		Lat: botRightPoint.Lat,
		Lon: topLeftPoint.Lon,
	}

	log.Printf("topLeftPoint = %v\nbotRightPoint = %v\ntopRightPoint = %v\nbotLeftPoint = %v\n",
		topLeftPoint, botRightPoint, topRightPoint, botLeftPoint)

	var limiter int64
	for lat := topLeftPoint.Lat; lat <= botLeftPoint.Lat; lat += config.Step {
		for lon := topLeftPoint.Lon; lon <= topRightPoint.Lon; lon += config.Step {
			coordinateList += fmt.Sprintf("%f,%f|", lat, lon)

			// TODO развернуть на докере
			limiter++
			if limiter == 100 {
				temp, err := elevationService.GetElevationPoints(coordinateList)
				if err != nil {
					log.Println("[GET/Points] can't get elevation points: ", err)
				}
				elevations = append(elevations, temp...)

				limiter = 0
				coordinateList = ""
			}
		}
	}
	log.Printf("[GET/Points] count elevations points: %d", len(elevations))

	osmService.GetTypePoints(osm.Box{
		MinLon: topLeftPoint.Lon, MinLat: topLeftPoint.Lat,
		MaxLon: botRightPoint.Lon, MaxLat: botRightPoint.Lat,
	})

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(elevations)
	if err != nil {
		log.Println("[GET/Points] can't encode to json")
	}
}

func getTopLeftPoint(r *http.Request) (models.Point, error) {
	latTopLeftPoint, err := strconv.ParseFloat(r.URL.Query().Get(LatTopLeftPoint), 64)
	if err != nil {
		return models.Point{}, fmt.Errorf("can't get top_left_lat value: %w", err)
	}
	lonTopLeftPoint, err := strconv.ParseFloat(r.URL.Query().Get(LonTopLeftPoint), 64)
	if err != nil {
		return models.Point{}, fmt.Errorf("can't get top_left_lon value: %w", err)
	}

	return models.Point{
		Lat: latTopLeftPoint,
		Lon: lonTopLeftPoint,
	}, nil
}

func getBotRightPoint(r *http.Request) (models.Point, error) {
	latBotRightPoint, err := strconv.ParseFloat(r.URL.Query().Get(LatBotRightPoint), 64)
	if err != nil {
		return models.Point{}, fmt.Errorf("can't get bot_right_lat value: %w", err)
	}
	lonBotRightPoint, err := strconv.ParseFloat(r.URL.Query().Get(LonBotRightPoint), 64)
	if err != nil {
		return models.Point{}, fmt.Errorf("can't get bot_right_lon value: %w", err)
	}

	return models.Point{
		Lat: latBotRightPoint,
		Lon: lonBotRightPoint,
	}, nil
}
