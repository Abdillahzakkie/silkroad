package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abdillahzakkie/silkroad/models"
	"github.com/gorilla/mux"
)

type CategorySignUpForm struct {
	Name    string 		`schema:"name"`
}

// POST "/categories/new"
// CreateNewCategory creates new category
func CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var form CategorySignUpForm
	// parse url-encoded form data
	if err := ParseForm(r, &form); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	category := models.Category {
		Name: 	form.Name,
	}

	// insert record into DB
	if err := cs.CreateNewCategory(&category); err != nil {
		switch err {
			case models.ErrAlreadyExists:
				RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			default:
				RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
		}
	}
	

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// GET "/categories"
// GetAllCategories queries and returns all categories
func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := cs.GetAllCategories(); if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

// GET "/categories/:category_id"
// GetCategoryById gets category by ID
func GetCategoryById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err  := strconv.Atoi(vars["category_id"]); if err != nil {
		RespondWithError(w, http.StatusNotFound, "invalid category id received")
		return
	}
	category, err := cs.GetCategoryById(uint(id)); if err != nil {
		switch err {
			case models.ErrNotFound:
				RespondWithError(w, http.StatusNotFound, err.Error())
				return
			default:
				RespondWithError(w, http.StatusInternalServerError, models.ErrInternalServerError.Error())
				return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}