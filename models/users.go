package models

import (
	"errors"

	"github.com/abdillahzakkie/silkroad/database"
	"gorm.io/gorm"
)

var (
	ErrorUserNotFound = errors.New("user not found")
	ErrorUserAlreadyExists = errors.New("user already exists")

)

type User struct {
	Model
	ID       uint    `gorm:"primaryKey" json:"id" schema:"-"`
	Wallet   string `gorm:"not null;uniqueIndex" json:"wallet" schema:"wallet,required"`
	Username string `gorm:"not null;uniqueIndex" json:"username" schema:"username,required"`
	Email    string `gorm:"not null;uniqueIndex" json:"-" schema:"email,required"`
	Password string `gorm:"not null" json:"password" schema:"password,required"`
	Product []Product `gorm:"foreignkey:user_id" json:"products" schema:"-"`
}
type UserService struct {
	db *gorm.DB
}

// create new UserService
func NewUserService() *UserService {
	us := UserService{
		db: database.DB,
	}
	return &us
}

