package processors

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"os"
	"time"

	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/proto"
)

//GetMenu Calls Sea API, retrieves the current day's menu in realtime
func GetMenu(ctx context.Context, client resty.Client, key string) *sea_dinner.DinnerMenuArray {
	var (
		currentArr *sea_dinner.DinnerMenuArray
	)
	txn := App.StartTransaction("get_menu")
	defer txn.End()

	_, err := client.R().
		SetHeader("Authorization", MakeToken(ctx, key)).
		SetResult(&currentArr).
		Get(MakeURL(int(sea_dinner.URLType_URL_MENU), proto.Int64(GetDayID(ctx))))

	if err != nil {
		log.Error(ctx, err.Error())
	}
	log.Info(ctx, "GetMenu | Query status of today's menu: %v", currentArr.GetStatus())
	return currentArr
}

//GetMenuUsingCache Calls Sea API, retrieves the current day's menu. Supports cache
func GetMenuUsingCache(ctx context.Context, key string) *sea_dinner.DinnerMenuArray {
	var (
		cacheKey   = fmt.Sprint(common.MENU_CACHE_KEY_PREFIX, ConvertTimeStamp(time.Now().Unix()))
		currentArr *sea_dinner.DinnerMenuArray
	)
	txn := App.StartTransaction("get_menu_using_cache")
	defer txn.End()

	//check cache
	val, redisErr := RedisClient.Get(cacheKey).Result()
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

	currentArr = GetMenu(ctx, Client, key)

	//set back into cache
	data, err := json.Marshal(currentArr)
	if err != nil {
		log.Error(ctx, "GetMenuUsingCache | Failed to marshal JSON results: %v\n", err.Error())
	}

	if err := RedisClient.Set(cacheKey, data, 0).Err(); err != nil {
		log.Error(ctx, "GetMenuUsingCache | Error while writing to redis: %v", err.Error())
	} else {
		log.Info(ctx, "GetMenuUsingCache | Successful | Written %v to redis", cacheKey)
	}
	log.Info(ctx, "GetMenuUsingCache | Query status of today's menu: %v", currentArr.GetStatus())
	return currentArr
}

//OutputMenuWithButton Sends menu and callback buttons
func OutputMenuWithButton(ctx context.Context, key string) ([]string, []tgbotapi.InlineKeyboardMarkup) {
	var (
		texts           []string
		out             []tgbotapi.InlineKeyboardMarkup
		dayText         = "today"
		skipFillButtons bool
	)
	txn := App.StartTransaction("output_menu_with_button")
	defer txn.End()

	if !IsWeekDay() {
		texts = append(texts, "We are done for this week! You can order again next week 😀")
		return texts, out
	}

	m := GetMenuUsingCache(ctx, key)

	if m.Status == nil {
		texts = append(texts, "There is no dinner order today! 😕")
		return texts, out
	}

	tz, _ := time.LoadLocation(TimeZone)
	if time.Now().In(tz).Unix() > GetLunchTime().Unix() {
		if IsNotEOW(time.Now().In(tz)) {
			dayText = "tomorrow"
		} else {
			skipFillButtons = true
		}
	}

	for _, d := range m.GetFood() {
		texts = append(texts, fmt.Sprintf(common.Config.Prefix.URLPrefix+"%v\n%v(%v) %v\nAvailable: %v", d.GetImageUrl(), d.GetCode(), d.GetId(), d.GetName(), d.GetQuota()))

		if !skipFillButtons {
			var buttons []tgbotapi.InlineKeyboardButton
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Snatch %v %v", d.GetCode(), dayText), fmt.Sprint(d.GetId())))
			out = append(out, tgbotapi.NewInlineKeyboardMarkup(buttons))
		}
	}

	//Follows the same conditions
	if !skipFillButtons {
		var rows []tgbotapi.InlineKeyboardButton
		texts = append(texts, fmt.Sprintf("Other Options👇🏻\n\n🎲 If you're feeling lucky\n🙅 If you don't need it / not coming to office %v", dayText))
		randomButton := tgbotapi.NewInlineKeyboardButtonData("🎲", "RAND")
		rows = append(rows, randomButton)
		skipButton := tgbotapi.NewInlineKeyboardButtonData("🙅", "-1")
		rows = append(rows, skipButton)
		out = append(out, tgbotapi.NewInlineKeyboardMarkup(rows))
	}
	log.Info(ctx, "OutputMenuWithButton | Success")
	return texts, out
}

//MenuRefresher Periodically refreshes cached menu with the updated live menu
func MenuRefresher(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(common.Config.Runtime.MenuRefreshIntervalSeconds) * time.Second)

	for range ticker.C {
		func() {
			if !IsActiveDay() {
				log.Warn(ctx, "MenuRefresher | Inactive day | Resumes check tomorrow.")
				time.Sleep(time.Duration(GetEOD().Unix()-time.Now().Unix()) * time.Second)
				return
			}
			key := os.Getenv("TOKEN")
			log.Info(ctx, "MenuRefresher | Comparing Live and Cached menu.")

			liveMenu := GetMenu(ctx, Client, key)
			cacheMenu := GetMenuUsingCache(ctx, key)

			if !CompareSliceStruct(ctx, liveMenu.GetFood(), cacheMenu.GetFood()) {
				log.Warn(ctx, "MenuRefresher | Live and Cached menu are inconsistent.")
				cacheKey := fmt.Sprint(common.MENU_CACHE_KEY_PREFIX, ConvertTimeStamp(time.Now().Unix()))

				data, err := json.Marshal(liveMenu)
				if err != nil {
					log.Error(ctx, "MenuRefresher | Failed to marshal JSON results: %v\n", err.Error())
				}

				//Use live menu as the source of truth
				if err := RedisClient.Set(cacheKey, data, 0).Err(); err != nil {
					log.Error(ctx, "MenuRefresher | Error while writing to redis: %v", err.Error())
				} else {
					log.Info(ctx, "MenuRefresher | Successful | Written %v to redis", cacheKey)
				}
			}
			log.Info(ctx, "MenuRefresher | Live and Cached menu are consistent.")
		}()
	}
}
