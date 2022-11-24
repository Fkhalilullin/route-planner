package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Fkhalilullin/route-planner/internal/models"
)

func GetRoute(w http.ResponseWriter, r *http.Request) {
	beginPoint, err := getBeginPoint(r)
	if err != nil {
		log.Printf("[GET/Route] can't parse begin point: %w\n", err)
		return
	}

	endPoint, err := getEndPoint(r)
	if err != nil {
		log.Printf("[GET/Route] can't parse end point: %w\n", err)
		return
	}
	route := getRoute(beginPoint, endPoint)
	log.Printf("[GET/Route] output route: %v", route)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(route)
}

func getRoute(beginPoint models.Point, endPoint models.Point) models.Route {
	var route models.Route
	route = append(route, beginPoint, endPoint)

	return route
}

func getBeginPoint(r *http.Request) (models.Point, error) {
	latBegin, err := strconv.ParseFloat(r.URL.Query().Get(FromLat), 64)
	if err != nil {
		return models.Point{}, fmt.Errorf("can't get from_lat value: %w", err)
	}
	lonBegin, err := strconv.ParseFloat(r.URL.Query().Get(FromLon), 64)
	if err != nil {
		return models.Point{}, fmt.Errorf("can't get from_lon value: %w", err)
	}

	return models.Point{
		Lat: latBegin,
		Lon: lonBegin,
	}, nil
}

func getEndPoint(r *http.Request) (models.Point, error) {
	latEnd, err := strconv.ParseFloat(r.URL.Query().Get(ToLat), 64)
	if err != nil {
		return models.Point{}, fmt.Errorf("can't get to_lat value: %w", err)
	}
	lonEnd, err := strconv.ParseFloat(r.URL.Query().Get(ToLon), 64)
	if err != nil {
		return models.Point{}, fmt.Errorf("can't get to_lon value: %w", err)
	}

	return models.Point{
		Lat: latEnd,
		Lon: lonEnd,
	}, nil
}
