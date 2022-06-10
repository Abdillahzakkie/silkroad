package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/abdillahzakkie/silkroad/helpers"
	"github.com/abdillahzakkie/silkroad/models"
	"github.com/gorilla/mux"
)

// POST "/products/{seller_id}/new"
// CreateNewProduct creates new product
func CreateNewProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	// parse url-encoded form data into struct
	err := helpers.ParseForm(r, &product); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user  := models.User{
		ID: product.UserID,
	}

	category := models.Category{
		ID: product.CategoryID,
	}

	// lookup seller by Id
	err = user.GetUser(); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "seller does not exist")
		return
	}

	// checks if category ID exists
	err = category.GetCategoryById(); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "category does not exist")
		return
	}

	// checks if product has already existed
	// discard error if product does not exist
	_isExistingProduct, _ := product.GetProduct(); if _isExistingProduct.ID != 0 {
		helpers.RespondWithError(w, http.StatusBadRequest, "product already exists")
		return
	}

	// insert new product record
	err = product.CreateNewProduct(); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// GET "/products"
// GetAllProducts get all products
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	products, err := product.GetAllProducts(); if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("%v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// GET "/products/{product_id}"
// GetProductById get product by Product ID
func GetProductById(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	var err error
	vars := mux.Vars(r)
	
	product.ID, err = strconv.Atoi(vars["product_id"]); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	err = product.GetProductById(); if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "product does not exist")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// GET "/products/{seller_id}"
// GetProductsBySellerId get product by Seller ID
func GetProductsBySellerId(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	var user models.User
	var err error
	vars := mux.Vars(r)
	
	product.UserID, err = strconv.Atoi(vars["seller_id"]); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid seller ID")
		return
	}

	// checks if seller exists
	err = user.GetUser(); if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "seller does not exist")
		return
	}

	// get all products by seller ID
	products, err := product.GetProductsBySellerId(); if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("%v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}