package points

import (
	"encoding/json"
	"fmt"
	"github.com/Fkhalilullin/route-planner/internal/config"
	"github.com/Fkhalilullin/route-planner/internal/models"
	"github.com/Fkhalilullin/route-planner/internal/route"
	"github.com/Fkhalilullin/route-planner/internal/services/osm"
	"log"
	"net/http"
	"strconv"
)

func GetPoints(w http.ResponseWriter, r *http.Request) {
	//var (
	//	elevations     models.Elevations
	//	//coordinateList string
	//)

	//elevationService := topodata.NewService()
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

	box := osm.Box{
		MinLon: topLeftPoint.Lon, MinLat: topLeftPoint.Lat,
		MaxLon: botRightPoint.Lon, MaxLat: botRightPoint.Lat,
	}

	log.Printf("topLeftPoint = %v\nbotRightPoint = %v\ntopRightPoint = %v\nbotLeftPoint = %v\n",
		topLeftPoint, botRightPoint, topRightPoint, botLeftPoint)

	// TODO Расчет высот
	//var limiter int64
	//for lat := topLeftPoint.Lat; lat <= botLeftPoint.Lat; lat += config.Step {
	//	for lon := topLeftPoint.Lon; lon <= topRightPoint.Lon; lon += config.Step {
	//		coordinateList += fmt.Sprintf("%f,%f|", lat, lon)
	//
	//		// TODO развернуть на докере
	//		limiter++
	//		if limiter == 100 {
	//			temp, err := elevationService.GetElevationPoints(coordinateList, box)
	//			if err != nil {
	//				log.Println("[GET/Points] can't get elevation points: ", err)
	//				return
	//			}
	//			elevations = append(elevations, temp...)
	//
	//			limiter = 0
	//			coordinateList = ""
	//		}
	//	}
	//}
	dLat := int((botLeftPoint.Lat - topLeftPoint.Lat) / config.Step)
	models.Mesh = make(models.Elevations, dLat)

	var counter int
	for lat := topLeftPoint.Lat; lat <= botLeftPoint.Lat; lat += config.Step {
		var lonElevations []*models.Elevation
		if counter == dLat {
			break
		}
		for lon := topLeftPoint.Lon; lon <= topRightPoint.Lon; lon += config.Step {
			lonElevations = append(lonElevations, &models.Elevation{
				Value: 0,
				Point: models.Point{
					Lat: lat,
					Lon: lon,
				},
				Type:              config.TypeLand,
				NeighboringPoints: nil,
			})
		}
		models.Mesh[counter] = lonElevations
		counter++
	}

	for i, e := range models.Mesh {
		for j := range e {
			models.Mesh[i][j].Y = i
			models.Mesh[i][j].X = j
		}
	}
	log.Printf("[GET/Points] count elevations points: %d", len(models.Mesh))

	models.Mesh, err = osmService.GetTypePoints(models.Mesh, box)
	if err != nil {
		log.Printf("[GET/Points] can't get type points: %w", err)
		return
	}

	models.Mesh = setNeighboringPoints(models.Mesh)

	var ro []*models.Elevation
	for i, e := range models.Mesh {
		for j := range e {
			ro = append(ro, &models.Elevation{
				Value:             models.Mesh[i][j].Value,
				Point:             models.Mesh[i][j].Point,
				Type:              models.Mesh[i][j].Type,
				NeighboringPoints: models.Mesh[i][j].NeighboringPoints,
				X:                 models.Mesh[i][j].X,
				Y:                 models.Mesh[i][j].Y,
			})
		}
	}
	//fmt.Println(models.Mesh[0][0], models.Mesh[5][5])
	path, dist, found := route.Path(models.Mesh[0][0], models.Mesh[6][7])
	fmt.Println(path, dist, found)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(path)
	if err != nil {
		log.Printf("[GET/Points] can't encode to json: %w", err)
		return
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

func setNeighboringPoints(elevations models.Elevations) models.Elevations {
	for i, e := range elevations {
		for j := range e {
			var bufElevations []*models.Elevation
			if i-1 >= 0 && j-1 >= 0 {
				bufElevations = append(bufElevations, &models.Elevation{
					Value:             elevations[i-1][j-1].Value,
					Point:             elevations[i-1][j-1].Point,
					Type:              elevations[i-1][j-1].Type,
					X:                 elevations[i-1][j-1].X,
					Y:                 elevations[i-1][j-1].Y,
					NeighboringPoints: nil,
				})
			}
			if i+1 < len(elevations) && j+1 < len(e) {
				bufElevations = append(bufElevations, &models.Elevation{
					Value:             elevations[i+1][j+1].Value,
					Point:             elevations[i+1][j+1].Point,
					Type:              elevations[i+1][j+1].Type,
					X:                 elevations[i+1][j+1].X,
					Y:                 elevations[i+1][j+1].Y,
					NeighboringPoints: nil,
				})
			}
			if i-1 >= 0 && j+1 < len(e) {
				bufElevations = append(bufElevations, &models.Elevation{
					Value:             elevations[i-1][j+1].Value,
					Point:             elevations[i-1][j+1].Point,
					Type:              elevations[i-1][j+1].Type,
					X:                 elevations[i-1][j+1].X,
					Y:                 elevations[i-1][j+1].Y,
					NeighboringPoints: nil,
				})
			}
			if i+1 < len(elevations) && j-1 >= 0 {
				bufElevations = append(bufElevations, &models.Elevation{
					Value:             elevations[i+1][j-1].Value,
					Point:             elevations[i+1][j-1].Point,
					Type:              elevations[i+1][j-1].Type,
					X:                 elevations[i+1][j-1].X,
					Y:                 elevations[i+1][j-1].Y,
					NeighboringPoints: nil,
				})
			}
			if i+1 < len(elevations) {
				bufElevations = append(bufElevations, &models.Elevation{
					Value:             elevations[i+1][j].Value,
					Point:             elevations[i+1][j].Point,
					Type:              elevations[i+1][j].Type,
					X:                 elevations[i+1][j].X,
					Y:                 elevations[i+1][j].Y,
					NeighboringPoints: nil,
				})
			}
			if i-1 >= 0 {
				bufElevations = append(bufElevations, &models.Elevation{
					Value:             elevations[i-1][j].Value,
					Point:             elevations[i-1][j].Point,
					Type:              elevations[i-1][j].Type,
					X:                 elevations[i-1][j].X,
					Y:                 elevations[i-1][j].Y,
					NeighboringPoints: nil,
				})
			}
			if j+1 < len(e) {
				bufElevations = append(bufElevations, &models.Elevation{
					Value:             elevations[i][j+1].Value,
					Point:             elevations[i][j+1].Point,
					Type:              elevations[i][j+1].Type,
					X:                 elevations[i][j+1].X,
					Y:                 elevations[i][j+1].Y,
					NeighboringPoints: nil,
				})
			}
			if j-1 >= 0 {
				bufElevations = append(bufElevations, &models.Elevation{
					Value:             elevations[i][j-1].Value,
					Point:             elevations[i][j-1].Point,
					Type:              elevations[i][j-1].Type,
					X:                 elevations[i][j-1].X,
					Y:                 elevations[i][j-1].Y,
					NeighboringPoints: nil,
				})
			}
			elevations[i][j].NeighboringPoints = append(elevations[i][j].NeighboringPoints, bufElevations...)
		}
	}

	return elevations
}
