package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abdillahzakkie/silkroad/helpers"
	"github.com/abdillahzakkie/silkroad/models"
	"github.com/gorilla/mux"
)

// POST "/categories/new"
// CreateNewCategory creates new category
func CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var categoryService models.CategoryService
	var category models.Category

	// parse url-encoded form data
	err := helpers.ParseForm(r, &category); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, helpers.ErrorString(err))
		return
	}

	// insert record into DB
	response, err := categoryService.CreateNewCategory(category)
	switch err {
		case nil:
			break
		case models.ErrorCategoryAlreadyExists:
			helpers.RespondWithError(w, http.StatusBadRequest, helpers.ErrorString(err))
			return
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, helpers.ErrorString(err))
			return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GET "/categories"
// GetAllCategories queries and returns all categories
func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	var categoryService models.CategoryService

	categories, err := categoryService.GetAllCategories(); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, helpers.ErrorString(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

// GET "/categories/:category_id"
// GetCategoryById gets category by ID
func GetCategoryById(w http.ResponseWriter, r *http.Request) {
	var categoryService models.CategoryService
	
	vars := mux.Vars(r)
	id, err  := strconv.Atoi(vars["category_id"]); if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "invalid category id received")
		return
	}

	status, err := categoryService.IsExistingCategory(uint(id))
	switch err {
		case nil:
			break
		case models.ErrorCategoryNotFound:
			helpers.RespondWithError(w, http.StatusNotFound, helpers.ErrorString(err))
			return
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, helpers.ErrorString(err))
			return
	}

	if !status {
		helpers.RespondWithError(w, http.StatusNotFound, helpers.ErrorString(models.ErrorCategoryNotFound))
		return
	}

	category, err := categoryService.GetCategoryById(uint(id)); if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, helpers.ErrorString(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}