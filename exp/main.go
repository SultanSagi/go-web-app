package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
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
	FirstName string `gorm:"type:varchar(100);not null"`
	LastName string `gorm:"type:varchar(100);not null"`
	BirthDate time.Time `gorm:"not null"`
	Gender string `gorm:"not null"`
	Email string `gorm:"not null"`
	Address string `gorm:"type:varchar(200)"`
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