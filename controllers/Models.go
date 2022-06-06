package controllers

import (
	"time"

	"github.com/abdillahzakkie/silkroad/database"
)

type Model struct {
	ID        uint       `gorm:"primaryKey" json:"id" schema:"-"`
	CreatedAt time.Time  `json:"created_at" schema:"-"`
	UpdatedAt time.Time  `json:"-" schema:"-"`
	DeletedAt *time.Time `gorm:"index;<-:update" json:"-" schema:"-"`
}

type User struct {
	Model
	Wallet     string `gorm:"<-;not null;uniqueIndex" json:"wallet" schema:"wallet,required"`
	Username     string `gorm:"<-;not null;uniqueIndex" json:"username" schema:"username,required"`
	Email     string `gorm:"<-;not null;uniqueIndex" json:"-" schema:"email,required"`
	Password  string `gorm:"<-;not null" json:"password" schema:"password,required"`
}

type Product struct {
	Model
	SellerId    uint    `json:"seller_id" schema:"-"`
	Name        string  `gorm:"<-;not null" json:"name" schema:"name,required"`
	Description string  `gorm:"<-;not null" json:"description" schema:"description,required"`
	Category            // `gorm:"embedded" json:"category" schema:"-"`
	Price       float64 `gorm:"<-;not null" json:"price" schema:"price,required"`
	Quantity    int     `gorm:"<-;not null" json:"quantity" schema:"quantity,required"`
}

type Category struct {
	CategoryId		uint       		`gorm:"primaryKey" json:"category_id" schema:"category_id"`
	CreatedAt		time.Time  		`json:"-" schema:"-"`
	UpdatedAt 		time.Time  		`json:"-" schema:"-"`
	DeletedAt 		*time.Time 		`gorm:"index;<-:update" json:"-" schema:"-"`
	Name      		string     		`gorm:"<-;not null;uniqueIndex" json:"name" schema:"name"`
}

func init() {
	// migrate database
	database.DB.AutoMigrate(
		&User{}, 
		&Product{}, 
		&Category{},
	)
}