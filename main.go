package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/abdillahzakkie/silkroad/controllers"
	"github.com/abdillahzakkie/silkroad/database"
	"github.com/gorilla/mux"
)

func main() {
	var product controllers.Product
	var category controllers.Category

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

	
	// POST "/products/{seller_id}/new"
	// Creates new product
	router.HandleFunc("/products/{seller_id}/new", product.CreateNewProduct).Methods(http.MethodPost)

	// POST "/categories/new"
	// Creates new categories
	router.HandleFunc("/categories/new", category.CreateNewCategory).Methods(http.MethodPost)

	// GET "/categories/"
	// Creates new product
	router.HandleFunc("/categories", category.GetAllCategories).Methods(http.MethodGet)

	log.Println("Server listening on port: ", port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf("localhost:%s", port), router))
}

func notFoundRoute(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
}