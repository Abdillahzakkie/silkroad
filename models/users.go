package models

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const passwordPepper = "test-pepper"


type UserService struct {
	db *gorm.DB
}

// create new UserService
func NewUserService(psqlInfo string) (*UserService, error) {
	db, err := ConnectDatabase(psqlInfo); if err != nil {
		return nil, err
	}

	us := UserService{
		db: db,
	}
	// auto migrate table
	if err = us.AutoMigrate(); err != nil {
		return nil, err
	}
	return &us, nil
}

// Close method close the database connection
func (us *UserService) Close() error {
	sql, err := us.db.DB(); if err != nil {
		return err
	}
	sql.Close()
	return nil
}

// AutoMigrate will try and automatically migrate table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}); err != nil {
		return errors.New("error while creating new 'users' table")
	}
	return nil
}

// clear existing "users" table and auto migrate the table
// to create new "users" table
func (us *UserService) DestructiveReset() error {
	if err := us.db.Migrator().DropTable("users"); err != nil {
		return errors.New("unable to delete 'users' records")
	}
	// create new tables and index
	if err := us.AutoMigrate(); err != nil {
		return errors.New("error while creating new 'users' table")
	}
	return nil
}

// HashPassword hashes the provided password passed in
func (us *UserService) HashPassword(user *User) error {
	passwordByte := []byte(user.Password + passwordPepper)
	passwordHash, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = ""
	user.PasswordHash = string(passwordHash)
	return nil
}

func (us *UserService) VerifyHashedPassword(password, passwordHash string) error {
	passwordByte := []byte(password + passwordPepper)
	password = ""
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), passwordByte); err != nil {
		switch err {
			case bcrypt.ErrMismatchedHashAndPassword:
				return err
			default:
				return err
		}
	}
	return nil
}

func (us *UserService) Authenticate(email, password string) (User, error) {
	user := User {
		Email: email,
	}

	// find user by  Username or Email
	if err := us.GetUser(&user); err != nil {
		switch err {
			case ErrNotFound:
				return User{}, ErrNotFound
			default:
				return User{}, err
		}
	}
	// verify user's password
	if err := us.VerifyHashedPassword(password, user.PasswordHash); err != nil {
		fmt.Println(err)
		return User{}, ErrInvalidCredentials
	}

	return user, nil
}

func (us *UserService) CreateNewUser(user *User) error {
	if err := us.HashPassword(user); err != nil {
		return err
	}
	// save new user to database
	if err := us.db.Create(&user).Error; err != nil {
		switch err.(type) {
			case *pgconn.PgError:
				return ErrAlreadyExists
			default:
				return ErrInternalServerError
		}
	}
	return nil
}

func (us *UserService) GetAllUsers() ([]User, error) {
	var users []User
	if err := us.db.Find(&users).Error; err != nil {
		return nil, ErrInternalServerError
	}
	return users, nil
}

func (us *UserService) GetUser(user *User) error {
	// checks if user exists
	if err := us.IsExistingUser(*user); err != nil {	
		switch err {
			case ErrNotFound:
				return err
			default:
				return ErrInternalServerError
		}
	}

	if err := us.db.Where(user).First(&user).Error; err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetUserById(id uint) (User, error) {
	user := User { ID: id }

	if err := us.GetUser(&user); err != nil {
		return User{}, err
	}
	return user, nil
}

func (us *UserService) DeleteUserById(id uint) error {
	// checks if user exists
	user, err := us.GetUserById(id); if err != nil {
		return err
	}
	// delete record from database
	if err := us.db.Where(user).Delete(&user).Error; err != nil {
		return ErrInternalServerError
	}
	return nil
}

func (us *UserService) IsExistingUser(user User) error {
	if err := us.db.Where(user).First(&user).Error; err != nil {
		switch err {
			case gorm.ErrRecordNotFound:
				return ErrNotFound
			default:
				return err
		}
	}
	return nil
}