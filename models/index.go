package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "test"
	dbname   = "silkroad"
)

var (
	// ErrInternalServerError is returned for internal server errors
	ErrInternalServerError = errors.New("models: internal server error")
	// ErrNotFound is returned when a resources cannot be found in the database
	ErrNotFound = errors.New("models: resource not found")
	// ErrAlreadyExists is returned when a resources already existed in the database
	ErrAlreadyExists = errors.New("models: resource already exists")
	// ErrInvalidCredentials is returned is credentials are invalid
	ErrInvalidCredentials = errors.New("models: invalid credentials provided")
	// ErrDatabaseConnectionFailed is returned if database connection cannot be established
	ErrDatabaseConnectionFailed = errors.New("database connection failed")
	// ErrPasswordTooShort is returned if password length is too short
	ErrPasswordTooShort = errors.New("models: password is too short")
	
)

type Model struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	Model
	ID       			uint    	`gorm:"primaryKey" json:"id"`
	Remember 			string		`gorm:"-" json:"-"`
	RememberHash		string		`gorm:"not null;unique_index"`
	Wallet   			string 		`gorm:"not null;uniqueIndex" json:"wallet"`
	Username 			string 		`gorm:"not null;uniqueIndex" json:"username"`
	Email    			string 		`gorm:"not null;uniqueIndex" json:"-"`
	Password 			string 		`gorm:"-:all" json:"-"`
	PasswordHash 		string 		`gorm:"not null;column:password" json:"-"`
	// Product 			[]Product 	`gorm:"foreignkey:user_id" json:"products"`
}

type Category struct {
	Model
	ID 			uint   		`gorm:"primaryKey" json:"category_id"`
	Name       	string 		`gorm:"not null;uniqueIndex" json:"name"`
}

// ConnectDatabase connects to the database
func ConnectDatabase(psqlInfo string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: psqlInfo,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return db, nil
}

// GetPsqlInfo returns a string containing the PostgreSQL connection information
func GetPsqlInfo(dbname string) (string, error) {
	if dbname == "" {
		return "", errors.New("models: dbname is empty")
	}
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", 
		host, 
		port, 
		user, 
		password, 
		dbname,
	)
	return psqlInfo, nil
}