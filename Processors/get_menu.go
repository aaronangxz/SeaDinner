package Processors

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/proto"
)

//GetMenu Calls Sea API, retrieves the current day's menu in realtime
func GetMenu(client resty.Client, key string) *sea_dinner.DinnerMenuArray {
	var (
		currentarr *sea_dinner.DinnerMenuArray
	)
	txn := App.StartTransaction("get_menu")
	defer txn.End()

	_, err := client.R().
		SetHeader("Authorization", MakeToken(key)).
		SetResult(&currentarr).
		Get(MakeURL(int(sea_dinner.URLType_URL_MENU), proto.Int64(GetDayId())))

	if err != nil {
		log.Println(err)
	}

	log.Printf("GetMenu | Query status of today's menu: %v", currentarr.GetStatus())
	return currentarr
}

//GetMenuUsingCache Calls Sea API, retrieves the current day's menu. Supports cache with TTL of 60 mins
func GetMenuUsingCache(client resty.Client, key string) *sea_dinner.DinnerMenuArray {
	var (
		cacheKey   = fmt.Sprint(Common.MENU_CACHE_KEY_PREFIX, ConvertTimeStamp(time.Now().Unix()))
		expiry     = 7200 * time.Second
		currentarr *sea_dinner.DinnerMenuArray
	)
	txn := App.StartTransaction("get_menu_using_cache")
	defer txn.End()

	//check cache
	val, redisErr := RedisClient.Get(cacheKey).Result()
	if redisErr != nil {
		if redisErr == redis.Nil {
			log.Printf("GetMenuUsingCache | No result of %v in Redis, reading from API", cacheKey)
		} else {
			log.Printf("GetMenuUsingCache | Error while reading from redis: %v", redisErr.Error())
		}
	} else {
		redisResp := &sea_dinner.DinnerMenuArray{}
		err := json.Unmarshal([]byte(val), &redisResp)
		if err != nil {
			log.Printf("GetMenuUsingCache | Fail to unmarshal Redis value of key %v : %v, reading from API", cacheKey, err)
		} else {
			log.Printf("GetMenuUsingCache | Successful | Cached %v", cacheKey)
			return redisResp
		}
	}

	currentarr = GetMenu(Client, key)

	//set back into cache
	data, err := json.Marshal(currentarr)
	if err != nil {
		log.Printf("GetMenuUsingCache | Failed to marshal JSON results: %v\n", err.Error())
	}

	if err := RedisClient.Set(cacheKey, data, expiry).Err(); err != nil {
		log.Printf("GetMenuUsingCache | Error while writing to redis: %v", err.Error())
	} else {
		log.Printf("GetMenuUsingCache | Successful | Written %v to redis", cacheKey)
	}

	log.Printf("GetMenuUsingCache | Query status of today's menu: %v", currentarr.GetStatus())
	return currentarr
}

//OutputMenuWithButton Sends menu and callback buttons
func OutputMenuWithButton(key string, id int64) ([]string, []tgbotapi.InlineKeyboardMarkup) {
	var (
		texts           []string
		out             []tgbotapi.InlineKeyboardMarkup
		dayText         string = "today"
		skipFillButtons bool
	)
	txn := App.StartTransaction("output_menu_with_button")
	defer txn.End()

	m := GetMenuUsingCache(Client, key)

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
		texts = append(texts, fmt.Sprintf(Common.Config.Prefix.UrlPrefix+"%v\n%v(%v) %v\nAvailable: %v", d.GetImageUrl(), d.GetCode(), d.GetId(), d.GetName(), d.GetQuota()))

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
		randomBotton := tgbotapi.NewInlineKeyboardButtonData("🎲", "RAND")
		rows = append(rows, randomBotton)
		skipBotton := tgbotapi.NewInlineKeyboardButtonData("🙅", "-1")
		rows = append(rows, skipBotton)
		out = append(out, tgbotapi.NewInlineKeyboardMarkup(rows))
	}

	return texts, out
}

func MenuRefresher() {
	ticker := time.NewTicker(time.Duration(Common.Config.Runtime.MenuRefreshIntervalSeconds) * time.Second)

	for range ticker.C {
		func() {
			key := os.Getenv("TOKEN")
			log.Println("MenuRefresher | Comparing Live and Cached menu.")

			liveMenu := GetMenu(Client, key)
			cacheMenu := GetMenuUsingCache(Client, key)

			if !CompareSliceStruct(liveMenu.GetFood(), cacheMenu.GetFood()) {
				log.Println("MenuRefresher | Live and Cached menu are inconsistent.")
				cacheKey := fmt.Sprint(Common.MENU_CACHE_KEY_PREFIX, ConvertTimeStamp(time.Now().Unix()))
				expiry := 7200 * time.Second

				data, err := json.Marshal(liveMenu)
				if err != nil {
					log.Printf("MenuRefresher | Failed to marshal JSON results: %v\n", err.Error())
				}

				//Use live menu as the source of truth
				if err := RedisClient.Set(cacheKey, data, expiry).Err(); err != nil {
					log.Printf("MenuRefresher | Error while writing to redis: %v", err.Error())
				} else {
					log.Printf("MenuRefresher | Successful | Written %v to redis", cacheKey)
				}
			}
		}()
	}
}
