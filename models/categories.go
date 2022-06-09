package models

import (
	"github.com/abdillahzakkie/silkroad/database"
)

type Category struct {
	Model
	CategoryID int   `gorm:"primaryKey" json:"category_id" schema:"category_id"`
	Name       string `gorm:"<-;not null;uniqueIndex" json:"name" schema:"name"`
}

// POST "/categories/new"
// CreateNewCategory creates new category
func (c *Category) CreateNewCategory() error {
	return database.DB.Create(&c).Error
}

func (c Category) GetAllCategories() (categories []Category, err error) {
	err = database.DB.Find(&categories).Error; if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *Category) GetCategoryById() error {
	return database.DB.Where(c).First(&c).Error
}

func (c *Category) GetCategoryByName() error {
	return database.DB.Where("name = ?", c.Name).First(&c).Error
}