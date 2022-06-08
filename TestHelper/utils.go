package TestHelper

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/Processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
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

func LoadConfig() {
	Common.ConfigPath = "../config.toml"
	if os.Getenv("HEROKU_DEPLOY") == "TRUE" || os.Getenv("TEST_DEPLOY") == "TRUE" {
		Common.ConfigPath = "config.toml"
	}

	if _, err := toml.DecodeFile(Common.ConfigPath, &Common.Config); err != nil {
		log.Fatalln("Reading config failed | ", err, Common.ConfigPath)
		return
	}
	log.Println("Reading config OK", Common.ConfigPath)
}

func InitTest() {
	LoadEnv()

	if Processors.DB == nil {
		ConnectTestMySQL()
	}

	if Processors.RedisClient == nil {
		ConnectTestRedis()
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
	Processors.RedisClient = rdb
}

func GetLiveMenuDetails() []*sea_dinner.Food {
	InitClient()
	LoadConfig()
	InitTest()
	key := os.Getenv("TOKEN")
	if key == "" {
		log.Println("GetLiveMenuDetails | unable to fetch TOKEN from env")
		return nil
	}
	log.Println("GetLiveMenuDetails | Success")
	return Processors.GetMenuUsingCache(Processors.Client, key).GetFood()
}

func IsInSlice(a interface{}, slice interface{}) bool {
	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice, reflect.Ptr:
		values := reflect.Indirect(reflect.ValueOf(slice))
		val := reflect.Indirect(reflect.ValueOf(a))

		if values.Len() == 0 || val.Len() == 0 {
			return false
		}

		if val.Index(0).Kind() != values.Index(0).Kind() {
			return false
		}

		for i := 0; i < values.Len(); i++ {
			for j := 0; j < val.Len(); j++ {
				if reflect.DeepEqual(values.Index(i).Interface(), val.Index(j).Interface()) {
					return true
				}
			}
		}
	}
	return false
}