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

	log.Println("Server listening on port: ", port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf("localhost:%s", port), router))
}