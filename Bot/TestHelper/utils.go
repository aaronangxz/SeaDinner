package TestHelper

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/aaronangxz/SeaDinner/Processors"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func RandomInt(max int) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int64(r.Intn(max))
}
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("unable to load .env file")
	}
}

func InitTest() {
	LoadEnv()
	ConnectTestMySQL()
}

func ConnectTestMySQL() {
	URL := fmt.Sprintf("%v:%v@tcp(%v)/%v", os.Getenv("TEST_DB_USERNAME"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_URL"), os.Getenv("TEST_DB_NAME"))

	log.Printf("Connecting to %v", URL)
	// sqlDB, err := sql.Open("mysql", URL)
	db, err := gorm.Open(mysql.Open(URL), &gorm.Config{})

	if err != nil {
		log.Printf("Error while establishing Test DB Connection: %v", err)
		panic("Failed to connect to test database!")
	}

	log.Println("NewMySQL: Test Database connection established")
	Processors.DB = db
}
