package Processors

import (
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
)

var (
	Client resty.Client
	DB     *gorm.DB
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("unable to load .env file")
	}
}

func Init() resty.Client {
	// Path to config file can be passed in.
	LoadConfig()

	Client = *resty.New()
	return Client
}

func ConnectDataBase() {
	dbName := "store.db"
	if os.Getenv("HEROKU_DEPLOY") == "FALSE" {
		dbName = "../store.db"
	}
	database, err := gorm.Open("sqlite3", dbName)

	if err != nil {
		panic("Failed to connect to database!")
	}

	log.Println("Connected to DB")
	DB = database
}

func ConnectMySQL() {
	URL := fmt.Sprintf("%v:%v@tcp(%v)/%v", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_URL"), os.Getenv("DB_NAME"))

	log.Printf("Connecting to %v", URL)
	db, err := gorm.Open("mysql", URL)

	if err != nil {
		log.Printf("Error while establishing DB Connection: %v", err)
		panic("Failed to connect to database!")
	}

	log.Println("NewMySQL: Database connection established")
	DB = db
}

func ConnectTestMySQL() {
	URL := fmt.Sprintf("%v:%v@tcp(%v)/%v", os.Getenv("TEST_DB_USERNAME"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_URL"), os.Getenv("TEST_DB_NAME"))

	log.Printf("Connecting to %v", URL)
	db, err := gorm.Open("mysql", URL)

	if err != nil {
		log.Printf("Error while establishing Test DB Connection: %v", err)
		panic("Failed to connect to test database!")
	}

	log.Println("NewMySQL: Test Database connection established")
	DB = db
}
