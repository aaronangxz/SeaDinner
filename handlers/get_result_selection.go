package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetResultSelection() (string, []tgbotapi.InlineKeyboardMarkup) {
	return "Pick a time range:", []tgbotapi.InlineKeyboardMarkup{
		tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Week", "WEEKRESULT"),
				tgbotapi.NewInlineKeyboardButtonData("Month", "MONTHRESULT"),
				//Yearly result will be too long to send via TG. Need to figure out how to paginate it.
			}),
	}
}
