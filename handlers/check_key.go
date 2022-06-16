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
	"os"
	"time"
)

//CheckKey Checks if user's API key exists.
//Reads from cache first, then user_key_tab.
func CheckKey(ctx context.Context, id int64) (string, bool) {
	var (
		existingRecord *sea_dinner.UserKey
		cacheKey       = fmt.Sprint(common.USER_KEY_PREFIX, id)
		expiry         = 604800 * time.Second
	)
	txn := processors.App.StartTransaction("check_key")
	defer txn.End()

	if id <= 0 {
		log.Error(ctx, "UpdateKey | Id must be > 1.")
		return "", false
	}

	//Read from cache
	val, redisErr := processors.RedisClient.Get(cacheKey).Result()
	if redisErr != nil {
		if redisErr == redis.Nil {
			log.Warn(ctx, "CheckKey | No result of %v in Redis, reading from DB", cacheKey)
		} else {
			log.Error(ctx, "CheckKey | Error while reading from redis: %v", redisErr.Error())
		}
	} else {
		redisResp := &sea_dinner.UserKey{}
		err := json.Unmarshal([]byte(val), &redisResp)
		if err != nil {
			log.Error(ctx, "CheckKey | Fail to unmarshal Redis value of key %v : %v, reading from DB", cacheKey, err)
		} else {
			log.Info(ctx, "CheckKey | Successful | Cached %v", cacheKey)
			decrypt := processors.DecryptKey(redisResp.GetUserKey(), os.Getenv("AES_KEY"))
			log.Info(ctx, "CheckKey | Successful")
			return fmt.Sprintf("I have your key %v***** that you told me on %v! But I won't leak it 😀", decrypt[:5], processors.ConvertTimeStamp(redisResp.GetMtime())), true
		}
	}

	//Read from DB
	if err := processors.DB.Table(common.DB_USER_KEY_TAB).Where("user_id = ?", id).First(&existingRecord).Error; err != nil {
		//Write into new user set
		return "I don't have your key, let me know in /newkey 😊", false
	}
	//set back into cache
	data, err := json.Marshal(existingRecord)
	if err != nil {
		log.Error(ctx, "CheckKey | Failed to marshal JSON results: %v\n", err.Error())
	}

	if err := processors.RedisClient.Set(cacheKey, data, expiry).Err(); err != nil {
		log.Error(ctx, "CheckKey | Error while writing to redis: %v", err.Error())
	} else {
		log.Info(ctx, "CheckKey | Successful | Written %v to redis", cacheKey)
	}
	decrypt := processors.DecryptKey(existingRecord.GetUserKey(), os.Getenv("AES_KEY"))
	log.Info(ctx, "CheckKey | Successful")
	return fmt.Sprintf("I have your key %v***** that you told me on %v! But I won't leak it 😀", decrypt[:5], processors.ConvertTimeStamp(existingRecord.GetMtime())), true

}
