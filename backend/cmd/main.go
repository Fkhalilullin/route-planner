package main

import (
	"github.com/Fkhalilullin/route-planner/internal/points"
	"log"
	"net/http"

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

	// TODO Переделать Queries
	r.HandleFunc("/points", points.GetPoints).
		Queries(route.FromLat, "{from_lat}",
			route.FromLon, "{from_lon}",
			route.ToLat, "{to_lat}",
			route.ToLon, "{to_lon}").
		Methods("GET")

	log.Println("Server start")
	log.Fatal(http.ListenAndServe(":8000", r))
}
