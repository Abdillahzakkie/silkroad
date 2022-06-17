package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	gorm "gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "test"
	dbname   = "silkroad"
)

var DB *gorm.DB


func init() {
	var err error
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", 
		host, 
		port, 
		user, 
		password, 
		dbname,
	)

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: psqlInfo,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Database connected!")
}
