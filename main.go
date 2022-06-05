package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/abdillahzakkie/silkroad/database"
	"github.com/gorilla/mux"
)

func main() {
	port := "8080"
	router := mux.NewRouter()
	sqlDB, err := database.DB.DB(); if err != nil {
		log.Fatalln(err)
	}
	defer sqlDB.Close()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(5)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(10)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// sqlDB.SetConnMaxLifetime(time.Second)

	// GET "*" handles all unknown routes
	router.NotFoundHandler = http.HandlerFunc(notFoundRoute)

	log.Println("Server listening on port: ", port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf("localhost:%s", port), router))
}

func notFoundRoute(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
}