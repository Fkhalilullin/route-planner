package main

import (
	"log"
	"net/http"

	"github.com/Fkhalilullin/route-planner/internal/points"
	"github.com/Fkhalilullin/route-planner/internal/route"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/route", route.GetRoute).
		Queries(route.FromLat, "{from_lat}",
			route.FromLon, "{from_lon}",
			route.ToLat, "{to_lat}",
			route.ToLon, "{to_lon}").
		Methods("GET")

	r.HandleFunc("/points", points.GetPoints).
		Queries(points.LatTopLeftPoint, "{top_left_lat}",
			points.LonTopLeftPoint, "{top_left_lon}",
			points.LatBotRightPoint, "{bot_right_lat}",
			points.LonBotRightPoint, "{bot_right_lon}").
		Methods("GET")

	log.Println("Server start")
	log.Fatal(http.ListenAndServe(":8000", r))
}
