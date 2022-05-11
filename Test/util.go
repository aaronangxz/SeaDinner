package Test

import (
	"log"

	"github.com/jinzhu/gorm"
)

var (
	TestDB *gorm.DB
)

func ConnectDataBase() {
	database, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		panic("Failed to connect to database!")
	}

	log.Println("Connected to DB")
	TestDB = database
}
