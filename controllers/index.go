package controllers

import (
	"log"

	"github.com/abdillahzakkie/silkroad/models"
)

var (
	us *models.UserService
	cs *models.CategoryService
)

func init() {
	psqlInfo, err := models.GetPsqlInfo("silkroad"); if err != nil {
		log.Fatal(err)
	}

	us, err = models.NewUserService(psqlInfo); if err != nil {
		log.Fatal(models.ErrDatabaseConnectionFailed.Error())
	}

	cs, err = models.NewCategoryService(psqlInfo); if err != nil {
		log.Fatal(models.ErrDatabaseConnectionFailed.Error())
	}

	// clear all tables
	// us.DestructiveReset()
	// cs.DestructiveReset()
}