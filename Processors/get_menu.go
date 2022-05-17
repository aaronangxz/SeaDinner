package Processors

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetMenu(client resty.Client, ID int, key string) DinnerMenuArr {
	var currentarr DinnerMenuArr

	log.Println("key:", key)
	log.Println("header:", MakeToken(key))
	log.Println("url:", MakeURL(URL_MENU, &ID))

	if ID == 0 {
		log.Println("GetMenu | Invalid id:", ID)
		return currentarr
	}

	_, err := client.R().
		SetHeader("Authorization", MakeToken(key)).
		SetResult(&currentarr).
		Get(MakeURL(URL_MENU, &ID))

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Query status of today's menu: %v\n\n", currentarr.Status)

	return currentarr
}

func OutputMenu(key string) string {
	var (
		output string
	)

	m := GetMenu(Client, GetDayId(key), key)

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
		texts []string
		out   []tgbotapi.InlineKeyboardMarkup
	)

	m := GetMenu(Client, GetDayId(key), key)

	if m.Status == nil {
		texts = append(texts, "There is no dinner order today! 😕")
		return texts, out
		//return []string{"There is no dinner order today! 😕"}, []tgbotapi.InlineKeyboardMarkup{tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("", "")})})
	}

	for _, d := range m.DinnerArr {
		texts = append(texts, fmt.Sprintf(Config.Prefix.UrlPrefix+"%v\n%v(%v) %v\nAvailable:%v/%v", d.ImageURL, d.Code, d.Id, d.Name, d.Ordered, d.Quota))
		var buttons []tgbotapi.InlineKeyboardButton
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Order %v", d.Code), fmt.Sprint(d.Id)))
		out = append(out, tgbotapi.NewInlineKeyboardMarkup(buttons))
	}
	return texts, out
}
