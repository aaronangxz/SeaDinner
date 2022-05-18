package Processors

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var (
	Client      resty.Client
	DB          *gorm.DB
	RedisClient *redis.Client
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("unable to load .env file")
	}
}

func Init() resty.Client {
	LoadEnv()
	LoadConfig()
	//For testing only, update in config.toml
	if Config.Adhoc {
		ConnectTestMySQL()
	} else {
		ConnectMySQL()
	}
	ConnectRedis()
	Client = *resty.New()
	return Client
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

func ConnectRedis() {
	redisAddress := fmt.Sprintf("%v:%v", os.Getenv("REDIS_URL"), os.Getenv("REDIS_PORT"))
	redisPassword := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       0, // use default DB
	})

	if err := rdb.Ping().Err(); err != nil {
		log.Printf("Error while establishing Redis Client: %v", err)
	}
	log.Println("NewRedisClient: Redis connection established")
	RedisClient = rdb
}
