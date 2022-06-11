package models

import (
	"errors"

	"github.com/abdillahzakkie/silkroad/database"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

var (
	ErrorCategoryNotFound = errors.New("category not found")
	ErrorCategoryAlreadyExists = errors.New("category already exists")
)

type Category struct {
	Model
	ID uint   `gorm:"primaryKey" json:"category_id" schema:"category_id"`
	Name       string `gorm:"not null;uniqueIndex" json:"name" schema:"name"`
}

type CategoryService struct {
	db *gorm.DB
}

// create new UserService
func NewCategoryService() *UserService {
	us := UserService{
		db: database.DB,
	}
	return &us
}

// Close method close the database connection
func (cs *CategoryService) Close() error {
	sql, err := cs.db.DB(); if err != nil {
		return err
	}
	sql.Close()
	return nil
}

// POST "/categories/new"
// CreateNewCategory creates new category
func (cs *CategoryService) CreateNewCategory(category Category) (Category, error) {
	err := database.DB.Create(&category).Error;
	switch err.(type) {
		case *pgconn.PgError:
			return category, ErrorCategoryAlreadyExists
	}
	return category, nil
}

func (cs CategoryService) GetAllCategories() ([]Category, error) {
	var categories []Category
	err := database.DB.Find(&categories).Error; if err != nil {
		return nil, err
	}
	return categories, nil
}

func (cs *CategoryService) GetCategoryById(id uint) (Category, error) {
	var category Category

	err := database.DB.Where("id = ?", id).Find(&category).Error;
	switch err {
		case gorm.ErrRecordNotFound:
			return category, ErrorCategoryNotFound
	}
	return category, err
}

func (cs *CategoryService) GetCategoryByName(categoryName string) (Category, error) {
	var category Category
	err := database.DB.Where("name = ?", category.Name).First(&category).Error;
	switch err {
		case gorm.ErrRecordNotFound:
			return category, ErrorCategoryNotFound
	}
	return category, nil
}

func (cs *CategoryService) IsExistingCategory(id uint) (bool, error) {
	var category Category
	err := database.DB.Where("id = ?", id).First(&category).Error; if err != nil {
		switch err {
			case gorm.ErrRecordNotFound:
				return false, nil
			default:
				return false, err
		}
	}
	return true, nil
}