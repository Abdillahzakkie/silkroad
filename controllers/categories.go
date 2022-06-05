package controllers

import (
	"time"

	"github.com/abdillahzakkie/silkroad/database"
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

