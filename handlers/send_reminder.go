package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

//SendReminder Sends out daily reminder at 10.30 SGT on weekdays / working days
func SendReminder(ctx context.Context) {
	txn := processors.App.StartTransaction("send_reminder")
	defer txn.End()

	bot, err := tgbotapi.NewBotAPI(common.GetTGToken(ctx))
	if err != nil {
		log.Error(ctx, err.Error())
	}
	bot.Debug = true
	log.Info(ctx, "Authorized on account %s", bot.Self.UserName)

	res := BatchGetUsersChoice(ctx)
	log.Info(ctx, "SendReminder | size: %v", len(res))

	menu := MakeMenuNameMap(ctx)
	code := MakeMenuCodeMap(ctx)

	for _, r := range res {
		msg := tgbotapi.NewMessage(r.GetUserId(), "")

		var (
			mk     tgbotapi.InlineKeyboardMarkup
			out    [][]tgbotapi.InlineKeyboardButton
			rows   []tgbotapi.InlineKeyboardButton
			msgTxt string
		)

		if processors.IsSOW(time.Now()) {
			//Everyone except "MUTE" will receive weekly reminders
			msgTxt = "Good Morning! It's a brand new week with a brand new menu! Check it out at /menu 😋"
			randomBotton := tgbotapi.NewInlineKeyboardButtonData("🎲", "RAND")
			rows = append(rows, randomBotton)
			skipBotton := tgbotapi.NewInlineKeyboardButtonData("🙅", "-1")
			rows = append(rows, skipBotton)
			out = append(out, rows)
			mk.InlineKeyboard = out
			msg.ReplyMarkup = mk
		} else {
			//Only skips on non-mondays
			if r.GetUserChoice() == "-1" {
				log.Info(ctx, "SendReminder | skip -1 records | %v", r.GetUserId())
				continue
			}

			_, ok := menu[r.GetUserChoice()]
			if !ok {
				msgTxt = fmt.Sprintf("Good Morning. Your previous order %v is not available today! I will not proceed to order. Choose another dish from /menu 😃 /mute to shut me up 🫢 ", r.GetUserChoice())
			} else {
				if r.GetUserChoice() != "-1" {
					//If choice was updated after yesterdays' lunch time
					if r.GetMtime() > processors.GetPreviousDayLunchTime().Unix() {
						msgTxt = fmt.Sprintf("Good Morning. I will order %v %v today! If you changed your mind, you can choose from /menu 😋", code[r.GetUserChoice()], menu[r.GetUserChoice()])
						if r.GetUserChoice() == "RAND" {
							msgTxt = "Good Morning. I will order a random dish today! If you changed your mind, you can choose from /menu 😋"
						}
					} else {
						msgTxt = fmt.Sprintf("Good Morning. I will order %v %v again, just like yesterday! If you changed your mind, you can choose from /menu 😋", code[r.GetUserChoice()], menu[r.GetUserChoice()])
						if r.GetUserChoice() == "RAND" {
							msgTxt = "Good Morning. I will order a random dish again today! If you changed your mind, you can choose from /menu 😋"
						}
					}

					//If choice is already RAND, don't show RAND button again
					if r.GetUserChoice() != "RAND" {
						randomBotton := tgbotapi.NewInlineKeyboardButtonData("🎲", "RAND")
						rows = append(rows, randomBotton)
					}

					ignoreBotton := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v is good!", code[r.GetUserChoice()]), "SAME")
					rows = append(rows, ignoreBotton)
					skipBotton := tgbotapi.NewInlineKeyboardButtonData("🙅", "-1")
					rows = append(rows, skipBotton)
					out = append(out, rows)
					mk.InlineKeyboard = out
					msg.ReplyMarkup = mk
				}
			}
		}
		msg.Text = msgTxt
		if _, err := bot.Send(msg); err != nil {
			log.Error(ctx, err.Error())
		}
	}
	log.Info(ctx, "SendReminder | Success")
}
