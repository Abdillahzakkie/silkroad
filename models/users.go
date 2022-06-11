package models

import (
	"errors"

	"github.com/abdillahzakkie/silkroad/database"
	"github.com/abdillahzakkie/silkroad/helpers"
	"github.com/abdillahzakkie/silkroad/vendor/github.com/jackc/pgconn"
	"github.com/abdillahzakkie/silkroad/vendor/golang.org/x/crypto/bcrypt"
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

// Close method close the database connection
func (us *UserService) Close() error {
	sql, err := us.db.DB(); if err != nil {
		return err
	}
	sql.Close()
	return nil
}

func (us *UserService) CreateNewUser(user *User) (err error) {
	// hash user's password
	user.Password, err = helpers.HashPassword(user.Password)
	switch err {
		case bcrypt.ErrHashTooShort:
			return ErrorPasswordTooShort
	}

	err = database.DB.Create(&user).Error;
	switch err.(type) {
		case *pgconn.PgError:
			return ErrorUserAlreadyExists
	}
	return nil
}

func (us *UserService) GetAllUsers() ([]User, error) {
	var users []User
	err := database.DB.Find(&users).Error; if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UserService) GetUserById(id uint) (User, error) {
	var user User

	err := database.
		DB.
		Where("id = ?", id).
		Preload("Product").
		First(&user).Error

	switch err {
		case gorm.ErrRecordNotFound:
			return user, ErrorUserNotFound
	}

	return user, nil
}

func (us *UserService) GetUser() (user User, err error) {
	// err = database.DB.Where(us.User).Preload("Product").First(&user).Error;

	// switch err.(type) {
	// 	case *pgconn.PgError:
	// 		return user, ErrorNotFound
	// 	default:
	// 		return user, ErrorInternalServerError
	// }
	return user, nil
}

func (us *UserService) DeleteUser(id uint) error {
	user := User{
		ID: id,
	}

	err := database.DB.Delete(&user).Error;
	switch err.(type) {
		case *pgconn.PgError:
			return ErrorNotFound
	}
	return nil
}

func (us UserService) IsExistingUser(user User) (bool, error) {
	err := database.DB.Where(user).First(&user).Error; if err != nil {
		switch err {
			case gorm.ErrRecordNotFound:
				return false, nil
			default:
				return false, err
		}
	}
	return true, nil
}