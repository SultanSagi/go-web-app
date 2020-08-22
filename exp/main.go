package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host = "localhost"
	port = 5432
	user = "sultan"
	password = "narivo979"
	dbname = "gowebapp_dev"
)

type Customer struct {
	gorm.Model
	FirstName string
	LastName string
	Gender string
	Email string `gorm:"not null"`
	Address string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.LogMode(true)
	db.AutoMigrate(&Customer{})
	var c Customer
	db.First(&c)
	fmt.Println(c)
}