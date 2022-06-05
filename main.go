package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	port := "8080"
	router := mux.NewRouter()

	// GET "*" handles all unknown routes
	router.NotFoundHandler = http.HandlerFunc(notFoundRoute)

	log.Println("Server listening on port: ", port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf("localhost:%s", port), router))
}

func notFoundRoute(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
}