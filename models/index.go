package models

import (
	"time"

	"github.com/abdillahzakkie/silkroad/database"
	"gorm.io/gorm"
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
}