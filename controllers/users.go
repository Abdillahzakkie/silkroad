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
	var userService models.UserService
	var user models.User
	
	// parse url-encoded form data
	err := helpers.ParseForm(r, &user); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, helpers.ErrorString(err))
		return 
	}

	// save user to database
	user, err = userService.CreateNewUser(user)
	switch err {
		case nil:
			break
		case models.ErrorUserAlreadyExists:
			helpers.RespondWithError(w, http.StatusBadRequest, helpers.ErrorString(err))
			return
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, helpers.ErrorString(err))
			return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GET "/users/all"
// GetAllUsers queries and returns all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var userService models.UserService
	users, err := userService.GetAllUsers(); if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, helpers.ErrorString(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GET "/users?id=<id>"
// GetUserById gets user by ID
func GetUserById(w http.ResponseWriter, r *http.Request) {
	var userService models.UserService

	vars := mux.Vars(r)
	id, err := (strconv.Atoi(vars["id"])); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := userService.GetUserById(uint(id))
	switch err {
		case nil:
			break
		case models.ErrorUserNotFound:
			helpers.RespondWithError(w, http.StatusNotFound, helpers.ErrorString(err))
			return
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, helpers.ErrorString(err))
			return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// DELETE "/users/:id"
// DeleteUser deletes user by ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var userService models.UserService
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["user_id"]); if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "invalid user id")
		return
	}

	// checks if user exists
	user := models.User{
		ID: uint(id),
	}

	err = userService.IsExistingUser(user)
	switch err {
		case nil:
			break
		case models.ErrorUserNotFound:
			helpers.RespondWithError(w, http.StatusNotFound, helpers.ErrorString(err))
			return
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, helpers.ErrorString(err))
			return
	}

	err = userService.DeleteUser(user.ID); if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("%v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}