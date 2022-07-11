package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetExitSelection() (string, []tgbotapi.InlineKeyboardMarkup) {
	return "What happened? Where are you going 😥", []tgbotapi.InlineKeyboardMarkup{
		tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Opting Out", "ATTEMPTOPTOUT"),
				tgbotapi.NewInlineKeyboardButtonData("Leaving Shopee", "ATTEMPTRESIGN"),
			}),
	}
}
