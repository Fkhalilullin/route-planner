package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func GetRoute(w http.ResponseWriter, r *http.Request) {
	beginPoint, err := getBeginPoint(r)
	if err != nil {
		log.Println("[GET/Route] can't parse begin point:", err)
		return
	}

	endPoint, err := getEndPoint(r)
	if err != nil {
		log.Println("[GET/Route] can't parse end point:", err)
		return
	}
	route := getRoute(beginPoint, endPoint)
	log.Println("[GET/Route] output route", route)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(route)
}

func getRoute(beginPoint Point, endPoint Point) Route {
	var route Route
	route = append(route, beginPoint, endPoint)

	return route
}

func getBeginPoint(r *http.Request) (Point, error) {
	latBegin, err := strconv.ParseFloat(r.URL.Query().Get(FromLat), 64)
	if err != nil {
		return Point{}, fmt.Errorf("can't get from_lat value: %w", err)
	}
	lonBegin, err := strconv.ParseFloat(r.URL.Query().Get(FromLon), 64)
	if err != nil {
		return Point{}, fmt.Errorf("can't get from_lon value: %w", err)
	}

	return Point{
		Lat: latBegin,
		Lon: lonBegin,
	}, nil
}

func getEndPoint(r *http.Request) (Point, error) {
	latEnd, err := strconv.ParseFloat(r.URL.Query().Get(ToLat), 64)
	if err != nil {
		return Point{}, fmt.Errorf("can't get to_lat value: %w", err)
	}
	lonEnd, err := strconv.ParseFloat(r.URL.Query().Get(ToLon), 64)
	if err != nil {
		return Point{}, fmt.Errorf("can't get to_lon value: %w", err)
	}

	return Point{
		Lat: latEnd,
		Lon: lonEnd,
	}, nil
}
