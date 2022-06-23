package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"reflect"
	"testing"
)

func TestGetResultSelection(t *testing.T) {
	expectedText := "Pick a time range:"
	expectedKb := []tgbotapi.InlineKeyboardMarkup{
		tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Week", "WEEKRESULT"),
				tgbotapi.NewInlineKeyboardButtonData("Month", "MONTHRESULT"),
			}),
	}

	tests := []struct {
		name  string
		want  string
		want1 []tgbotapi.InlineKeyboardMarkup
	}{
		{
			"HappyCase",
			expectedText,
			expectedKb,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetResultSelection()
			if got != tt.want {
				t.Errorf("GetResultSelection() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetResultSelection() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
