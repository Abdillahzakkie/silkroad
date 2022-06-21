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
	"github.com/gorilla/mux"
)

func notFoundRoute(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
}


func main() {
	router := mux.NewRouter()
	// GET "*" handles all unknown routes
	router.NotFoundHandler = http.HandlerFunc(notFoundRoute)

	router.HandleFunc("/users/signup", controllers.CreateNewUser).Methods(http.MethodPost)
	router.HandleFunc("/users/login", controllers.LoginUser).Methods(http.MethodPost)
	router.HandleFunc("/users/all", controllers.GetAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/users", controllers.GetUserById).Queries("id", "{id}").Methods(http.MethodGet)
	router.HandleFunc("/users/{user_id}", controllers.DeleteUserById).Methods(http.MethodDelete)


	router.HandleFunc("/categories/new", controllers.CreateNewCategory).Methods(http.MethodPost)
	router.HandleFunc("/categories", controllers.GetAllCategories).Methods(http.MethodGet)
	router.HandleFunc("/categories/{category_id}", controllers.GetCategoryById).Methods(http.MethodGet)

	// // POST "/products/new"
	// // Create new category
	// router.HandleFunc("/products/new", controllers.CreateNewProduct).Methods(http.MethodPost)

	// // GET "/products"
	// // Get all products
	// router.HandleFunc("/products/all", controllers.GetAllProducts).Methods(http.MethodGet)

	// // GET "/products/:product_id"
	// // Create new category
	// router.HandleFunc("/products/{product_id}", controllers.GetProductById).Methods(http.MethodGet)

	// // GET "/products?seller_id=<seller_id>"
	// // Create new category
	// router.HandleFunc("/products", controllers.GetProductsBySellerId).Queries("seller_id", "{seller_id}").Methods(http.MethodGet)
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
        Handler: r,
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