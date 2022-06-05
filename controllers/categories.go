package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/abdillahzakkie/silkroad/database"
	"github.com/abdillahzakkie/silkroad/helpers"
	"github.com/gorilla/mux"
)

type Category struct {
	CategoryId		uint       		`gorm:"primaryKey" json:"category_id" schema:"category_id"`
	CreatedAt		time.Time  		`json:"-" schema:"-"`
	UpdatedAt 		time.Time  		`json:"-" schema:"-"`
	DeletedAt 		*time.Time 		`gorm:"index;<-:update" json:"-" schema:"-"`
	Name      		string     		`gorm:"<-;not null;uniqueIndex" json:"name" schema:"name"`
}

func init() {
	// migrate database
	database.DB.AutoMigrate(&Category{})
}

// POST "/categories/new"
// CreateNewCategory creates new category
func (c *Category) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	category := Category {
		CategoryId: 0,
	}

	err := helpers.ParseForm(r, &category); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	if category.Name == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "name is required")
		return
	}

	result := database.DB.Create(&category) 

	if result.Error != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}


func (c *Category) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	var categories []Category

	result := database.DB.Find(&categories); if result.Error != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, result.Error.Error())
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func (c *Category) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryId, err  := strconv.Atoi(vars["category_id"]); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid category id received")
		return
	}

	var category Category
	result := database.DB.First(&category, categoryId); if result.Error != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, result.Error.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}