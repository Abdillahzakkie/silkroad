package models

import (
	"errors"
	"fmt"

	"github.com/abdillahzakkie/silkroad/database"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

var (
	ErrorProductAlreadyExists = errors.New("product already exists")
	ErrorProductNotFound = errors.New("product not found")
)

type Product struct {
	Model
	ID 				uint 		`gorm:"primaryKey;column:product_id" json:"id" schema:"-"`
	UserID			uint			`json:"user_id" schema:"seller_id,required"`	
	CategoryID      uint     	`json:"category_id" schema:"category_id,required"`
	ProductName 	string  	`gorm:"not null" json:"product_name" schema:"name,required"`
	Description 	string  	`gorm:"not null" json:"description" schema:"description,required"`
	Price       	float64 	`gorm:"not null" json:"price" schema:"price,required"`
	Quantity    	int     	`gorm:"not null" json:"quantity" schema:"quantity,required"`
}

type ProductService struct {
	db *gorm.DB
}

// create new ProductService
func NewProductService() *ProductService {
	us := ProductService{
		db: database.DB,
	}
	return &us
}

// Close method close the database connection
func (us *ProductService) Close() error {
	sql, err := us.db.DB(); if err != nil {
		return err
	}
	sql.Close()
	return nil
}

func (ps *ProductService) CreateNewProduct(product Product) (Product, error) {
	err := database.DB.Create(&product).Error
	switch err.(type) {
		case *pgconn.PgError:
			return product, ErrorProductAlreadyExists
	}
	return product, err
}

func (ps *ProductService) GetAllProducts() (products []Product, err error) {
	err = database.DB.Find(&products).Error; if err != nil {
		return nil, err
	}
	return products, nil
}

func (ps ProductService) GetProduct(product Product) (Product, error) {
	err := database.DB.Where(product).First(&product).Error
	fmt.Println(err)
	switch err {
		case gorm.ErrRecordNotFound:
			return product, ErrorProductNotFound
	}
	return product, err
}

func (p *Product) GetProductById() error {
	return database.DB.Where("product_id = ?", p.ID).First(&p).Error
}

func (p *Product) GetProductsBySellerId() (products []Product, err error) {
	err = database.DB.Where("user_id = ?", p.UserID).Find(&products).Error; if err != nil {
		return nil, err
	}
	return products, nil
}

func (ps ProductService) IsExistingProduct(product Product) error {
	err := database.DB.Where(product).First(&product).Error; if err != nil {
		switch err {
			case gorm.ErrRecordNotFound:
				return ErrorUserNotFound
		}
	}
	return err
}