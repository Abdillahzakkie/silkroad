package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abdillahzakkie/silkroad/helpers"
	"github.com/abdillahzakkie/silkroad/models"
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
		ID: product.SellerID,
	}

	category := models.Category{
		CategoryID: product.CategoryID,
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
	_isExistingProduct, _ := product.GetProduct(); if _isExistingProduct.ProductID != 0 {
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