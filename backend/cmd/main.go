package main

import (
	"github.com/Fkhalilullin/route-planner/internal/mesh"
	"github.com/Fkhalilullin/route-planner/internal/route"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/route", route.GetPoints).Methods(http.MethodPost)
	r.HandleFunc("/mesh", mesh.GetMesh).Methods(http.MethodPost)

	headersOk := handlers.AllowedHeaders([]string{"*"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{
		http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodOptions})

	log.Println("Server start")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
