package models

type User struct {
	Model
	ID       int   `gorm:"primaryKey" json:"id" schema:"-"`
	Wallet   string `gorm:"<-;not null;uniqueIndex" json:"wallet" schema:"wallet,required"`
	Username string `gorm:"<-;not null;uniqueIndex" json:"username" schema:"username,required"`
	Email    string `gorm:"<-;not null;uniqueIndex" json:"-" schema:"email,required"`
	Password string `gorm:"<-;not null" json:"password" schema:"password,required"`
}
