package processors

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	_ "github.com/go-sql-driver/mysql" //Required for gorm
	"github.com/joho/godotenv"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	defaultLog "log"
	"os"
	"time"
)

var (
	//Client resty client object
	Client resty.Client
	//DB gorm database object
	DB *gorm.DB
	//RedisClient redis client object
	RedisClient *redis.Client
	//App New Relic object
	App *newrelic.Application
	//Ctx context used for logging
	Ctx context.Context

	newLogger = logger.New(
		defaultLog.New(os.Stdout, "\r\n", defaultLog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
)

//LoadEnv Loads the environment variables from file
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Warn(Ctx, "unable to load .env file")
	}
}

//InitClient Initialize resty client
func InitClient() resty.Client {
	Client = *resty.New()
	return Client
}

//Init Main initialization
func Init() {
	log.InitializeLogger()
	Ctx = context.TODO()
	LoadEnv()
	common.LoadConfig()
	//For testing only, update in config.toml
	if os.Getenv("TEST_DEPLOY") == "TRUE" || common.Config.Adhoc {
		ConnectTestMySQL()
		ConnectTestRedis()
	} else {
		ConnectMySQL()
		ConnectRedis()
	}
	InitRelic()
}

//InitRelic Initialize New Relic
func InitRelic() {
	var (
		appName = "sea-dinner"
		appKey  = os.Getenv("NEWRELIC_KEY")
	)

	if os.Getenv("TEST_DEPLOY") == "TRUE" || common.Config.Adhoc {
		appName = "sea-dinner-test"
		appKey = os.Getenv("TEST_NEWRELIC_KEY")
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(appKey),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		log.Error(Ctx, "Error initializing newRelic | %v", err.Error())
	}
	log.Info(Ctx, "Successfuly initialized newRelic | %v", appName)
	App = app
}

//ConnectMySQL Establish connection to DB
func ConnectMySQL() {
	URL := fmt.Sprintf("%v:%v@tcp(%v)/%v", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_URL"), os.Getenv("DB_NAME"))

	log.Info(Ctx, "Connecting to %v", URL)
	db, err := gorm.Open(mysql.Open(URL), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Error(Ctx, "Error while establishing Live DB Connection: %v", err)
		panic("Failed to connect to live database!")
	}
	log.Info(Ctx, "Live Database connection established")
	DB = db
}

//ConnectTestMySQL Establish connection to test DB
func ConnectTestMySQL() {
	URL := fmt.Sprintf("%v:%v@tcp(%v)/%v", os.Getenv("TEST_DB_USERNAME"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_URL"), os.Getenv("TEST_DB_NAME"))

	log.Info(Ctx, "Connecting to %v", URL)

	db, err := gorm.Open(mysql.Open(URL), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Error(Ctx, "Error while establishing Test DB Connection: %v", err)
		panic("Failed to connect to test database!")
	}
	log.Info(Ctx, "Test Database connection established")
	DB = db
}

//ConnectRedis Establish connection to redis
func ConnectRedis() {
	redisAddress := fmt.Sprintf("%v:%v", os.Getenv("REDIS_URL"), os.Getenv("REDIS_PORT"))
	redisPassword := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       0, // use default DB
	})

	if err := rdb.Ping().Err(); err != nil {
		log.Error(Ctx, "Error while establishing Live Redis client: %v", err)
	}
	log.Info(Ctx, "Live Redis connection established")
	RedisClient = rdb
}

//ConnectTestRedis Establish connection to test redis
func ConnectTestRedis() {
	redisAddress := fmt.Sprintf("%v:%v", os.Getenv("TEST_REDIS_URL"), os.Getenv("TEST_REDIS_PORT"))
	redisPassword := os.Getenv("TEST_REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       0, // use default DB
	})

	if err := rdb.Ping().Err(); err != nil {
		log.Error(Ctx, "Error while establishing Test Redis client: %v", err)
	}
	log.Info(Ctx, "Test Redis connection established")
	RedisClient = rdb
}
