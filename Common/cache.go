package Common

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aaronangxz/SeaDinner/Log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/go-redis/redis"
)

var (
	Redis *redis.Client
)

type RedisWrite struct {
	key    string
	value  interface{}
	expiry time.Duration
}

type RedisRead struct {
	key    string
	entity interface{}
}

type RedisIO struct {
	key    string
	value  interface{}
	entity interface{}
	expiry time.Duration
}

func NewCacheIO() *RedisIO {
	return &RedisIO{
		key:    "",
		value:  nil,
		entity: nil,
		expiry: 0,
	}
}

func NewCacheWrite() *RedisWrite {
	return &RedisWrite{
		key:    "",
		value:  nil,
		expiry: 0,
	}
}

func NewCacheRead() *RedisRead {
	return &RedisRead{
		key:    "",
		entity: nil,
	}
}

func (r *RedisIO) SetKey(key string) *RedisIO {
	r.key = key
	return r
}

func (r *RedisIO) SetEntity(entity interface{}) *RedisIO {
	r.entity = entity
	return r
}

func (r *RedisWrite) SetKey(key string) *RedisWrite {
	r.key = key
	return r
}

func (r *RedisIO) SetValue(value interface{}) *RedisIO {
	r.value = value
	return r
}

func (r *RedisIO) SetExpiry(expiry time.Duration) *RedisIO {
	r.expiry = expiry
	return r
}

func (r *RedisIO) Set() error {
	return SetCache(r.key, r.expiry, r.value)
}

func (r *RedisIO) Get() (interface{}, error) {
	return GetCache(r.key, r.entity)
}

func SetCache(cacheKey string, expiry time.Duration, value interface{}) error {
	v, _ := value.(*sea_dinner.DinnerMenuArray)
	log.Println(v)

	data, err := json.Marshal(value)
	if err != nil {
		Log.Error(ctx, "Failed to marshal results: %v\n", err.Error())
		return err
	}

	if err := Redis.Set(cacheKey, data, expiry).Err(); err != nil {
		Log.Error(ctx, "Error while writing to redis: %v", err.Error())
		return err
	} else {
		Log.Info(ctx, "Successful | Written %v to redis", cacheKey)
		return nil
	}
}

func GetCache(cacheKey string, entity interface{}) (interface{}, error) {
	val, redisErr := Redis.Get(cacheKey).Result()
	if redisErr != nil {
		if redisErr == redis.Nil {
			Log.Warn(ctx, "No result of %v in Redis", cacheKey)
			return nil, redisErr
		} else {
			Log.Error(ctx, "Error while reading from redis: %v", redisErr.Error())
			return nil, redisErr
		}
	} else {
		redisResp := entity
		unmarshalErr := json.Unmarshal([]byte(val), &redisResp)
		if unmarshalErr != nil {
			Log.Warn(ctx, "Fail to unmarshal Redis value of key %v : %v", cacheKey, unmarshalErr)
			return nil, unmarshalErr
		} else {
			Log.Info(ctx, "Successful | Cached %v", cacheKey)
			return redisResp, nil
		}
	}
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
	}
	Log.Info(Ctx, "Live Redis connection established")
	Redis = rdb
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
	}
	Log.Info(Ctx, "Test Redis connection established")
	Redis = rdb
}
