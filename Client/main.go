package main

import (
	"log"
	"os"
	"time"

	"github.com/aaronangxz/SeaDinner/Bot"
	"github.com/aaronangxz/SeaDinner/Processors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	startListenKey   = false
	startListenChope = false
	Id               int64
)

func main() {
	Processors.Init()
	Processors.LoadEnv()
	if Processors.Config.Adhoc {
		Processors.ConnectTestMySQL()
	} else {
		Processors.ConnectMySQL()
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// var b = tgbotapi.NewInlineKeyboardMarkup(
	// 	[]tgbotapi.InlineKeyboardButton{
	// 		tgbotapi.NewInlineKeyboardButtonData("Select me", "sss"),
	// 		tgbotapi.NewInlineKeyboardButtonData("Select me", "sss"),
	// 	},
	// )

	// var a []tgbotapi.InlineKeyboardButton
	// a = append(a, tgbotapi.NewInlineKeyboardButtonData("Select This", "sss"))
	// a = append(a, tgbotapi.NewInlineKeyboardButtonData("Select That", "sss"))
	// mk := tgbotapi.NewInlineKeyboardMarkup(a)

	for update := range updates {
		if update.CallbackQuery != nil {
			msg := tgbotapi.NewMessage(Id, "")
			msg.Text, _ = Bot.CallbackQueryHandler(Id, update.CallbackQuery)
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
			continue
		}

		if update.Message.Chat.ID != 0 {
			Id = update.Message.Chat.ID
		}

		if time.Now().Unix() >= Processors.GetLunchTime().Unix()-60 && time.Now().Unix() <= Processors.GetLunchTime().Unix()+210 {
			if _, err := bot.Send(tgbotapi.NewMessage(Id, "Omw to order, wait for my good news! 🏃")); err != nil {
				log.Println(err)
			}
			continue
		}

		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			if startListenKey {
				//Capture key
				msg, _ := Bot.UpdateKey(Id, update.Message.Text)
				if _, err := bot.Send(tgbotapi.NewMessage(Id, msg)); err != nil {
					log.Println(err)
				}
				startListenKey = false
				continue
			} else if startListenChope {
				msg := tgbotapi.NewMessage(Id, "")
				ok := false
				msg.Text, ok = Bot.GetChope(Id, update.Message.Text)

				if !ok {
					if _, err := bot.Send(msg); err != nil {
						log.Println(err)
					}
					continue
				}
				//Capture chope
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
				}
				startListenChope = false
				continue
			} else {
				continue
			}
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(Id, "")
		// Extract the command from the Message.
		switch update.Message.Command() {
		case "start":
			s, ok := Bot.CheckKey(Id)
			if !ok {
				msg.Text = s
			} else {
				msg.Text = "Hello! " + update.Message.Chat.UserName
			}
		case "menu":
			s, ok := Bot.CheckKey(Id)
			if !ok {
				msg.Text = s
			} else {
				//msg.Text = Processors.OutputMenu(Bot.GetKey(Id))
				txt, mp := Processors.OutputMenuWithButton(Bot.GetKey(Id), Id)
				for i, r := range txt {
					msg.Text = r
					msg.ReplyMarkup = mp[i]
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
				}
				continue
			}
		case "help":
			msg.Text = "Check the commands."
		case "key":
			msg.Text, _ = Bot.CheckKey(Id)
		case "newkey":
			msg.Text = "What's your key? \nGo to https://dinner.sea.com/accounts/token, copy the Key under Generate Auth Token and paste it here:"
			startListenKey = true
		case "status":
			msg.Text = Bot.GetLatestResultByUserId(Id)
		case "chope":
			msg.Text = "What do you want to order? Tell me the Food ID 😋"
			//msg.ReplyMarkup = numericKeyboard
			startListenChope = true
		case "choice":
			msg.Text, _ = Bot.CheckChope(Id)
		case "ret":
			return
		default:
			msg.Text = "I don't understand this command :("
		}
		if msg.Text != "" {
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}
