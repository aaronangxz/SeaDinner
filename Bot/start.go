package Bot

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aaronangxz/SeaDinner/Processors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	startListenKey   = false
	startListenChope = false
	Key              = ""
	Id               int64
	block            = false
)

func InitBot() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		go func() {
			fmt.Println("time: ", time.Now().Unix())
			if time.Now().Unix() >= Processors.GetLunchTime().Unix()-300 && time.Now().Unix() <= Processors.GetLunchTime().Unix()+300 {
				block = true
				return
			} else {
				block = false
			}
		}()

		if block {
			if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Omw to order, wait for my good news!")); err != nil {
				log.Println(err)
			}
			bot.StopReceivingUpdates()
			return
		}

		if update.Message.Chat.ID != 0 {
			Id = update.Message.Chat.ID
			fmt.Println("chat id: ", update.Message.Chat.ID)
		}

		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			if startListenKey {
				//Capture key
				if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					updateKey(update.Message.Chat.ID, update.Message.Text))); err != nil {
					log.Println(err)
				}
				startListenKey = false
				continue
			} else if startListenChope {
				//Capture chope
				if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					getChope(update.Message.Chat.ID, update.Message.Text))); err != nil {
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
			s, ok := checkKey(update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				msg.Text = "Hello! " + update.Message.Chat.UserName
			}
		case "menu":
			s, ok := checkKey(update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				msg.Text = Processors.OutputMenu(Key)
			}
		case "help":
			msg.Text = "I understand /sayhi and /status."
		case "key":
			msg.Text, _ = checkKey(update.Message.Chat.ID)
		case "newkey":
			msg.Text = "What's your key?"
			startListenKey = true
		case "status":
			msg.Text = "I'm ok."
		case "chope":
			msg.Text = "What do you want to order? Tell me the Food ID"
			startListenChope = true
		case "choice":
			msg.Text, _ = checkChope(Id)
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

func checkKey(id int64) (string, bool) {
	var (
		existingRecord UserRecord
	)

	if err := Processors.DB.Where("user_id = ?", id).First(&existingRecord).Error; err != nil {
		return "I don't have your key, let me know in /newkey 😊", false
	} else {
		Key = existingRecord.Key
		return "I have your key! But I won't leak it 😀", true
	}
}

func updateKey(id int64, s string) string {
	var (
		existingRecord UserRecord
		r              = UserRecord{
			UserID: id,
			Key:    s,
		}
	)

	Key = s

	if err := Processors.DB.Where("user_id = ?", id).Where("key = ?", s).First(&existingRecord).Error; err != nil {
		//Insert new row
		if err := Processors.DB.Create(&r).Error; err != nil {
			log.Println("Failed to insert DB")
			return err.Error()
		}
		return "Okay got it. I remember your key now! 😙"
	} else {
		//Update key if user_id exists
		if err := Processors.DB.Exec("UPDATE user_records SET key = ? WHERE user_id = ?", s, id).Error; err != nil {
			log.Println("Failed to update DB")
			return err.Error()
		}
		return "Okay got it. I will take note of your new key 😙"
	}
}

func checkChope(id int64) (string, bool) {
	var (
		existingRecord UserRecord
	)

	if err := Processors.DB.Where("user_id = ?", id).First(&existingRecord).Error; err != nil {
		return "I have yet to receive your order 🥲 Tell me at /chope", false
	} else {
		return fmt.Sprintf("I'm tasked to snatch %v for you 😀 Changed your mind? Tell me at /chope", existingRecord.Choice), true
	}
}

func getChope(id int64, s string) string {
	if err := Processors.DB.Exec("UPDATE user_records SET choice = ? WHERE user_id = ?", s, id).Error; err != nil {
		log.Println("Failed to update DB")
		return err.Error()
	}
	return fmt.Sprintf("Okay got it. I will order %v for you 😙", s)
}
