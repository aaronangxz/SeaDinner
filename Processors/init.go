package Processors

import (
	"context"
	"fmt"
	"os"

	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/Log"
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
	Ctx         context.Context
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		Log.Warn(Ctx, "unable to load .env file")
		// log.Printf("unable to load .env file")
	}
}

func InitClient() resty.Client {
	Client = *resty.New()
	return Client
}

func Init() {
	Log.InitializeLogger()
	Ctx = context.TODO()
	LoadEnv()
	Common.LoadConfig()
	//For testing only, update in config.toml
	if os.Getenv("TEST_DEPLOY") == "TRUE" || Common.Config.Adhoc {
		ConnectTestMySQL()
		Common.ConnectTestRedis()
		ConnectTestRedis()
	} else {
		ConnectMySQL()
		Common.ConnectRedis()
		ConnectRedis()
	}
	InitRelic()
}

func InitRelic() {
	var (
		appName = "sea-dinner"
		appKey  = os.Getenv("NEWRELIC_KEY")
	)

	if os.Getenv("TEST_DEPLOY") == "TRUE" || Common.Config.Adhoc {
		appName = "sea-dinner-test"
		appKey = os.Getenv("TEST_NEWRELIC_KEY")
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(appKey),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		Log.Error(Ctx, "Error initializing newRelic | %v", err.Error())
		// log.Printf("Error initializing newRelic | %v", err.Error())
	}
	Log.Info(Ctx, "Successfuly initialized newRelic | %v", appName)
	// log.Printf("Successfuly initialized newRelic | %v", appName)
	App = app
}

func ConnectMySQL() {
	URL := fmt.Sprintf("%v:%v@tcp(%v)/%v", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_URL"), os.Getenv("DB_NAME"))

	Log.Info(Ctx, "Connecting to %v", URL)
	// log.Printf("Connecting to %v", URL)
	db, err := gorm.Open(mysql.Open(URL), &gorm.Config{})

	if err != nil {
		Log.Error(Ctx, "Error while establishing Live DB Connection: %v", err)
		// log.Printf("Error while establishing DB Connection: %v", err)
		panic("Failed to connect to live database!")
	}
	Log.Info(Ctx, "Live Database connection established")
	// log.Println("NewMySQL: Database connection established")
	DB = db
}

func ConnectTestMySQL() {
	URL := fmt.Sprintf("%v:%v@tcp(%v)/%v", os.Getenv("TEST_DB_USERNAME"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_URL"), os.Getenv("TEST_DB_NAME"))

	Log.Info(Ctx, "Connecting to %v", URL)
	// log.Printf("Connecting to %v", URL)
	db, err := gorm.Open(mysql.Open(URL), &gorm.Config{})

	if err != nil {
		Log.Error(Ctx, "Error while establishing Test DB Connection: %v", err)
		// log.Printf("Error while establishing Test DB Connection: %v", err)
		panic("Failed to connect to test database!")
	}
	Log.Info(Ctx, "Test Database connection established")
	// log.Println("NewMySQL: Test Database connection established")
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
		Log.Error(Ctx, "Error while establishing Live Redis Client: %v", err)
		// log.Printf("Error while establishing Redis Client: %v", err)
	}
	Log.Info(Ctx, "Live Redis connection established")
	// log.Println("ConnectRedis: Redis connection established")
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
		Log.Error(Ctx, "Error while establishing Test Redis Client: %v", err)
		// log.Printf("Error while establishing Test Redis Client: %v", err)
	}
	Log.Info(Ctx, "Test Redis connection established")
	// log.Println("ConnectTestRedis: Redis connection established")
	RedisClient = rdb
}
