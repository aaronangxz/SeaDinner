package TestHelper

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/aaronangxz/SeaDinner/Processors"
	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
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

	if Processors.DB == nil {
		ConnectTestMySQL()
	}

	if Processors.RedisClient == nil {
		ConnectRedis()
	}
}

func InitClient() resty.Client {
	Processors.Client = *resty.New()
	return Processors.Client
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
	Processors.DB = db
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
	Processors.RedisClient = rdb
}

func LoadConfig() {
	Processors.ConfigPath = "../config.toml"
	if os.Getenv("HEROKU_DEPLOY") == "TRUE" {
		Processors.ConfigPath = "config.toml"
	}

	if _, err := toml.DecodeFile(Processors.ConfigPath, &Processors.Config); err != nil {
		log.Fatalln("Reading config failed | ", err, Processors.ConfigPath)
		return
	}
	log.Println("Reading config OK", Processors.ConfigPath)
}

func GetFromCache(cacheKey string, model interface{}) interface{} {
	val, redisErr := Processors.RedisClient.Get(cacheKey).Result()
	if redisErr != nil {
		if redisErr == redis.Nil {
			log.Printf("GetCache | No result of %v in Redis, reading from DB", cacheKey)
		} else {
			log.Printf("GetCache | Error while reading from redis: %v", redisErr.Error())
		}
	} else {
		redisResp := model
		err := json.Unmarshal([]byte(val), &redisResp)
		if err != nil {
			log.Printf("GetCache | Fail to unmarshal Redis value of key %v : %v, reading from DB", cacheKey, err)
		} else {
			log.Printf("GetCache | Successful | Cached %v", cacheKey)
			return redisResp
		}
	}
	return nil
}

func GetLiveMenuDetails() []Processors.Food {
	InitClient()
	LoadConfig()
	InitTest()
	key := os.Getenv("TOKEN")
	return Processors.GetMenu(Processors.Client, Processors.GetDayId(), key).DinnerArr
}
