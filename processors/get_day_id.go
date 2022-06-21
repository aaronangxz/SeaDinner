package processors

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"os"
	"strconv"
	"time"

	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/go-redis/redis"
)

//GetDayID Calls Sea API, retrieves the current day's id
func GetDayID(ctx context.Context) (ID int64) {
	var (
		key = os.Getenv("TOKEN")
		//12 hours offset, so we don't try to check between 0000 ~ 1200 when day id isn't updated yet
		cacheKey = fmt.Sprint(common.DAY_ID_KEY_PREFIX, ConvertTimeStamp(time.Now().Unix()-43200))
		expiry   = 172800 * time.Second
	)

	//check cache
	//Not for unit test in case of weekends
	if !common.Config.UnitTest {
		redisResp, redisErr := CacheInstance().Get(cacheKey).Result()
		if redisErr != nil {
			if redisErr == redis.Nil {
				log.Warn(ctx, "GetDayId | No result of %v in Redis, reading from API", cacheKey)
			} else {
				log.Error(ctx, "GetDayId | Error while reading from redis: %v", redisErr.Error())
			}
		} else {
			redisRespInt, _ := strconv.Atoi(redisResp)
			return int64(redisRespInt)
		}
	}

	var (
		currentMenu sea_dinner.Current
		currentID   int64
	)

	_, err := Client.R().
		SetHeader("Authorization", MakeToken(ctx, key)).
		SetResult(&currentMenu).
		Get(MakeURL(int(sea_dinner.URLType_URL_CURRENT), nil))

	if err != nil {
		log.Error(ctx, err.Error())
	}

	currentID = currentMenu.GetMenu().GetId()

	if !common.Config.UnitTest && currentMenu.GetMenu().GetPollStart() != fmt.Sprint(ConvertTimeStamp(time.Now().Unix()), "T04:30:00Z") {
		log.Warn(ctx, "GetDayId | Today's ID not found: %v", currentMenu.GetMenu().GetPollStart())
		currentID = 0
		expiry = 1800 * time.Second
	}

	//Short TTL if day is invalid
	//Might due to late menu update, so we have some room to get the correct data
	if err := CacheInstance().Set(cacheKey, currentID, expiry).Err(); err != nil {
		log.Error(ctx, "GetDayId | Error while writing to redis: %v", err.Error())
	} else {
		log.Info(ctx, "GetDayId | Successful | Written %v to redis", cacheKey)
	}
	return currentID
}
