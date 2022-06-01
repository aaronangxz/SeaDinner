package Processors

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/go-redis/redis"
)

func GetDayId() (ID int64) {
	var (
		key      = os.Getenv("TOKEN")
		cacheKey = fmt.Sprint(Common.DAY_ID_KEY_PREFIX, ConvertTimeStamp(time.Now().Unix()))
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
		redisRespInt, _ := strconv.Atoi(redisResp)
		return int64(redisRespInt)
	}

	var (
		currentMenu sea_dinner.Current
		currentId   int64
	)

	_, err := Client.R().
		SetHeader("Authorization", MakeToken(key)).
		SetResult(&currentMenu).
		Get(MakeURL(int(sea_dinner.URLType_URL_CURRENT), nil))

	if err != nil {
		fmt.Println(err)
	}

	currentId = currentMenu.GetMenu().GetId()

	if currentMenu.GetMenu().GetPollStart() != fmt.Sprint(ConvertTimeStamp(time.Now().Unix()), "T04:30:00Z") {
		log.Println("GetDayId | Today's ID not found:", currentMenu.GetMenu().GetPollStart())
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
