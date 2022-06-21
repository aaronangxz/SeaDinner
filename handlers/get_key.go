package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/go-redis/redis"
	"time"
)

//GetKey Retrieves user's API key with user_id.
//Reads from cache first, then user_key_tab.
func GetKey(ctx context.Context, id int64) string {
	var (
		existingRecord *sea_dinner.UserKey
		cacheKey       = fmt.Sprint(common.USER_KEY_PREFIX, id)
		expiry         = 604800 * time.Second
	)
	txn := processors.App.StartTransaction("get_key")
	defer txn.End()

	if id <= 0 {
		log.Error(ctx, "UpdateKey | Id must be > 1.")
		return ""
	}

	//Read from cache
	val, redisErr := processors.CacheInstance().Get(cacheKey).Result()
	if redisErr != nil {
		if redisErr == redis.Nil {
			log.Warn(ctx, "GetKey | No result of %v in Redis, reading from DB", cacheKey)
		} else {
			log.Error(ctx, "GetKey | Error while reading from redis: %v", redisErr.Error())
		}
	} else {
		redisResp := &sea_dinner.UserKey{}
		err := json.Unmarshal([]byte(val), &redisResp)
		if err != nil {
			log.Error(ctx, "GetKey | Fail to unmarshal Redis value of key %v : %v, reading from DB", cacheKey, err)
		} else {
			log.Info(ctx, "GetKey | Successful | Cached %v", cacheKey)
			return redisResp.GetUserKey()
		}
	}

	//Read from DB
	if err := processors.DbInstance().Table(common.DB_USER_KEY_TAB).Where("user_id = ?", id).First(&existingRecord).Error; err != nil {
		log.Error(ctx, "GetKey | Failed to find record | %v", err.Error())
		return ""
	}

	//set back into cache
	data, err := json.Marshal(existingRecord)
	if err != nil {
		log.Error(ctx, "GetKey | Failed to marshal JSON results: %v\n", err.Error())
	}

	if err := processors.CacheInstance().Set(cacheKey, data, expiry).Err(); err != nil {
		log.Error(ctx, "GetKey | Error while writing to redis: %v", err.Error())
	} else {
		log.Info(ctx, "GetKey | Successful | Written %v to redis", cacheKey)
	}
	log.Info(ctx, "GetKey | Successful.")
	return existingRecord.GetUserKey()
}
