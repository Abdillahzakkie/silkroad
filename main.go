package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/abdillahzakkie/silkroad/controllers"
	"github.com/abdillahzakkie/silkroad/database"
	"github.com/gorilla/mux"
)


func main() {
	router := mux.NewRouter()
	sqlDB, err := database.DB.DB(); if err != nil {
		log.Fatalln(err)
	}
	// defer closing DB connection
	defer sqlDB.Close()
	
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(5)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(10)

	// GET "*" handles all unknown routes
	router.NotFoundHandler = http.HandlerFunc(notFoundRoute)

	// POST "/users/signup"
	// Creates new user
	router.HandleFunc("/users/signup", controllers.CreateNewUser).Methods(http.MethodPost)

	// GET "/users"
	// Get all users
	router.HandleFunc("/users/all", controllers.GetAllUsers).Methods(http.MethodGet)

	// GET "/users/:id"
	// Get user by ID
	router.HandleFunc("/users", controllers.GetUserById).Queries("id", "{id}").Methods(http.MethodGet)

	// DELETE "/users/:user_id"
	// Delete user by ID
	router.HandleFunc("/users/{user_id}", controllers.DeleteUser).Methods(http.MethodDelete)

	// POST "/categories/new"
	// Create new category
	router.HandleFunc("/categories/new", controllers.CreateNewCategory).Methods(http.MethodPost)

	// GET "/categories"
	// Get all categories
	router.HandleFunc("/categories", controllers.GetAllCategories).Methods(http.MethodGet)

	// GET "/categories/:category_id"
	// Get category by Category ID
	router.HandleFunc("/categories/{category_id}", controllers.GetCategoryById).Methods(http.MethodGet)

	// POST "/products/new"
	// Create new category
	router.HandleFunc("/products/new", controllers.CreateNewProduct).Methods(http.MethodPost)

	startServer(router)
}


func startServer(r *mux.Router) {
	var wait time.Duration
    flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
    flag.Parse()
	port := "8080"

    srv := &http.Server{
        Addr:        fmt.Sprintf( "localhost:%v", port),
        // Good practice to set timeouts to avoid Slowloris attacks.
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler: r, // Pass our instance of gorilla/mux in.
    }

    go func() {
		fmt.Println("Server listening on port:", port)
        log.Fatalln(srv.ListenAndServe())
    }()

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    <-c

    // Create a deadline to wait for.
    ctx, cancel := context.WithTimeout(context.Background(), wait)
    defer cancel()
    srv.Shutdown(ctx)
    log.Println("shutting down")
    os.Exit(0)
}

// func main() {
// 	var user controllers.User
// 	var product controllers.Product
// 	var category controllers.Category

// 	port := "8080"
// 	router := mux.NewRouter()
// 	sqlDB, err := database.DB.DB(); if err != nil {
// 		log.Fatalln(err)
// 	}
// 	// defer closing DB connection
// 	defer sqlDB.Close()
	
// 	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
// 	sqlDB.SetMaxIdleConns(5)

// 	// SetMaxOpenConns sets the maximum number of open connections to the database.
// 	sqlDB.SetMaxOpenConns(10)
// 	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
// 	// sqlDB.SetConnMaxLifetime(time.Second)

// 	// GET "*" handles all unknown routes
// 	router.NotFoundHandler = http.HandlerFunc(notFoundRoute)

// 	// POST "/users/signup"
// 	// Creates new user
// 	router.HandleFunc("/users/signup", user.CreateNewUser).Methods(http.MethodPost)
	
// 	// POST "/products/{seller_id}/new"
// 	// Creates new product
// 	router.HandleFunc("/products/{seller_id}/new", product.CreateNewProduct).Methods(http.MethodPost)

// 	// POST "/categories/new"
// 	// Creates new categories
// 	router.HandleFunc("/categories/new", category.CreateNewCategory).Methods(http.MethodPost)

// 	// GET "/categories/"
// 	// Get all categories
// 	router.HandleFunc("/categories", category.GetAllCategories).Methods(http.MethodGet)

// 	// GET "/categories/{category_id}"
// 	// Get category by ID
// 	router.HandleFunc("/categories/{category_id}", category.GetCategoryById).Methods(http.MethodGet)

// 	log.Println("Server listening on port: ", port)
// 	log.Fatalln(http.ListenAndServe(fmt.Sprintf("localhost:%s", port), router))
// }

func notFoundRoute(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
}