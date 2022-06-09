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

// POST "/users/new"
// CreateNewUser creates new user
func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// parse url-encoded form data
	err := helpers.ParseForm(r, &user); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return 
	}

	// save user to database
	err = user.CreateNewUser(); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "user has already existed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GET "/users/all"
// GetAllUsers queries and returns all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var user models.User
	users, err := user.GetAllUsers(); if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GET "/users?id=<id>"
// GetUserById gets user by ID
func GetUserById(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var err error

	vars := mux.Vars(r)
	user.ID, err = strconv.Atoi(vars["id"]); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	err = user.GetUser(); if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "user does not exist")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// DELETE "/users/:id"
// DeleteUser deletes user by ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var err error
	vars := mux.Vars(r)

	user.ID, err = strconv.Atoi(vars["user_id"]); if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "invalid user id")
		return
	}

	// checks if user exists
	err = user.GetUser(); if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "user does not exist")
		return
	}

	err = user.DeleteUser(); if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}