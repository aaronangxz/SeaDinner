package processors

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/go-redis/redis"
	"time"
)

//GetMenuUsingCache Calls Sea API, retrieves the current day's menu. Supports cache
func GetMenuUsingCache(ctx context.Context, key string) *sea_dinner.DinnerMenuArray {
	var (
		cacheKey   = fmt.Sprint(common.MENU_CACHE_KEY_PREFIX, ConvertTimeStamp(time.Now().Unix()))
		currentArr *sea_dinner.DinnerMenuArray
	)
	txn := App.StartTransaction("get_menu_using_cache")
	defer txn.End()

	//check cache
	val, redisErr := CacheInstance().Get(cacheKey).Result()
	if redisErr != nil {
		if redisErr == redis.Nil {
			log.Warn(ctx, "GetMenuUsingCache | No result of %v in Redis, reading from API", cacheKey)
		} else {
			log.Error(ctx, "GetMenuUsingCache | Error while reading from redis: %v", redisErr.Error())
		}
	} else {
		redisResp := &sea_dinner.DinnerMenuArray{}
		err := json.Unmarshal([]byte(val), &redisResp)
		if err != nil {
			log.Warn(ctx, "GetMenuUsingCache | Fail to unmarshal Redis value of key %v : %v, reading from API", cacheKey, err)
		} else {
			log.Info(ctx, "GetMenuUsingCache | Successful | Cached %v", cacheKey)
			return redisResp
		}
	}

	currentArr = GetMenu(ctx, key)

	//set back into cache
	data, err := json.Marshal(currentArr)
	if err != nil {
		log.Error(ctx, "GetMenuUsingCache | Failed to marshal JSON results: %v\n", err.Error())
	}

	if err := CacheInstance().Set(cacheKey, data, 0).Err(); err != nil {
		log.Error(ctx, "GetMenuUsingCache | Error while writing to redis: %v", err.Error())
	} else {
		log.Info(ctx, "GetMenuUsingCache | Successful | Written %v to redis", cacheKey)
	}
	log.Info(ctx, "GetMenuUsingCache | Query status of today's menu: %v", currentArr.GetStatus())
	return currentArr
}
