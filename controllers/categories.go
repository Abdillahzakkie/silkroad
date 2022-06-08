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

// POST "/categories/new"
// CreateNewCategory creates new category
func CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	// parse url-encoded form data
	err := helpers.ParseForm(r, &category); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	if category.Name == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "name is required")
		return
	}
	// insert record into DB
	err = category.CreateNewCategory(); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	var categories []models.Category

	categories, err := category.GetAllCategories(); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func GetCategoryById(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	var err error
	
	vars := mux.Vars(r)
	category.CategoryID, err  = strconv.Atoi(vars["category_id"]); if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "invalid category id received")
		return
	}

	err = category.GetCategoryById(); if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("%v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}