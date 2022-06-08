package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abdillahzakkie/silkroad/helpers"
	"github.com/abdillahzakkie/silkroad/models"
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

