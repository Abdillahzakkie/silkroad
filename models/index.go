package models

import (
	"errors"
	"time"

	"github.com/abdillahzakkie/silkroad/database"
	"gorm.io/gorm"
)

var (
	ErrorNotFound = errors.New("models: resource not found")
	ErrorAlreadyExists = errors.New("models: resource already exists")
	ErrorInternalServerError = errors.New("models: internal server error")

	// BCRYPT ERRORS
	ErrorPasswordTooShort = errors.New("models: password is too short")
)

type Model struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func init() {
	// migrate database
	database.DB.AutoMigrate(
		&User{}, 
		&Product{},
		&Category{},
	)

	// err := destructiveReset(); if err != nil {
	// 	log.Fatal(err)
	// }
}

// clear existing "users" table and auto migrate the table
// to create new "users" table
func destructiveReset() error {
	database.DB.Migrator().DropTable("users", "products", "categories")

	// create new tables and index
	err := database.DB.AutoMigrate(
		&User{},
		&Product{},
		&Category{},
	)
	
	if err != nil {
		return err
	}
	return nil
}