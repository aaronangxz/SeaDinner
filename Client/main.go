package main

import (
	"log"
	"os"
	"time"

	"github.com/aaronangxz/SeaDinner/Bot"
	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/Processors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	startListenKey   = false
	startListenChope = false
)

func main() {
	Processors.LoadEnv()
	Processors.Init()
	Processors.InitClient()

	bot, err := tgbotapi.NewBotAPI(Common.GetTGToken())
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
			msg.ParseMode = "MARKDOWN"
			msg.Text, _ = Bot.CallbackQueryHandler(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery)
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
			continue
		}

		//Stop responding from 12.29pm to 12.31pm or until dinner order has started (For occasional weird order timings)
		if time.Now().Unix() >= Processors.GetLunchTime().Unix()-60 &&
			(time.Now().Unix() <= Processors.GetLunchTime().Unix()+60 && !Processors.IsPollStart()) {
			if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Omw to order, wait for my good news! 🏃")); err != nil {
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
				msg, _ := Bot.UpdateKey(update.Message.Chat.ID, update.Message.Text)
				if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg)); err != nil {
					log.Println(err)
				}
				startListenKey = false
				continue
			} else if startListenChope {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				ok := false
				msg.Text, ok = Bot.GetChope(update.Message.Chat.ID, update.Message.Text)
				if !ok {
					if _, err := bot.Send(msg); err != nil {
						log.Println(err)
					}
					continue
				}
				//Capture chope
				msg.ParseMode = "MARKDOWN"
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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		// Extract the command from the Message.
		switch update.Message.Command() {
		case "start":
			s, ok := Bot.CheckKey(update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				msg.Text = "Hello! " + update.Message.Chat.UserName
			}
		case "menu":
			s, ok := Bot.CheckKey(update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				txt, mp := Processors.OutputMenuWithButton(Bot.GetKey(update.Message.Chat.ID), update.Message.Chat.ID)
				for i, r := range txt {
					msg.Text = r
					if len(mp) > 0 {
						msg.ReplyMarkup = mp[i]
					}
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
				}
				continue
			}
		case "help":
			msg.Text = Bot.MakeHelpResponse()
			msg.ParseMode = "MARKDOWN"
		case "key":
			msg.Text, _ = Bot.CheckKey(update.Message.Chat.ID)
		case "newkey":
			msg.Text = "What's your key? \nGo to https://dinner.sea.com/accounts/token, copy the Key under Generate Auth Token and paste it here:"
			startListenKey = true
		case "status":
			s, ok := Bot.CheckKey(update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				msg.Text = Bot.ListWeeklyResultByUserId(update.Message.Chat.ID)
				msg.ParseMode = "HTML"
			}
		case "chope":
			msg.Text = "This command is deprecated. Choose from /menu instead!😋"
		case "choice":
			s, ok := Bot.CheckKey(update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				msg.Text, _ = Bot.CheckChope(update.Message.Chat.ID)
			}
		case "reminder":
			//Backdoor for test env
			if os.Getenv("TEST_DEPLOY") == "TRUE" || Common.Config.Adhoc {
				Bot.SendReminder()
			}
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
