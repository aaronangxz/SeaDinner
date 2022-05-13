package Processors

import (
	"log"
	"os"

	"github.com/go-resty/resty/v2"
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

func ConnectTestDataBase() {
	database, err := gorm.Open("sqlite3", "../test.db")

	if err != nil {
		panic("Failed to connect to database!")
	}

	log.Println("Connected to DB")
	DB = database
}
