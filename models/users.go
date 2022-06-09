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
	return  database.DB.Create(&u).Error
}


func (u User) GetAllUsers() (users []User, err error) {
	err = database.DB.Find(&users).Error; if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *User) GetUser() error {
	return database.DB.Where(u).First(&u).Error;
}

func (u *User) DeleteUser() error {
	return database.DB.Where(u).Delete(&u).Error;
}