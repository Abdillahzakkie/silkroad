package models

import (
	"errors"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)


type CategoryService struct {
	db *gorm.DB
}

// create new CategoryService
func NewCategoryService(psqlInfo string) (*CategoryService, error) {
	db, err := ConnectDatabase(psqlInfo); if err != nil {
		return nil, err
	}
	cs := CategoryService{
		db: db,
	}

	// auto migrate table
	if err := cs.AutoMigrate(); err != nil {
		return nil, err
	}
	return &cs, nil
}

// Close method close the database connection
func (cs *CategoryService) Close() error {
	sql, err := cs.db.DB(); if err != nil {
		return err
	}
	sql.Close()
	return nil
}

// AutoMigrate will try and automatically migrate table
func (us *CategoryService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&Category{}); err != nil {
		return errors.New("error while creating new 'categories' table")
	}
	return nil
}

// clear existing "users" table and auto migrate the table
// to create new "users" table
func (cs *CategoryService) DestructiveReset() error {
	if err := cs.db.Migrator().DropTable("categories"); err != nil {
		return errors.New("unable to delete 'categories' records")
	}
	// create new tables and index
	if err := cs.AutoMigrate(); err != nil {
		return errors.New("error while creating new 'categories' table")
	}
	return nil
}

// POST "/categories/new"
// CreateNewCategory creates new category
func (cs *CategoryService) CreateNewCategory(category *Category) error {
	if err := cs.db.Create(category).Error; err != nil {
		switch err.(type) {
			case *pgconn.PgError:
				return ErrAlreadyExists
			default:
				return ErrInternalServerError
		}
	}
	
	return nil
}

func (cs CategoryService) GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := cs.db.Find(&categories).Error; err != nil {
		return nil, ErrInternalServerError
	}
	return categories, nil
}

func (cs *CategoryService) GetCategoryById(id uint) (Category, error) {
	category := Category{
		ID: id,
	}
	// checks if user exists
	if err := cs.IsExistingCategory(category); err != nil {
		return category, err
	}

	if err := cs.db.Where(category).First(&category).Error; err != nil {
		switch err {
			case gorm.ErrRecordNotFound:
				return category, ErrNotFound
			default:
				return category, err
		}
	}
	return category, nil
}

func (cs *CategoryService) GetCategoryByName(categoryName string) (Category, error) {
	category := Category{
		Name: categoryName,
	}
	if err := cs.db.Where(category).First(&category).Error; err != nil {
		switch err {
			case gorm.ErrRecordNotFound:
				return category, ErrNotFound
			default:
				return category, err
		}
	}
	
	return category, nil
}

func (cs *CategoryService) IsExistingCategory(category Category) error {
	if err := cs.db.Where(category).First(&category).Error; err != nil {
		switch err {
			case gorm.ErrRecordNotFound:
				return ErrNotFound
			default:
				return err
		}
	}
	return nil
}