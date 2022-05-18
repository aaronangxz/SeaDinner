package Processors

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func GetDayId() (ID int) {
	var (
		key      = os.Getenv("TOKEN")
		cacheKey = fmt.Sprint(DAY_ID_KEY_PREFIX, ConvertTimeStamp(time.Now().Unix()))
		expiry   = 86400 * time.Second
	)

	//check cache
	redisResp, redisErr := RedisClient.Get(cacheKey).Result()
	if redisErr != nil {
		if redisErr == redis.Nil {
			log.Printf("GetDayId | No result of %v in Redis, reading from API", cacheKey)
		} else {
			log.Printf("GetDayId | Error while reading from redis: %v", redisErr.Error())
		}
	} else {
		log.Printf("GetDayId | Successful | Cached %v", cacheKey)
		redisRespInt, _ := strconv.Atoi(redisResp)
		return redisRespInt
	}

	var (
		currentMenu Current
		currentId   int
	)

	_, err := Client.R().
		SetHeader("Authorization", MakeToken(key)).
		SetResult(&currentMenu).
		Get(MakeURL(URL_CURRENT, nil))

	if err != nil {
		fmt.Println(err)
	}

	currentId = currentMenu.Menu.GetId()

	if currentMenu.Menu.GetPollStart() != fmt.Sprint(ConvertTimeStamp(time.Now().Unix()), "T04:30:00Z") {
		log.Println("GetDayId | Today's ID not found:", currentMenu.Menu.GetPollStart())
		currentId = 0
	}

	//set back into cache
	if err := RedisClient.Set(cacheKey, currentId, expiry).Err(); err != nil {
		log.Printf("GetDayId | Error while writing to redis: %v", err.Error())
	} else {
		log.Printf("GetDayId | Successful | Written %v to redis", cacheKey)
	}

	return currentId
}
