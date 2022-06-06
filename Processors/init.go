package Processors

import (
	"fmt"
	"log"
	"os"

	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Client      resty.Client
	DB          *gorm.DB
	RedisClient *redis.Client
	App         *newrelic.Application
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("unable to load .env file")
	}
}

func InitClient() resty.Client {
	Client = *resty.New()
	return Client
}

func Init() {
	LoadEnv()
	Common.LoadConfig()
	//For testing only, update in config.toml
	if os.Getenv("TEST_DEPLOY") == "TRUE" || Common.Config.Adhoc {
		ConnectTestMySQL()
		ConnectTestRedis()
	} else {
		ConnectMySQL()
		ConnectRedis()
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("sea-dinner"),
		newrelic.ConfigLicense("76e67ea9ce3c0608c6a45dcf35496190fed8NRAL"),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		log.Printf("Error initializing newRelic | %v", err.Error())
	}
	App = app
}

func ConnectMySQL() {
	URL := fmt.Sprintf("%v:%v@tcp(%v)/%v", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_URL"), os.Getenv("DB_NAME"))

	log.Printf("Connecting to %v", URL)
	db, err := gorm.Open(mysql.Open(URL), &gorm.Config{})

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
	db, err := gorm.Open(mysql.Open(URL), &gorm.Config{})

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
	log.Println("ConnectRedis: Redis connection established")
	RedisClient = rdb
}

func ConnectTestRedis() {
	redisAddress := fmt.Sprintf("%v:%v", os.Getenv("TEST_REDIS_URL"), os.Getenv("TEST_REDIS_PORT"))
	redisPassword := os.Getenv("TEST_REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       0, // use default DB
	})

	if err := rdb.Ping().Err(); err != nil {
		log.Printf("Error while establishing Test Redis Client: %v", err)
	}
	log.Println("ConnectTestRedis: Redis connection established")
	RedisClient = rdb
}
