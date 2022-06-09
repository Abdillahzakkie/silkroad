package models

import (
	"github.com/abdillahzakkie/silkroad/database"
)

type Product struct {
	Model
	ProductID 		int 		`gorm:"primaryKey" json:"id" schema:"-"`
	SellerID		int			`gorm:"embedded" json:"seller_id" schema:"seller_id,required"`	
	ProductName 	string  	`gorm:"<-;not null" json:"name" schema:"name,required"`
	Description 	string  	`gorm:"<-;not null" json:"description" schema:"description,required"`
	CategoryID      int     	`gorm:"embedded" json:"-" schema:"category_id,required"`
	Price       	float64 	`gorm:"<-;not null" json:"price" schema:"price,required"`
	Quantity    	int     	`gorm:"<-;not null" json:"quantity" schema:"quantity,required"`
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
	return database.DB.Where("product_id = ?", p.ProductID).First(&p).Error
}