package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetResultSelection() (string, []tgbotapi.InlineKeyboardMarkup) {
	var (
		out []tgbotapi.InlineKeyboardMarkup
	)

	var rows []tgbotapi.InlineKeyboardButton
	text := "Pick a time range:"
	weekButton := tgbotapi.NewInlineKeyboardButtonData("Week", "WEEKRESULT")
	rows = append(rows, weekButton)
	monthButton := tgbotapi.NewInlineKeyboardButtonData("Month", "MONTHRESULT")
	rows = append(rows, monthButton)
	//Yearly result will be too long to send via TG. Need to figure out how to paginate it.
	//yearButton := tgbotapi.NewInlineKeyboardButtonData("Year", "YEARRESULT")
	//rows = append(rows, yearButton)
	out = append(out, tgbotapi.NewInlineKeyboardMarkup(rows))

	return text, out
}
