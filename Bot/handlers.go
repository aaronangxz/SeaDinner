package Bot

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aaronangxz/SeaDinner/Processors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetKey(id int64) string {
	var (
		existingRecord UserKey
	)

	if id <= 0 {
		log.Println("Id must be > 1.")
		return ""
	}

	if err := Processors.DB.Table(Processors.DB_USER_KEY_TAB).Where("user_id = ?", id).First(&existingRecord).Error; err != nil {
		return ""
	}
	return existingRecord.GetUserKey()
}

func CheckKey(id int64) (string, bool) {
	var (
		existingRecord UserKey
	)

	if id <= 0 {
		log.Println("Id must be > 1.")
		return "", false
	}

	if err := Processors.DB.Table(Processors.DB_USER_KEY_TAB).Where("user_id = ?", id).First(&existingRecord).Error; err != nil {
		return "I don't have your key, let me know in /newkey 😊", false
	} else {
		return fmt.Sprintf("I have your key that you told me on %v! But I won't leak it 😀", Processors.ConvertTimeStamp(existingRecord.GetMtime())), true
	}
}

func UpdateKey(id int64, s string) (string, bool) {
	hashedKey := Processors.EncryptKey(s, os.Getenv("AES_KEY"))

	var (
		existingRecord UserKey
		r              = UserKey{
			UserID:  Processors.Int64(id),
			UserKey: Processors.String(hashedKey),
			Ctime:   Processors.Int64(time.Now().Unix()),
			Mtime:   Processors.Int64(time.Now().Unix()),
		}
	)

	if id <= 0 {
		log.Println("Id must be > 1.")
		return "", false
	}

	if s == "" {
		log.Println("Key cannot be empty.")
		return "Key cannot be empty 😟", false
	}

	if len(s) != 40 {
		log.Printf("Key length invalid | length: %v", len(s))
		return "Are you sure this is a valid key? 😟", false
	}

	if err := Processors.DB.Raw("SELECT * FROM user_key_tab WHERE user_id = ?", id).Scan(&existingRecord).Error; err != nil {
		//Insert new row
		if existingRecord.UserID == nil {
			if err := Processors.DB.Table(Processors.DB_USER_KEY_TAB).Create(&r).Error; err != nil {
				log.Println("Failed to insert DB")
				return err.Error(), false
			}
			return "Okay got it. I remember your key now! 😙", true
		}
		return err.Error(), false
	} else {
		//Update key if user_id exists
		if err := Processors.DB.Exec("UPDATE user_key_tab SET user_key = ?, mtime = ? WHERE user_id = ?", hashedKey, time.Now().Unix(), id).Error; err != nil {
			log.Println("Failed to update DB")
			return err.Error(), false
		}
		return "Okay got it. I will take note of your new key 😙", true
	}
}

func CheckChope(id int64) (string, bool) {
	var (
		existingRecord UserChoice
	)

	if id <= 0 {
		log.Println("Id must be > 1.")
		return "", false
	}

	if err := Processors.DB.Raw("SELECT * FROM user_choice_tab WHERE user_id = ?", id).Scan(&existingRecord).Error; err != nil {
		return "I have yet to receive your order 🥲 You can choose from /menu", false
	} else {
		menu := MakeMenuMap()
		return fmt.Sprintf("I'm tasked to snatch %v for you 😀 Changed your mind? You can choose from /menu", menu[existingRecord.GetUserChoice()]), true
	}
}

func GetChope(id int64, s string) (string, bool) {
	var (
		existingRecord UserChoice
		r              = UserChoice{
			UserID:     Processors.Int64(id),
			UserChoice: Processors.String(s),
			Ctime:      Processors.Int64(time.Now().Unix()),
			Mtime:      Processors.Int64(time.Now().Unix()),
		}
	)

	if id <= 0 {
		log.Println("Id must be > 1.")
		return "", false
	}

	if Processors.IsNotNumber(s) {
		log.Printf("Selection contains illegal character | selection: %v", s)
		return "Are you sure that is a valid FoodID? Tell me another one. 😟", false
	}
	menu := MakeMenuMap()

	_, ok := menu[s]
	if !ok {
		log.Printf("Selection is invalid | selection: %v", s)
		return "This dish is not available today. Tell me another one. 😟", false
	}

	if err := Processors.DB.Raw("SELECT * FROM user_choice_tab WHERE user_id = ?", id).Scan(&existingRecord).Error; err != nil {
		//Insert new row
		if err := Processors.DB.Table(Processors.DB_USER_CHOICE_TAB).Create(&r).Error; err != nil {
			log.Println("Failed to insert DB")
			return err.Error(), false
		}
		return fmt.Sprintf("Okay got it. I will order %v for you 😙", menu[s]), true
	} else {
		//Update key if user_id exists
		if err := Processors.DB.Exec("UPDATE user_choice_tab SET user_choice = ?, mtime = ? WHERE user_id = ?", s, time.Now().Unix(), id).Error; err != nil {
			log.Println("Failed to update DB")
			return err.Error(), false
		}
		return fmt.Sprintf("Okay got it. I will order %v for you 😙", menu[s]), true
	}
}

func GetLatestResultByUserId(id int64) string {
	var (
		res Processors.OrderRecord
	)

	if id <= 0 {
		log.Println("Id must be > 1.")
		return ""
	}

	if err := Processors.DB.Raw("SELECT * FROM order_log_tab WHERE user_id = ? AND order_time BETWEEN ? AND ? ORDER BY order_time DESC LIMIT 1", id, Processors.GetLunchTime().Unix()-3600, Processors.GetLunchTime().Unix()+3600).Scan(&res).Error; err != nil {
		log.Printf("id : %v | Failed to retrieve record.", id)
		return "I have yet to order anything today 😕"
	}
	menu := MakeMenuMap()

	if res.GetStatus() == Processors.ORDER_STATUS_OK {
		return fmt.Sprintf("Successfully ordered %v at %v! 🥳", menu[res.GetFoodID()], Processors.ConvertTimeStampTime(res.GetOrderTime()))
	}
	return fmt.Sprintf("Failed to order %v today. 😔", menu[res.GetFoodID()])
}

func BatchGetLatestResult() []Processors.OrderRecord {
	var (
		res []Processors.OrderRecord
	)

	if err := Processors.DB.Raw("SELECT ol.* FROM order_log_tab ol INNER JOIN "+
		"(SELECT MAX(order_time) AS max_order_time FROM order_log_tab WHERE status <> ? AND order_time BETWEEN ? AND ? GROUP BY user_id) nestedQ "+
		"ON ol.order_time = nestedQ.max_order_time",
		Processors.ORDER_STATUS_OK, Processors.GetLunchTime().Unix()-300, Processors.GetLunchTime().Unix()+300).
		Scan(&res).Error; err != nil {
		log.Println("Failed to retrieve record.")
		return nil
	}
	log.Println("BatchGetLatestResult:", len(res))
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

	res := BatchGetLatestResult()
	menu := MakeMenuMap()
	log.Println("SendNotifications:", len(res))

	for _, r := range res {
		if r.GetStatus() == Processors.ORDER_STATUS_OK {
			msg = fmt.Sprintf("Successfully ordered %v! 🥳", menu[r.GetFoodID()])
		} else {
			msg = fmt.Sprintf("Failed to order %v today. %v😔", menu[r.GetFoodID()], r.GetErrorMsg())
		}

		if _, err := bot.Send(tgbotapi.NewMessage(r.GetUserID(), msg)); err != nil {
			log.Println(err)
		}
	}
}
func BatchGetUsersChoice() []UserChoice {
	var (
		res []UserChoice
	)
	if err := Processors.DB.Raw("SELECT * FROM user_choice_tab").Scan(&res).Error; err != nil {
		log.Println("Failed to retrieve record:", err.Error())
		return nil
	}
	log.Println("BatchGetUsersChoice:", len(res))
	log.Println(res)
	return res
}

func SendReminder() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	res := BatchGetUsersChoice()
	log.Println("SendReminder:", len(res))

	menu := MakeMenuMap()

	for _, r := range res {
		var msg string
		_, ok := menu[r.GetUserChoice()]
		if !ok {
			msg = fmt.Sprintf("Good Morning. Your previous order %v is not available today! I will not proceed to order. Choose another dish from /menu 😃 ", r.GetUserChoice())
		} else {
			msg = fmt.Sprintf("Good Morning. I will order %v again today! If you changed your mind, you can choose from /menu 😋", menu[r.GetUserChoice()])
		}
		if _, err := bot.Send(tgbotapi.NewMessage(r.GetUserID(), msg)); err != nil {
			log.Println(err)
		}
	}
}

func MakeMenuMap() map[string]string {
	key := os.Getenv("TOKEN")
	menuMap := make(map[string]string)

	menu := Processors.GetMenu(Processors.Client, Processors.GetDayId(key), key)

	for _, m := range menu.DinnerArr {
		menuMap[fmt.Sprint(m.Id)] = m.Name
	}
	return menuMap
}

func CallbackQueryHandler(id int64, callBack *tgbotapi.CallbackQuery) (string, bool) {
	log.Printf("id: %v | CallbackQueryHandler | callback: %v", id, callBack)
	return GetChope(id, callBack.Data)
}
