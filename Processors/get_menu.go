package Processors

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetMenu(client resty.Client, ID int, key string) DinnerMenuArr {
	var (
		cacheKey   = fmt.Sprint(MENU_CACHE_KEY_PREFIX, ConvertTimeStamp(time.Now().Unix()))
		expiry     = 3600 * time.Second
		currentarr DinnerMenuArr
	)

	if ID == 0 {
		log.Println("GetMenu | Invalid id:", ID)
		return currentarr
	}

	//check cache
	val, redisErr := RedisClient.Get(cacheKey).Result()
	if redisErr != nil {
		if redisErr == redis.Nil {
			log.Printf("GetMenu | No result of %v in Redis, reading from API", cacheKey)
		} else {
			log.Printf("GetMenu | Error while reading from redis: %v", redisErr.Error())
		}
	} else {
		redisResp := DinnerMenuArr{}
		err := json.Unmarshal([]byte(val), &redisResp)
		if err != nil {
			log.Printf("GetMenu | Fail to unmarshal Redis value of key %v : %v, reading from API", cacheKey, err)
		} else {
			log.Printf("GetMenu | Successful | Cached %v", cacheKey)
			return redisResp
		}
	}

	_, err := client.R().
		SetHeader("Authorization", MakeToken(key)).
		SetResult(&currentarr).
		Get(MakeURL(URL_MENU, &ID))

	if err != nil {
		log.Println(err)
	}

	//set back into cache
	data, err := json.Marshal(currentarr)
	if err != nil {
		log.Printf("GetMenu | Failed to marshal JSON results: %v\n", err.Error())
	}

	if err := RedisClient.Set(cacheKey, data, expiry).Err(); err != nil {
		log.Printf("GetMenu | Error while writing to redis: %v", err.Error())
	} else {
		log.Printf("GetMenu | Successful | Written %v to redis", cacheKey)
	}

	log.Printf("GetMenu | Query status of today's menu: %v", currentarr.GetStatus())
	return currentarr
}

func OutputMenu(key string) string {
	var (
		output string
	)

	m := GetMenu(Client, GetDayId(), key)

	if m.Status == nil {
		return "There is no dinner order today! 😕"
	}

	for _, d := range m.DinnerArr {
		output += fmt.Sprintf(Config.Prefix.UrlPrefix+"%v\nFood ID: %v\nName: %v\nQuota: %v\n\n",
			d.ImageURL, d.Id, d.Name, d.Quota)
	}
	return output
}

func OutputMenuWithButton(key string, id int64) ([]string, []tgbotapi.InlineKeyboardMarkup) {
	var (
		texts           []string
		out             []tgbotapi.InlineKeyboardMarkup
		buttonText      string = "Snatch %v today"
		skipFillButtons bool
	)

	m := GetMenu(Client, GetDayId(), key)

	if m.Status == nil {
		texts = append(texts, "There is no dinner order today! 😕")
		return texts, out
	}

	tz, _ := time.LoadLocation(TimeZone)
	if time.Now().In(tz).Unix() > GetLunchTime().Unix() {
		if IsNotEOW(time.Now().In(tz)) {
			buttonText = "Snatch %v tomorrow"
		} else {
			skipFillButtons = true
		}
	}

	for _, d := range m.DinnerArr {
		texts = append(texts, fmt.Sprintf(Config.Prefix.UrlPrefix+"%v\n%v(%v) %v\nAvailable: %v", d.ImageURL, d.Code, d.Id, d.Name, d.Quota))

		if !skipFillButtons {
			var buttons []tgbotapi.InlineKeyboardButton
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf(buttonText, d.Code), fmt.Sprint(d.Id)))
			out = append(out, tgbotapi.NewInlineKeyboardMarkup(buttons))
		}
	}

	//Follows the same conditions
	if !skipFillButtons {
		//Append for random
		randomBotton := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("I'm feeling lucky!", "RAND")}
		out = append(out, tgbotapi.NewInlineKeyboardMarkup(randomBotton))

		//Append for order skipping
		skipBotton := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Nah I'm good.", "-1")}
		out = append(out, tgbotapi.NewInlineKeyboardMarkup(skipBotton))
	}

	return texts, out
}
