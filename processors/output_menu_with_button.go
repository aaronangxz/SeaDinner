package processors

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

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
