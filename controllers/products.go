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

type Model struct {
	ID        uint       `gorm:"primaryKey" json:"id" schema:"-"`
	CreatedAt time.Time  `json:"created_at" schema:"-"`
	UpdatedAt time.Time  `json:"updated_at" schema:"-"`
	DeletedAt *time.Time `gorm:"index;<-:update" json:"-" schema:"-"`
}

type Product struct {
	Model
	SellerId 	uint 				`json:"seller_id" schema:"-"`
	Name 		string 				`gorm:"<-;not null" json:"name" schema:"name,required"`
	Description string 				`gorm:"<-;not null" json:"description" schema:"description,required"`
	Category 						// `gorm:"embedded" json:"category" schema:"-"`
	Price 		float64 			`gorm:"<-;not null" json:"price" schema:"price,required"`
	Quantity 	int 				`gorm:"<-;not null" json:"quantity" schema:"quantity,required"`
}

func init() {
	// migrate database
	database.DB.AutoMigrate(&Product{})
}

// POST "/products/{seller_id}/new"
// CreateNewProduct creates new product
func (p Product) CreateNewProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sellerId, err  := strconv.Atoi(vars["seller_id"]); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid seller id received")
		return
	}

	product := Product{
		SellerId: uint(sellerId),
	}

	err = helpers.ParseForm(r, &product); if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checks if category ID exists
	var category Category
	result := database.DB.Find(&category, product.Category.CategoryId); if result.Error != nil {
		helpers.RespondWithError(
			w, 
			http.StatusBadRequest, 
			fmt.Sprintf("category id %v does not exist", product.Category.CategoryId),
		)
		return
	}
	if category.CategoryId == 0 {
		helpers.RespondWithError(w, http.StatusBadRequest, "category id does not exist")
		return
	}

	result = database.DB.Create(&product); if result.Error != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}