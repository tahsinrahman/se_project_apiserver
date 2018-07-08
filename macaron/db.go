package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// connects with the db
func InitializeDB(host, port, dbname, username, password string) {
	conn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v", host, port, username, dbname, password)

	var err error
	db, err = gorm.Open("postgres", conn)

	if err != nil {
		log.Fatal(err)
	}

	CreateTables()
}

// create all necessary tables
func CreateTables() {
	//db.AutoMigrate(&User{})
	//db.AutoMigrate(&Company{})
	//db.AutoMigrate(&Job{})
}
