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
	val, redisErr := processors.CacheInstance().Get(cacheKey).Result()
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
	if err := processors.DbInstance().Table(common.DB_USER_KEY_TAB).Where("user_id = ?", id).First(&existingRecord).Error; err != nil {
		//Remove the previous record of this user
		val, _, redisErr := processors.CacheInstance().SScan(common.POTENTIAL_USER_SET, 0, fmt.Sprint("*", id, "*"), 1000).Result()
		if redisErr != nil {
			if redisErr == redis.Nil {
				log.Warn(ctx, "CheckKey | No result of *%v* pattern in Redis", id)
			} else {
				log.Error(ctx, "CheckKey | Error while reading previous potential_user set from redis: %v", redisErr.Error())
			}
		} else {
			//Remove from potential_user Set
			for _, r := range val {
				if err := processors.CacheInstance().SRem(common.POTENTIAL_USER_SET, r).Err(); err != nil {
					log.Error(ctx, "CheckKey | Error while deleting previous potential_user set from redis: %v", err.Error())
				} else {
					log.Info(ctx, "CheckKey | Successful | Removed %v from potential_user set", r)
				}
			}
		}

		//Write the new record into potential_user set
		toWrite := fmt.Sprint(id, ":", time.Now().Unix())
		if err := processors.CacheInstance().SAdd(common.POTENTIAL_USER_SET, toWrite).Err(); err != nil {
			log.Error(ctx, "CheckKey | Error while writing potential_user set to redis: %v", err.Error())
		} else {
			log.Info(ctx, "CheckKey | Successful | Written %v to potential_user set", toWrite)
		}
		return "I don't have your key, let me know in /newkey 😊", false
	}

	//user status check
	switch existingRecord.GetStatus() {
	case int64(sea_dinner.UserStatus_USER_STATUS_INACTIVE):
		if err := UpdateUserStatus(ctx, existingRecord.GetUserId(), int64(sea_dinner.UserStatus_USER_STATUS_ACTIVE)); err != nil {
			return "Unexpected Error. Try again!", false
		}
		return "Welcome back! It's good to see you here again 😊", false
	case int64(sea_dinner.UserStatus_USER_STATUS_RESIGNED):
		if err := UpdateUserStatus(ctx, existingRecord.GetUserId(), int64(sea_dinner.UserStatus_USER_STATUS_ACTIVE)); err != nil {
			return "Unexpected Error. Try again!", false
		}
		return "Welcome back! 😊 Double check your key /key in case it doesn't match our previous records.", false
	}

	//set back into cache
	data, err := json.Marshal(existingRecord)
	if err != nil {
		log.Error(ctx, "CheckKey | Failed to marshal JSON results: %v\n", err.Error())
	}

	if err := processors.CacheInstance().Set(cacheKey, data, expiry).Err(); err != nil {
		log.Error(ctx, "CheckKey | Error while writing to redis: %v", err.Error())
	} else {
		log.Info(ctx, "CheckKey | Successful | Written %v to redis", cacheKey)
	}
	decrypt := processors.DecryptKey(existingRecord.GetUserKey(), os.Getenv("AES_KEY"))
	log.Info(ctx, "CheckKey | Successful")
	return fmt.Sprintf("I have your key %v***** that you told me on %v! But I won't leak it 😀", decrypt[:5], processors.ConvertTimeStamp(existingRecord.GetMtime())), true

}
