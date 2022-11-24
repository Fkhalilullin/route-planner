package main

import (
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"

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
		Queries(points.LatTopLeftPoint, "{min_lat}",
			points.LonTopLeftPoint, "{min_lon}",
			points.LatBotRightPoint, "{max_lat}",
			points.LonBotRightPoint, "{max_lon}").
		Methods("GET")

	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Println("Server start")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
