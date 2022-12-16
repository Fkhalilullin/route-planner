package route

import (
	"encoding/json"
	"github.com/Fkhalilullin/route-planner/internal/config"
	"github.com/Fkhalilullin/route-planner/internal/models"
	"github.com/Fkhalilullin/route-planner/internal/pather"
	"github.com/Fkhalilullin/route-planner/internal/services/openelevation"
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

	elevationService := openelevation.NewService()
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

	pather.Mesh = pather.Coordinates{}
	for lat := topLeftPoint.Lat; lat <= botLeftPoint.Lat; lat += config.Step {
		var elevations []*pather.Coordinate
		for lon := topLeftPoint.Lon; lon <= topRightPoint.Lon; lon += config.Step {
			elevations = append(elevations, &pather.Coordinate{
				Value: 0,
				Point: models.Point{
					Lat: roundFloat(lat, 6),
					Lon: roundFloat(lon, 6),
				},
				Type: config.TypeLand,
			})
		}
		pather.Mesh = append(pather.Mesh, elevations)
	}

	for i, c := range pather.Mesh {
		for j, _ := range c {
			pather.Mesh[i][j].X = i
			pather.Mesh[i][j].Y = j
		}
	}

	log.Println("Get elevation...")
	pather.Mesh, err = elevationService.GetElevationPoints(pather.Mesh)
	if err != nil {
		log.Printf("[GET/Points] can't get elevaion: %w", err)
		return
	}

	log.Println("Get types...")
	pather.Mesh, err = osmService.GetTypePoints(pather.Mesh, box)
	if err != nil {
		log.Printf("[GET/Points] can't get type route: %w", err)
		return
	}

	beginX, beginY := getForeignPoint(models.Point{
		Lat: resp.BeginPoint.Lat,
		Lon: resp.BeginPoint.Lon,
	})
	endX, endY := getForeignPoint(models.Point{
		Lat: resp.EndPoint.Lat,
		Lon: resp.EndPoint.Lon,
	})

	log.Println("BeginPoint: ", pather.Mesh[beginX][beginY])
	log.Println("EndPoint: ", pather.Mesh[endX][endY])
	path, distance, _ := pather.Path(pather.Mesh[beginX][beginY], pather.Mesh[endX][endY])

	var result []models.Point
	for _, p := range path {
		converter := p.(*pather.Coordinate)
		log.Printf("Type=%s Lat=%f Lon=%f Elev=%f",
			converter.Type, converter.Point.Lat, converter.Point.Lon, converter.Value)
		result = append(result, models.Point{
			Lat: converter.Point.Lat,
			Lon: converter.Point.Lon,
		})
	}
	log.Printf("Total distance: %f", distance)

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("[GET/Points] can't encode to json: %w", err)
		return
	}
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

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
