package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/abdillahzakkie/silkroad/models"
	"github.com/gorilla/schema"
)

var (
	us *models.UserService
	cs *models.CategoryService
)

func init() {
	psqlInfo, err := models.GetPsqlInfo("silkroad"); if err != nil {
		log.Fatal(err)
	}

	us, err = models.NewUserService(psqlInfo); if err != nil {
		log.Fatal(models.ErrDatabaseConnectionFailed.Error())
	}

	cs, err = models.NewCategoryService(psqlInfo); if err != nil {
		log.Fatal(models.ErrDatabaseConnectionFailed.Error())
	}

	// clear all tables
	us.DestructiveReset()
	cs.DestructiveReset()
}

func ParseForm(r *http.Request, dst interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return errors.New("error while parsing form")
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(dst, r.PostForm)
	if err != nil {
		return errors.New("bad data received")
	}
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(map[string]string{
		"code": fmt.Sprintf("%d", code),
		"error": message,
	})
	
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
}