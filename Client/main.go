package main

import (
	"log"
	"math"
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
	userId           int64
)

func main() {
	Processors.Init()
	Processors.LoadEnv()
	Processors.ConnectDataBase()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			if update.MyChatMember != nil {
				Id = update.MyChatMember.Chat.ID
				userId = update.Message.From.ID
				log.Println("Id,userId:", Id, userId)
				if _, err := bot.Send(tgbotapi.NewMessage(Id, "Hello! Pm me to chat 😉")); err != nil {
					log.Println(err)
				}
			}
		} else {
			if update.Message.Chat.ID != 0 {
				if update.Message.From != nil {
					Id = update.Message.Chat.ID
					userId = update.Message.From.ID
					log.Println("Id,userId:", Id, userId)
				} else {
					Id = update.Message.Chat.ID
					log.Println("Id:", Id)
				}
			}
		}

		log.Println(time.Now().Unix(), Processors.GetLunchTime().Unix()-60)
		if time.Now().Unix() >= Processors.GetLunchTime().Unix()-60 && time.Now().Unix() <= Processors.GetLunchTime().Unix()+5 {
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
				//Capture chope
				if _, err := bot.Send(tgbotapi.NewMessage(Id,
					Bot.GetChope(Id, update.Message.Text))); err != nil {
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
		msg := tgbotapi.NewMessage(int64(math.Min(float64(Id), float64(userId))), "")
		name := ""
		if Id < 0 {
			name = update.Message.From.FirstName + " "
		}
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
				msg.Text = Processors.OutputMenu(Bot.GetKey(Id))
			}
		case "help":
			msg.Text = "Check the commands."
		case "key":
			msg.Text, _ = Bot.CheckKey(Id)
		case "newkey":
			msg.Text = "What's your key? \nGo to https://dinner.sea.com/accounts/token, copy the Key under Generate Auth Token and paste it here:"
			startListenKey = true
		case "status":
			msg.Text = name + Bot.GetLatestResultByUserId(int64(math.Max(float64(Id), float64(userId))))
			// msg.Entities = append(msg.Entities, tgbotapi.MessageEntity{
			// 	Type: "bold",
			// })
			log.Println("msg ID", update.Message.MessageID)
		case "chope":
			msg.Text = name + "What do you want to order? Tell me the Food ID 😋"
			startListenChope = true
		case "choice":
			msg.Text, _ = Bot.CheckChope(int64(math.Max(float64(Id), float64(userId))))
		case "ret":
			return
		default:
			msg.Text = "I don't understand this command :("
		}
		if msg.Text != "" {
			msg.BaseChat.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}
