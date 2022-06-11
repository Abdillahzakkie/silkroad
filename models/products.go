package models

import (
	"github.com/abdillahzakkie/silkroad/database"
)

type Product struct {
	Model
	ID 				int 		`gorm:"primaryKey;column:product_id" json:"id" schema:"-"`
	UserID		int			`json:"user_id" schema:"seller_id,required"`	
	CategoryID      int     	`json:"category_id" schema:"category_id,required"`
	ProductName 	string  	`gorm:"not null" json:"product_name" schema:"name,required"`
	Description 	string  	`gorm:"not null" json:"description" schema:"description,required"`
	Price       	float64 	`gorm:"not null" json:"price" schema:"price,required"`
	Quantity    	int     	`gorm:"not null" json:"quantity" schema:"quantity,required"`
}

func (p *Product) CreateNewProduct() error {
	return database.DB.Create(&p).Error
}

func (p *Product) GetAllProducts() (products []Product, err error) {
	err = database.DB.Find(&products).Error; if err != nil {
		return nil, err
	}
	return products, nil
}

func (p Product) GetProduct() (product Product, err error) {
	err = database.DB.Where(p).First(&product).Error; if err != nil {
		return product, err
	}
	return product, nil
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