package models

import (
	"errors"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

// import (
// 	"errors"

// 	"github.com/abdillahzakkie/silkroad/database"
// 	"github.com/jackc/pgconn"
// 	"gorm.io/gorm"
// )

// var (
// 	ErrorProductAlreadyExists = errors.New("product already exists")
// 	ErrorProductNotFound = errors.New("product not found")
// )

type Product struct {
	Model
	ID 				uint 		`gorm:"primaryKey;column:product_id" json:"id" schema:"-"`
	UserID			uint		`json:"user_id" schema:"seller_id,required"`
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
// create new UserService
func NewProductService(psqlInfo string) (*ProductService, error) {
	db, err := ConnectDatabase(psqlInfo); if err != nil {
		return nil, err
	}
	ps := ProductService{
		db: db,
	}

	// auto migrate table
	err = ps.AutoMigrate(); if err != nil {
		return nil, err
	}
	return &ps, nil
}

// Close method close the database connection
func (ps *ProductService) Close() error {
	sql, err := ps.db.DB(); if err != nil {
		return err
	}
	sql.Close()
	return nil
}


// AutoMigrate will try and automatically migrate table
func (us *ProductService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&Product{}); err != nil {
		return errors.New("error while creating new 'users' table")
	}
	return nil
}

// clear existing "users" table and auto migrate the table
// to create new "users" table
func (ps *ProductService) DestructiveReset() error {
	if err := ps.db.Migrator().DropTable("products"); err != nil {
		return errors.New("unable to delete 'products' records")
	}
	// create new tables and index
	if err := ps.AutoMigrate(); err != nil {
		return errors.New("error while creating new 'products' table")
	}
	return nil
}

func (ps *ProductService) CreateNewProduct(product Product) (Product, error) {
	if err := ps.db.Create(&product).Error; err != nil {
		switch err.(type) {
			case *pgconn.PgError:
				return product, ErrAlreadyExists
			default:
				return product, ErrInternalServerError
		}
	}
	return product, nil
}

func (ps *ProductService) GetAllProducts() ([]Product, error) {
	var products []Product
	if err := ps.db.Find(&products).Error; err != nil {
		return nil, ErrInternalServerError
	}
	return products, nil
}

func (ps *ProductService) GetProductById(id uint) (Product, error) {
	product := Product{
		ID: id,
	}

	// checks whether product already exists
	if err := ps.IsExistingProduct(product); err != nil {
		return product, err
	}

	// validate if seller exists
	psqlInfo, err :=  GetPsqlInfo(dbname); if err != nil {
		return product, err
	}

	if us, err := NewUserService(psqlInfo); err != nil {
		defer us.Close()
		user := User {
			ID: product.UserID,
		}

		if err := us.IsExistingUser(user); err != nil {
			return product, err
		}

	}
	
	if err := ps.db.Where(product).First(&product).Error; err != nil {
		switch err {
			case gorm.ErrRecordNotFound:
				return product, ErrNotFound
			default:
				return product, err
		}
	}
	return product, nil
}

func (ps *ProductService) IsExistingProduct(product Product) error {
	if err := ps.db.Where(product).First(&product).Error; err != nil {
		switch err {
			case gorm.ErrRecordNotFound:
				return ErrNotFound
			default:
				return err
		}
	}
	return nil
}