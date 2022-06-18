package models

import (
	"errors"
	"testing"
	"time"
)

func testingCategoryService() (*CategoryService, error) {
	psqlInfo, err := GetPsqlInfo("silkroad_test")
	if err != nil {
		return nil, err
	}

	cs, err := NewCategoryService(psqlInfo)
	if err != nil {
		return nil, err
	}
	// clear the users table between test
	err = cs.DestructiveReset()
	if err != nil {
		return cs, err
	}
	return cs, nil
}

func TestCreateNewCategory(t *testing.T) {
	cs, err := testingCategoryService(); if err != nil {
		t.Errorf(err.Error())
	}
	// close database connection
	defer cs.Close()
	
	// create new category
	category := Category{
		Name: "testing",
	}

	if err := cs.CreateNewCategory(&category); err != nil {
		t.Error(err)
	}

	if category.ID == 0 {
		t.Errorf(errorResponse(1, category.ID))
	}
	if time.Since(category.CreatedAt) > time.Duration(5 * time.Second) ||
	time.Since(category.UpdatedAt) > time.Duration(5 * time.Second) {
		t.Error("invalid createdAt & updatedAt timestamp")
	}
}

func TestGetAllCategories(t *testing.T) {
	cs, err := testingCategoryService(); if err != nil {
		t.Error(err)
	}
	// close database connection
	defer cs.Close()

	// create new category
	category := Category{
		Name: "testing category",
	}

	if err := cs.CreateNewCategory(&category); err != nil {
		t.Error(err)
	}

	// get all categories
	categories, err := cs.GetAllCategories(); if err != nil {
		t.Error(err)
	}

	if categoryLength:= len(categories);  categoryLength != 1 {
		t.Errorf("Expected categories to have 1 item, received %d", categoryLength)
	}
}

func TestGetCategoryById(t *testing.T) {
	cs, err := testingCategoryService(); if err != nil {
		t.Error(err)
	}
	// close database connection
	defer cs.Close()

	// create new category
	category := Category{
		Name: "testing",
	}

	if err := cs.CreateNewCategory(&category); err != nil {
		t.Error(err)
	}

	// get category by ID
	category, err = cs.GetCategoryById(category.ID); if err != nil {
		t.Error(err)
	}

	if category.ID == 0 {
		t.Errorf(errorResponse(0, 1))
	}
}

func TestCategoryNotFound(t *testing.T) {
	cs, err := testingCategoryService(); if err != nil {
		t.Error(err)
	}
	// close database connection
	defer cs.Close()

	// get category by ID
	category, err := cs.GetCategoryById(100); if err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Error("should reject invalid category")
		}
		return
	}

	if categoryID := category.ID; categoryID > 0 {
		t.Errorf(errorResponse(0, categoryID))
	}
}