package test_helper

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"log"
	"math/rand"
	"os"
	"reflect"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//RandomString Generates a random string of length n
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

//RandomInt Generates a random int of maximum max
func RandomInt(max int) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int64(r.Int63n(int64(max)))
}

//LoadEnv Loads env variable from file
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("unable to load .env file")
	}
}

//LoadConfig Loads config.toml
func LoadConfig() {
	common.ConfigPath = "../config.toml"
	if os.Getenv("HEROKU_DEPLOY") == "TRUE" || os.Getenv("TEST_DEPLOY") == "TRUE" {
		common.ConfigPath = "config.toml"
	}

	if _, err := toml.DecodeFile(common.ConfigPath, &common.Config); err != nil {
		log.Fatalln("Reading config failed | ", err, common.ConfigPath)
		return
	}
	log.Println("Reading config OK", common.ConfigPath)
}

//InitTest Initialize test environment
func InitTest() {
	LoadEnv()

	if processors.DB == nil {
		ConnectTestMySQL()
	}

	if processors.RedisClient == nil {
		ConnectTestRedis()
	}
}

//InitClient Initialize resty client
func InitClient() resty.Client {
	processors.Client = *resty.New()
	return processors.Client
}

//ConnectTestMySQL Establish connection fo test DB
func ConnectTestMySQL() {
	URL := fmt.Sprintf("%v:%v@tcp(%v)/%v", os.Getenv("TEST_DB_USERNAME"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_URL"), os.Getenv("TEST_DB_NAME"))

	log.Printf("Connecting to %v", URL)
	db, err := gorm.Open(mysql.Open(URL), &gorm.Config{})

	if err != nil {
		log.Printf("Error while establishing Test DB Connection: %v", err)
		panic("Failed to connect to test database!")
	}

	log.Println("NewMySQL: Test Database connection established")
	processors.DB = db
}

//ConnectTestRedis Establish connection fo test redis
func ConnectTestRedis() {
	redisAddress := fmt.Sprintf("%v:%v", os.Getenv("TEST_REDIS_URL"), os.Getenv("TEST_REDIS_PORT"))
	redisPassword := os.Getenv("TEST_REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       0, // use default DB
	})

	if err := rdb.Ping().Err(); err != nil {
		log.Printf("Error while establishing Test Redis client: %v", err)
	}
	log.Println("ConnectTestRedis: Redis connection established")
	processors.RedisClient = rdb
}

//GetLiveMenuDetails Retrieves live menu details from Sea API
func GetLiveMenuDetails() []*sea_dinner.Food {
	InitClient()
	LoadConfig()
	InitTest()

	//Dynamically set unit_test as true
	//To bypass day_id checks otherwise UT will fail without menu
	common.Config.UnitTest = true
	key := os.Getenv("TOKEN")
	if key == "" {
		log.Println("GetLiveMenuDetails | unable to fetch TOKEN from env")
		return nil
	}
	log.Println("GetLiveMenuDetails | Success")
	return processors.GetMenu(context.TODO(), processors.Client, key).GetFood()
}

//IsInSlice Check if element is in slice
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
