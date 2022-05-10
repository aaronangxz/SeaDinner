package Processors

import (
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
)

var (
	Client resty.Client
	DB     *gorm.DB
)

const (
	UrlPrefix = "https://dinner.sea.com"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("unable to load .env file")
	}
}

func Init() resty.Client {
	Client = *resty.New()
	return Client
}

func ConnectDataBase() {
	database, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		panic("Failed to connect to database!")
	}

	log.Println("Connected to DB")
	DB = database
}
