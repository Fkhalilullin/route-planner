package route

import (
	"encoding/json"
	"github.com/Fkhalilullin/route-planner/internal/config"
	"github.com/Fkhalilullin/route-planner/internal/models"
	"github.com/Fkhalilullin/route-planner/internal/pather"
	"github.com/Fkhalilullin/route-planner/internal/services/osm"
	"log"
	"math"
	"net/http"
)

func GetPoints(w http.ResponseWriter, r *http.Request) {
	var (
		resp Response
		err  error
	)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&resp)

	//elevationService := topodata.NewService()
	osmService := osm.NewService()

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

	box := osm.Box{
		MinLon: topLeftPoint.Lon, MinLat: topLeftPoint.Lat,
		MaxLon: botRightPoint.Lon, MaxLat: botRightPoint.Lat,
	}

	log.Printf("topLeftPoint = %v\nbotRightPoint = %v\ntopRightPoint = %v\nbotLeftPoint = %v\n",
		topLeftPoint, botRightPoint, topRightPoint, botLeftPoint)

	for lat := topLeftPoint.Lat; lat <= botLeftPoint.Lat; lat += config.Step {
		var elevations []*pather.Coordinate
		for lon := topLeftPoint.Lon; lon <= topRightPoint.Lon; lon += config.Step {
			elevations = append(elevations, &pather.Coordinate{
				Value: 0,
				Point: models.Point{
					Lat: lat,
					Lon: lon,
				},
				Type:              config.TypeLand,
				NeighboringPoints: nil,
			})
		}
		pather.Mesh = append(pather.Mesh, elevations)
	}

	for i, e := range pather.Mesh {
		for j := range e {
			pather.Mesh[i][j].Y = i
			pather.Mesh[i][j].X = j
		}
	}

	pather.Mesh, err = osmService.GetTypePoints(pather.Mesh, box)
	if err != nil {
		log.Printf("[GET/Points] can't get type route: %w", err)
		return
	}
	pather.Mesh = setNeighboringPoints(pather.Mesh)

	beginX, beginY := getForeignPoint(models.Point{
		Lat: resp.BeginPoint.Lat,
		Lon: resp.BeginPoint.Lon,
	})
	endX, endY := getForeignPoint(models.Point{
		Lat: resp.EndPoint.Lat,
		Lon: resp.EndPoint.Lon,
	})
	path, _, _ := pather.Path(pather.Mesh[beginX][beginY], pather.Mesh[endX][endY])

	var result []models.Point
	for _, p := range path {
		converter := p.(*pather.Coordinate)
		result = append(result, models.Point{
			Lat: converter.Point.Lat,
			Lon: converter.Point.Lon,
		})
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("[GET/Points] can't encode to json: %w", err)
		return
	}
}

func setNeighboringPoints(elevations pather.Coordinates) pather.Coordinates {
	for i, e := range elevations {
		for j := range e {
			var bufElevations []*pather.Coordinate
			if i-1 >= 0 && j-1 >= 0 {
				bufElevations = append(bufElevations, &pather.Coordinate{
					Value:             elevations[i-1][j-1].Value,
					Point:             elevations[i-1][j-1].Point,
					Type:              elevations[i-1][j-1].Type,
					X:                 elevations[i-1][j-1].X,
					Y:                 elevations[i-1][j-1].Y,
					NeighboringPoints: nil,
				})
			}
			if i+1 < len(elevations) && j+1 < len(e) {
				bufElevations = append(bufElevations, &pather.Coordinate{
					Value:             elevations[i+1][j+1].Value,
					Point:             elevations[i+1][j+1].Point,
					Type:              elevations[i+1][j+1].Type,
					X:                 elevations[i+1][j+1].X,
					Y:                 elevations[i+1][j+1].Y,
					NeighboringPoints: nil,
				})
			}
			if i-1 >= 0 && j+1 < len(e) {
				bufElevations = append(bufElevations, &pather.Coordinate{
					Value:             elevations[i-1][j+1].Value,
					Point:             elevations[i-1][j+1].Point,
					Type:              elevations[i-1][j+1].Type,
					X:                 elevations[i-1][j+1].X,
					Y:                 elevations[i-1][j+1].Y,
					NeighboringPoints: nil,
				})
			}
			if i+1 < len(elevations) && j-1 >= 0 {
				bufElevations = append(bufElevations, &pather.Coordinate{
					Value:             elevations[i+1][j-1].Value,
					Point:             elevations[i+1][j-1].Point,
					Type:              elevations[i+1][j-1].Type,
					X:                 elevations[i+1][j-1].X,
					Y:                 elevations[i+1][j-1].Y,
					NeighboringPoints: nil,
				})
			}
			if i+1 < len(elevations) {
				bufElevations = append(bufElevations, &pather.Coordinate{
					Value:             elevations[i+1][j].Value,
					Point:             elevations[i+1][j].Point,
					Type:              elevations[i+1][j].Type,
					X:                 elevations[i+1][j].X,
					Y:                 elevations[i+1][j].Y,
					NeighboringPoints: nil,
				})
			}
			if i-1 >= 0 {
				bufElevations = append(bufElevations, &pather.Coordinate{
					Value:             elevations[i-1][j].Value,
					Point:             elevations[i-1][j].Point,
					Type:              elevations[i-1][j].Type,
					X:                 elevations[i-1][j].X,
					Y:                 elevations[i-1][j].Y,
					NeighboringPoints: nil,
				})
			}
			if j+1 < len(e) {
				bufElevations = append(bufElevations, &pather.Coordinate{
					Value:             elevations[i][j+1].Value,
					Point:             elevations[i][j+1].Point,
					Type:              elevations[i][j+1].Type,
					X:                 elevations[i][j+1].X,
					Y:                 elevations[i][j+1].Y,
					NeighboringPoints: nil,
				})
			}
			if j-1 >= 0 {
				bufElevations = append(bufElevations, &pather.Coordinate{
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

func getForeignPoint(point models.Point) (int, int) {
	var (
		minDistance = math.MaxFloat64
		x           int
		y           int
	)

	for i, e := range pather.Mesh {
		for j, ee := range e {
			distance := math.Sqrt(
				(point.Lat-ee.Point.Lat)*(point.Lat-ee.Point.Lat) +
					(point.Lon-ee.Point.Lon)*(point.Lon-ee.Point.Lon),
			)
			if distance < minDistance {
				minDistance = distance
				x = i
				y = j
			}
		}
	}

	return x, y
}
