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

func (u *User) CreateNewUser() (err error) {
	// hash user's password
	password, err := helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(password)

	err = database.DB.Create(&u).Error; if err != nil {
		return err
	}

	return nil
}


func (u User) GetAllUsers() (users []User, err error) {
	err = database.DB.Find(&users).Error; if err != nil {
		return nil, err
	}
	return users, nil
}


func (u *User) GetUserById() (err error) {
	err = database.DB.Where(u).First(&u).Error; if err != nil {
		return err
	}
	return nil
}

func (u *User) GetUser() (err error) {
	err = database.DB.Where(u).First(&u).Error; if err != nil {
		return err
	}
	return nil
}

func (u *User) DeleteUser() (err error) {
	err = database.DB.Where(u).Delete(&u).Error; if err != nil {
		return err
	}
	return nil
}

func (u *User) IsExisted() bool {
	database.DB.Where("id = ? OR username = ? OR wallet = ?", u.ID, u.Username, u.Wallet).First(&u)
	return u.ID != 0
}