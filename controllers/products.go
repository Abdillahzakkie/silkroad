package controllers

import (
	"time"

	"github.com/abdillahzakkie/silkroad/database"
)

type Model struct {
	ID        	uint 			`gorm:"primaryKey" json:"id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
	DeletedAt 	*time.Time 	`gorm:"index;<-:update" json:"-"`
}

type Product struct {
	Model
	Name 		string 		`gorm:"<-;not null" json:"name"`
	Description string 		`gorm:"<-;not null" json:"description"`
	Price 		float64 	`gorm:"<-;not null" json:"price"`
	Quantity 	int 		`gorm:"<-;not null" json:"quantity"`
}

func init() {
	// migrate database
	database.DB.AutoMigrate(&Product{})
}