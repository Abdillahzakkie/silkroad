package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abdillahzakkie/silkroad/helpers"
	"github.com/abdillahzakkie/silkroad/models"
	"github.com/gorilla/mux"
)

type UserSignUpForm struct {
	Wallet   string 	`schema:"wallet,required"`
	Username string 	`schema:"username,required"`
	Email    string 	`schema:"email,required"`
	Password string 	`schema:"password,required"`
}

type UserLoginForm struct {
	Email    string 	`schema:"email"`
	Password string 	`schema:"password,required"`
}

// CreateNewUser creates new user
// POST "/users/new"
func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	var form UserSignUpForm
	// parse url-encoded form data
	if err := helpers.ParseForm(r, &form); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return 
	}

	user := models.User {
		Wallet: 	form.Wallet,
		Username:   form.Username,
		Email: 		form.Email,
		Password:   form.Password,
	}

	// save user to database
	userErr := make(chan error)
	go func() {
		userErr <- us.CreateNewUser(&user)
	}()

	if err, ok := <- userErr; err != nil {
		if !ok {
			helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		switch err {
			case models.ErrInternalServerError:
				helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			default:
				helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
		}
	}

	// if err := us.CreateNewUser(&user); err != nil {
	// 	switch err {
	// 		case models.ErrInternalServerError:
	// 			helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
	// 			return
	// 		default:
	// 			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
	// 			return
	// 	}
	// }


	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// LoginUser logins user with the provided credentials
// POST "/users/login"
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var form UserLoginForm
	if err := helpers.ParseForm(r, &form); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// authenticate user 
	user, err := us.Authenticate(form.Email, form.Password); if err != nil {
		switch err {
			case models.ErrInvalidCredentials:
				helpers.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			default:
				helpers.RespondWithError(w, http.StatusInternalServerError, models.ErrInternalServerError.Error())
				return
		}
	}

	cookies := http.Cookie{
		Name: "email",
		Value: user.Email,
	}
	http.SetCookie(w, &cookies)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// GetAllUsers queries and returns all users
// GET "/users/all"
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := us.GetAllUsers(); if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetUserById gets user by ID
// GET "/users?id={id}"
func GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := (strconv.Atoi(vars["id"])); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := us.GetUserById(uint(id)); if err != nil {
		switch err {
			case models.ErrNotFound:
				helpers.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			default:
				helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// DeleteUserById deletes user by ID
// DELETE "/users/:id"
func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["user_id"]); if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "invalid user id")
		return
	}

	if err := us.DeleteUserById(uint(id)); err != nil {
		switch err {
			case models.ErrNotFound:
				helpers.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			default:
				helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}