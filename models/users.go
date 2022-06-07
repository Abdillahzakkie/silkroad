package models

import (
	"github.com/abdillahzakkie/silkroad/database"
	"github.com/abdillahzakkie/silkroad/helpers"
)

type User struct {
	Model
	ID       int    `gorm:"primaryKey" json:"id" schema:"-"`
	Wallet   string `gorm:"<-;not null;uniqueIndex" json:"wallet" schema:"wallet,required"`
	Username string `gorm:"<-;not null;uniqueIndex" json:"username" schema:"username,required"`
	Email    string `gorm:"<-;not null;uniqueIndex" json:"-" schema:"email,required"`
	Password string `gorm:"<-;not null" json:"password" schema:"password,required"`
}

func (u *User) CreateNewUser() error {
	// hash user's password
	password, err := helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(password)

	result := database.DB.Create(&u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}


func (u User) GetAllUsers() ([]User, error) {
	var users []User
	result := database.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}


func (u *User) GetUserById() error {
	result := database.DB.Where("id = ?", u.ID).First(&u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) DeleteUser() error {
	result := database.DB.Delete(&u, u.ID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

