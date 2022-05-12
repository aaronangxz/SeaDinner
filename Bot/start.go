package Bot

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	u.Timeout = 1
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if time.Now().Unix() > Processors.GetLunchTime().Unix()-30 && time.Now().Unix() < Processors.GetLunchTime().Unix()+30 {
			block = true
			if _, err := bot.Send(tgbotapi.NewMessage(Id, "Omw to order, wait for my good news!")); err != nil {
				log.Println(err)
			}
			return
		} else {
			//check result
			block = false
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
			msg.Text = "Check the commands."
		case "key":
			msg.Text, _ = checkKey(update.Message.Chat.ID)
		case "newkey":
			msg.Text = "What's your key? \nGo to https://dinner.sea.com/accounts/token, copy the Key under Generate Auth Token and paste it here:"
			startListenKey = true
		case "status":
			msg.Text = getLatestResultByUserId(update.Message.Chat.ID)
		case "chope":
			msg.Text = "What do you want to order? Tell me the Food ID 😋"
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
		existingRecord UserKey
	)

	if err := Processors.DB.Table("user_key").Where("user_id = ?", id).First(&existingRecord).Error; err != nil {
		return "I don't have your key, let me know in /newkey 😊", false
	} else {
		Key = existingRecord.Key
		t := time.Unix(int64(existingRecord.Mtime), 0).UTC()
		return fmt.Sprintf("I have your key that you told me on %v! But I won't leak it 😀", t.Format("2006-01-02")), true
	}
}

func updateKey(id int64, s string) string {
	var (
		existingRecord UserKey
		r              = UserKey{
			UserID: id,
			Key:    s,
			Ctime:  time.Now().Unix(),
			Mtime:  time.Now().Unix(),
		}
	)

	Key = s

	if err := Processors.DB.Raw("SELECT * FROM user_key WHERE user_id = ? AND key = ?", id, s).Scan(&existingRecord).Error; err != nil {
		//Insert new row
		if err := Processors.DB.Table("user_key").Create(&r).Error; err != nil {
			log.Println("Failed to insert DB")
			return err.Error()
		}
		return "Okay got it. I remember your key now! 😙"
	} else {
		//Update key if user_id exists
		if err := Processors.DB.Exec("UPDATE user_key SET key = ?, mtime = ? WHERE user_id = ?", s, time.Now().Unix(), id).Error; err != nil {
			log.Println("Failed to update DB")
			return err.Error()
		}
		return "Okay got it. I will take note of your new key 😙"
	}
}

func checkChope(id int64) (string, bool) {
	var (
		existingRecord UserChoice
	)

	if err := Processors.DB.Raw("SELECT * FROM user_choice WHERE user_id = ?", id).Scan(&existingRecord).Error; err != nil {
		return "I have yet to receive your order 🥲 Tell me at /chope", false
	} else {
		return fmt.Sprintf("I'm tasked to snatch %v for you 😀 Changed your mind? Tell me at /chope", existingRecord.Choice), true
	}
}

func getChope(id int64, s string) string {
	n, _ := strconv.ParseInt(s, 10, 64)
	var (
		existingRecord UserChoice
		r              = UserChoice{
			UserID: id,
			Choice: n,
			Ctime:  time.Now().Unix(),
			Mtime:  time.Now().Unix(),
		}
	)

	if err := Processors.DB.Raw("SELECT * FROM user_choice WHERE user_id = ?", id).Scan(&existingRecord).Error; err != nil {
		//Insert new row
		if err := Processors.DB.Table("user_choice").Create(&r).Error; err != nil {
			log.Println("Failed to insert DB")
			return err.Error()
		}
		return fmt.Sprintf("Okay got it. I will order %v for you 😙", s)
	} else {
		//Update key if user_id exists
		if err := Processors.DB.Exec("UPDATE user_choice SET choice = ?, mtime = ? WHERE user_id = ?", s, time.Now().Unix(), id).Error; err != nil {
			log.Println("Failed to update DB")
			return err.Error()
		}
		return fmt.Sprintf("Okay got it. I will order %v for you 😙", s)
	}
}

func getLatestResultByUserId(id int64) string {
	var (
		res Processors.OrderRecord
	)
	if err := Processors.DB.Raw("SELECT * FROM order_log WHERE user_id = ? AND order_time BETWEEN ? AND ? ORDER BY order_time DESC LIMIT 1", id, Processors.GetLunchTime().Unix()-3600, Processors.GetLunchTime().Unix()+3600).Scan(&res).Error; err != nil {
		log.Printf("id : %v | Failed to retrieve record.", id)
		return "Unable to find record for today."
	}

	if res.GetStatus() == Processors.ORDER_STATUS_OK {
		return fmt.Sprintf("Successfully ordered %v!", res.GetFoodID())
	}
	return fmt.Sprintf("Failed to order %v today.", res.GetFoodID())
}

func batchGetLatestResult() []Processors.OrderRecord {
	var (
		res []Processors.OrderRecord
	)
	if err := Processors.DB.Raw("SELECT * FROM order_log WHERE order_time BETWEEN ? AND ? ORDER BY order_time DESC LIMIT 1", Processors.GetLunchTime().Unix()-3600, Processors.GetLunchTime().Unix()+3600).Scan(&res).Error; err != nil {
		log.Println("Failed to retrieve record.")
		return nil
	}
	return res
}

func SendNotifications() {
	var (
		msg string
	)
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	res := batchGetLatestResult()

	for _, r := range res {
		if r.GetStatus() == Processors.ORDER_STATUS_OK {
			msg = fmt.Sprintf("Successfully ordered %v!", r.GetFoodID())
		} else {
			msg = fmt.Sprintf("Failed to order %v today.", r.GetFoodID())
		}

		if _, err := bot.Send(tgbotapi.NewMessage(r.GetUserID(), msg)); err != nil {
			log.Println(err)
		}
	}
}
