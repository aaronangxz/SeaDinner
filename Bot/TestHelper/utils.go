package TestHelper

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aaronangxz/SeaDinner/Processors"
	"github.com/jinzhu/gorm"
)

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func RandomInt(max int) int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(max))
}

func InitTest() {
	ConnectTestDataBase()
}

func ConnectTestDataBase() {
	database, err := gorm.Open("sqlite3", "../test.db")

	if err != nil {
		panic("Failed to connect to test database!")
	}

	log.Println("Connected to Test DB")
	Processors.DB = database
}
