package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/go-redis/redis"
	"google.golang.org/protobuf/proto"
	"math"
	"time"
)

//UpdateChope Updates the current food choice made by user.
//With basic parameter verifications
//Supports Button Callbacks
func UpdateChope(ctx context.Context, id int64, s string) (string, bool) {
	var (
		existingRecord sea_dinner.UserChoice
		r              = &sea_dinner.UserChoice{
			UserId:     proto.Int64(id),
			UserChoice: proto.String(s),
			Ctime:      proto.Int64(time.Now().Unix()),
			Mtime:      proto.Int64(time.Now().Unix()),
		}
		key = fmt.Sprint(common.USER_CHOICE_PREFIX, r.GetUserId())
	)
	txn := processors.App.StartTransaction("get_chope")
	defer txn.End()

	if id <= 0 {
		log.Error(ctx, "Id must be > 1.")
		return "", false
	}

	//When it is Friday after 12.30pm, we don't accept any orders (except -1) because we don't know next week's menu yet
	if !processors.IsNotEOW(time.Now()) && time.Now().Unix() > processors.GetLunchTime().Unix() && s != "-1" {
		return "We are done for this week! You can tell me your order again next week 😀", false
	}

	if processors.IsNotNumber(s) {
		//RAND is passed from CallBack
		if s != "RAND" && s != "SAME" {
			log.Error(ctx, "Selection contains illegal character | selection: %v", s)
			return "Are you sure that is a valid FoodID? Tell me another one. 😟", false
		}
	}

	menu := MakeMenuNameMap(ctx)

	if s == "SAME" {
		//Set back to DB using cache data (user_choice:<id>)
		//To handle Morning Reminder callback
		val, redisErr := processors.CacheInstance().Get(key).Result()
		if redisErr != nil {
			if redisErr == redis.Nil {
				log.Info(ctx, "GetChope | No result of %v in Redis, reading from API", key)
			} else {
				log.Error(ctx, "GetChope | Error while reading from redis: %v", redisErr.Error())
			}
			return "The selection has expired, you can choose from /menu again 😀", true
		}

		if val == "" {
			log.Error(ctx, "GetChope | empty in redis: %v", key)
			return "The selection has expired, you can choose from /menu again 😀", true
		}

		if err := processors.DbInstance().Exec("UPDATE user_choice_tab SET user_choice = ?, mtime = ? WHERE user_id = ?", val, time.Now().Unix(), id).Error; err != nil {
			log.Error(ctx, "Failed to update DB | %v", err.Error())
			return err.Error(), false
		}

		if val == "-1" {
			return "Okay got it. I will not order anything for you instead.😀", true
		}

		if val == "RAND" {
			return "Okay got it. I will give you a surprise 😙", true
		}
		return fmt.Sprintf("Okay got it! I will order %v 😙", menu[val]), true
	}

	_, ok := menu[s]
	if !ok {
		log.Error(ctx, "Selection is invalid | selection: %v", s)
		return "This dish is not available today. Tell me another one.😟", false
	}

	if err := processors.DbInstance().Raw("SELECT * FROM user_choice_tab WHERE user_id = ?", id).Scan(&existingRecord).Error; err != nil {
		log.Error(ctx, "GetChope | %v", err.Error())
		return err.Error(), false
	}
	if existingRecord.UserId == nil {
		if err := processors.DbInstance().Table(common.DB_USER_CHOICE_TAB).Create(&r).Error; err != nil {
			log.Error(ctx, "Failed to update DB | %v", err.Error())
			return err.Error(), false
		}

		//To stop ordering
		if s == "-1" {
			return fmt.Sprintf("Okay got it. I will order %v for you and stop sending reminders in the morning for the rest of the week.😀", menu[s]), true
		}

		if s == "RAND" {
			return "Okay got it. I will give you a surprise😙", true
		}

		//Orders placed before lunch time
		if time.Now().Unix() < processors.GetLunchTime().Unix() {
			return fmt.Sprintf("Okay got it. I will order %v for you today 😙", menu[s]), true
		}

		return fmt.Sprintf("Okay got it. I will order %v for you tomorrow 😙", menu[s]), true
	}
	//Update key if user_id exists
	if err := processors.DbInstance().Exec("UPDATE user_choice_tab SET user_choice = ?, mtime = ? WHERE user_id = ?", s, time.Now().Unix(), id).Error; err != nil {
		log.Error(ctx, "Failed to update DB | %v", err.Error())
		return err.Error(), false
	}

	//To stop ordering
	if s == "-1" {
		return fmt.Sprintf("Okay got it. I will order %v for you and stop sending reminders in the morning for the rest of the week.😀", menu[s]), true
	}

	if s == "RAND" {
		return "Okay got it. I will give you a surprise 😙", true
	}

	//Orders placed before lunch time
	if time.Now().Unix() < processors.GetLunchTime().Unix() {
		//Set into cache for Morning reminder callback. TTL is always until 12.30
		//Minimum TTL is 1 second
		expiry := time.Duration(math.Max(1, float64(processors.GetLunchTime().Unix()-time.Now().Unix())))
		if err := processors.CacheInstance().Set(key, s, expiry).Err(); err != nil {
			log.Error(ctx, "GetChope | Error while writing to redis: %v", err.Error())
		} else {
			log.Info(ctx, "GetChope | Successful | Written %v to redis", key)
		}
		return fmt.Sprintf("Okay got it. I will order %v for you today 😙", menu[s]), true
	}
	return fmt.Sprintf("Okay got it. I will order %v for you tomorrow 😙", menu[s]), true
}
