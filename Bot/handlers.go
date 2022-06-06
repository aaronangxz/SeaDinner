package Bot

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/Processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/proto"
)

//GetKey Retrieves user's API key with user_id.
//Reads from cache first, then user_key_tab.
func GetKey(id int64) string {
	var (
		existingRecord *sea_dinner.UserKey
		cacheKey       = fmt.Sprint(Common.USER_KEY_PREFIX, id)
		expiry         = 604800 * time.Second
	)

	if id <= 0 {
		log.Println("GetKey | Id must be > 1.")
		return ""
	}

	//Read from cache
	val, redisErr := Processors.RedisClient.Get(cacheKey).Result()
	if redisErr != nil {
		if redisErr == redis.Nil {
			log.Printf("GetKey | No result of %v in Redis, reading from DB", cacheKey)
		} else {
			log.Printf("GetKey | Error while reading from redis: %v", redisErr.Error())
		}
	} else {
		redisResp := &sea_dinner.UserKey{}
		err := json.Unmarshal([]byte(val), &redisResp)
		if err != nil {
			log.Printf("GetKey | Fail to unmarshal Redis value of key %v : %v, reading from DB", cacheKey, err)
		} else {
			log.Printf("GetKey | Successful | Cached %v", cacheKey)
			return redisResp.GetUserKey()
		}
	}

	//Read from DB
	if err := Processors.DB.Table(Common.DB_USER_KEY_TAB).Where("user_id = ?", id).First(&existingRecord).Error; err != nil {
		return ""
	}

	//set back into cache
	data, err := json.Marshal(existingRecord)
	if err != nil {
		log.Printf("GetKey | Failed to marshal JSON results: %v\n", err.Error())
	}

	if err := Processors.RedisClient.Set(cacheKey, data, expiry).Err(); err != nil {
		log.Printf("GetKey | Error while writing to redis: %v", err.Error())
	} else {
		log.Printf("GetKey | Successful | Written %v to redis", cacheKey)
	}

	return existingRecord.GetUserKey()
}

//CheckKey Checks if user's API key exists.
//Reads from cache first, then user_key_tab.
func CheckKey(id int64) (string, bool) {
	var (
		existingRecord *sea_dinner.UserKey
		cacheKey       = fmt.Sprint(Common.USER_KEY_PREFIX, id)
		expiry         = 604800 * time.Second
	)

	if id <= 0 {
		log.Println("CheckKey | Id must be > 1.")
		return "", false
	}

	//Read from cache
	val, redisErr := Processors.RedisClient.Get(cacheKey).Result()
	if redisErr != nil {
		if redisErr == redis.Nil {
			log.Printf("CheckKey | No result of %v in Redis, reading from DB", cacheKey)
		} else {
			log.Printf("CheckKey | Error while reading from redis: %v", redisErr.Error())
		}
	} else {
		redisResp := &sea_dinner.UserKey{}
		err := json.Unmarshal([]byte(val), &redisResp)
		if err != nil {
			log.Printf("CheckKey | Fail to unmarshal Redis value of key %v : %v, reading from DB", cacheKey, err)
		} else {
			log.Printf("CheckKey | Successful | Cached %v", cacheKey)
			decrypt := Processors.DecryptKey(redisResp.GetUserKey(), os.Getenv("AES_KEY"))
			return fmt.Sprintf("I have your key %v***** that you told me on %v! But I won't leak it 😀", decrypt[:5], Processors.ConvertTimeStamp(redisResp.GetMtime())), true
		}
	}

	//Read from DB
	if err := Processors.DB.Table(Common.DB_USER_KEY_TAB).Where("user_id = ?", id).First(&existingRecord).Error; err != nil {
		return "I don't have your key, let me know in /newkey 😊", false
	} else {
		//set back into cache
		data, err := json.Marshal(existingRecord)
		if err != nil {
			log.Printf("CheckKey | Failed to marshal JSON results: %v\n", err.Error())
		}

		if err := Processors.RedisClient.Set(cacheKey, data, expiry).Err(); err != nil {
			log.Printf("CheckKey | Error while writing to redis: %v", err.Error())
		} else {
			log.Printf("CheckKey | Successful | Written %v to redis", cacheKey)
		}
		decrypt := Processors.DecryptKey(existingRecord.GetUserKey(), os.Getenv("AES_KEY"))
		return fmt.Sprintf("I have your key %v***** that you told me on %v! But I won't leak it 😀", decrypt[:5], Processors.ConvertTimeStamp(existingRecord.GetMtime())), true
	}
}

//UpdateKey Creates record to store user's key if not exists, or update the existing record.
//With basic parameter verifications
func UpdateKey(id int64, s string) (string, bool) {
	hashedKey := Processors.EncryptKey(s, os.Getenv("AES_KEY"))

	var (
		cacheKey       = fmt.Sprint(Common.USER_KEY_PREFIX, id)
		existingRecord sea_dinner.UserKey
		r              = &sea_dinner.UserKey{
			UserId:  proto.Int64(id),
			UserKey: proto.String(hashedKey),
			Ctime:   proto.Int64(time.Now().Unix()),
			Mtime:   proto.Int64(time.Now().Unix()),
		}
	)

	if id <= 0 {
		log.Println("UpdateKey | Id must be > 1.")
		return "", false
	}

	if s == "" {
		log.Println("UpdateKey | Key cannot be empty.")
		return "Key cannot be empty 😟", false
	}

	if len(s) != 40 {
		log.Printf("UpdateKey | Key length invalid | length: %v", len(s))
		return "Are you sure this is a valid key? 😟", false
	}

	if err := Processors.DB.Raw("SELECT * FROM user_key_tab WHERE user_id = ?", id).Scan(&existingRecord).Error; err != nil {
		log.Printf("UpdateKey | %v", err.Error())
		return err.Error(), false
	} else {
		if existingRecord.UserId == nil {
			if err := Processors.DB.Table(Common.DB_USER_KEY_TAB).Create(&r).Error; err != nil {
				log.Println("UpdateKey | Failed to insert DB")
				return err.Error(), false
			}
			return "Okay got it. I remember your key now! 😙\n Disclaimer: I will never disclose your key. Your key is safely encrypted.", true
		}
		//Update key if user_id exists
		if err := Processors.DB.Exec("UPDATE user_key_tab SET user_key = ?, mtime = ? WHERE user_id = ?", hashedKey, time.Now().Unix(), id).Error; err != nil {
			log.Printf("UpdateKey | %v", err.Error())
			return err.Error(), false
		}

		//Invalidate cache after successful update
		if _, err := Processors.RedisClient.Del(cacheKey).Result(); err != nil {
			log.Printf("UpdateKey | Failed to invalidate cache: %v. %v", cacheKey, err)
		}
		log.Printf("UpdateKey | Successfully invalidated cache: %v", cacheKey)

		return "Okay got it. I will take note of your new key 😙", true
	}
}

//CheckChope Retrieves the current food choice made by user.
func CheckChope(id int64) (string, bool) {
	var (
		existingRecord sea_dinner.UserChoice
		dayText        = "today"
	)

	if id <= 0 {
		log.Println("Id must be > 1.")
		return "", false
	}

	if err := Processors.DB.Raw("SELECT * FROM user_choice_tab WHERE user_id = ?", id).Scan(&existingRecord).Error; err != nil {
		return "I have yet to receive your order 🥲 You can choose from /menu", false
	} else {
		if existingRecord.UserChoice == nil {
			return "I have yet to receive your order 🥲 You can choose from /menu", false
		} else if existingRecord.GetUserChoice() == "-1" {
			//Dynamic text based on time - shows tomorrow if current time is past lunch
			tz, _ := time.LoadLocation(Processors.TimeZone)
			if time.Now().In(tz).Unix() > Processors.GetLunchTime().Unix() {
				if Processors.IsNotEOW(time.Now().In(tz)) {
					dayText = "tomorrow"
				} else {
					//On fridays ~ sundays
					return "We are done for this week! You can tell me your order again next week 😀", false
				}
			}
			return fmt.Sprintf("Not placing dinner order for you %v 🙅 Changed your mind? You can choose from /menu", dayText), false
		}
		menu := MakeMenuNameMap()

		_, ok := menu[existingRecord.GetUserChoice()]

		if !ok {
			return fmt.Sprintf("Your choice %v is not available this week, so I will not order anything 🥲 Choose a new dish from /menu", existingRecord.GetUserChoice()), true
		}

		if existingRecord.GetUserChoice() == "RAND" {
			return "I'm tasked to snatch a random dish for you 😀 Changed your mind? You can choose from /menu", true
		}
		return fmt.Sprintf("I'm tasked to snatch %v for you 😀 Changed your mind? You can choose from /menu", menu[existingRecord.GetUserChoice()]), true
	}
}

//GetChope Updates the current food choice made by user.
//With basic parameter verifications
//Supports Button Callbacks
func GetChope(id int64, s string) (string, bool) {
	var (
		existingRecord sea_dinner.UserChoice
		r              = &sea_dinner.UserChoice{
			UserId:     proto.Int64(id),
			UserChoice: proto.String(s),
			Ctime:      proto.Int64(time.Now().Unix()),
			Mtime:      proto.Int64(time.Now().Unix()),
		}
		key = fmt.Sprint(Common.USER_CHOICE_PREFIX, r.GetUserId())
	)

	if id <= 0 {
		log.Println("Id must be > 1.")
		return "", false
	}

	//When it is Friday after 12.30pm, we don't accept any orders (except -1) because we don't know next week's menu yet
	if !Processors.IsNotEOW(time.Now()) && time.Now().Unix() > Processors.GetLunchTime().Unix() && s != "-1" {
		return "We are done for this week! You can tell me your order again next week 😀", false
	}

	if Processors.IsNotNumber(s) {
		//RAND is passed from CallBack
		if s != "RAND" && s != "SAME" {
			log.Printf("Selection contains illegal character | selection: %v", s)
			return "Are you sure that is a valid FoodID? Tell me another one. 😟", false
		}
	}

	menu := MakeMenuNameMap()

	if s == "SAME" {
		//Set back to DB using cache data (user_choice:<id>)
		//To handle Morning Reminder callback
		val, redisErr := Processors.RedisClient.Get(key).Result()
		if redisErr != nil {
			if redisErr == redis.Nil {
				log.Printf("GetChope | No result of %v in Redis, reading from API", key)
			} else {
				log.Printf("GetChope | Error while reading from redis: %v", redisErr.Error())
			}
			return "The selection has expired, you can choose from /menu again 😀", true
		}

		if val == "" {
			log.Printf("GetChope | empty in redis: %v", key)
			return "The selection has expired, you can choose from /menu again 😀", true
		}

		if err := Processors.DB.Exec("UPDATE user_choice_tab SET user_choice = ?, mtime = ? WHERE user_id = ?", val, time.Now().Unix(), id).Error; err != nil {
			log.Println("Failed to update DB")
			return err.Error(), false
		}

		if val == "-1" {
			return "Okay got it. I will not order anything for you instead.😀", true
		}

		if val == "RAND" {
			return "Okay got it. I will give you a surprise 😙", true
		}
		return fmt.Sprintf("Okay got it! I will order %v 😙", menu[val]), true
	}

	_, ok := menu[s]
	if !ok {
		log.Printf("Selection is invalid | selection: %v", s)
		return "This dish is not available today. Tell me another one.😟", false
	}

	if err := Processors.DB.Raw("SELECT * FROM user_choice_tab WHERE user_id = ?", id).Scan(&existingRecord).Error; err != nil {
		log.Printf("GetChope | %v", err.Error())
		return err.Error(), false
	} else {
		if existingRecord.UserId == nil {
			if err := Processors.DB.Table(Common.DB_USER_CHOICE_TAB).Create(&r).Error; err != nil {
				log.Println("Failed to insert DB")
				return err.Error(), false
			}

			//To stop ordering
			if s == "-1" {
				return fmt.Sprintf("Okay got it. I will order %v for you and stop sending reminders in the morning for the rest of the week.😀", menu[s]), true
			}

			if s == "RAND" {
				return "Okay got it. I will give you a surprise😙", true
			}

			//Orders placed before lunch time
			if time.Now().Unix() < Processors.GetLunchTime().Unix() {
				return fmt.Sprintf("Okay got it. I will order %v for you today 😙", menu[s]), true
			}

			return fmt.Sprintf("Okay got it. I will order %v for you tomorrow 😙", menu[s]), true
		}
		//Update key if user_id exists
		if err := Processors.DB.Exec("UPDATE user_choice_tab SET user_choice = ?, mtime = ? WHERE user_id = ?", s, time.Now().Unix(), id).Error; err != nil {
			log.Println("Failed to update DB")
			return err.Error(), false
		}

		//To stop ordering
		if s == "-1" {
			return fmt.Sprintf("Okay got it. I will order %v for you and stop sending reminders in the morning for the rest of the week.😀", menu[s]), true
		}

		if s == "RAND" {
			return "Okay got it. I will give you a surprise 😙", true
		}

		//Orders placed before lunch time
		if time.Now().Unix() < Processors.GetLunchTime().Unix() {
			//Set into cache for Morning reminder callback. TTL is always until 12.30
			if err := Processors.RedisClient.Set(key, s, time.Duration(Processors.GetLunchTime().UnixMilli()-time.Now().UnixMilli())).Err(); err != nil {
				log.Printf("GetChope | Error while writing to redis: %v", err.Error())
			} else {
				log.Printf("GetChope | Successful | Written %v to redis", key)
			}

			return fmt.Sprintf("Okay got it. I will order %v for you today 😙", menu[s]), true
		}

		return fmt.Sprintf("Okay got it. I will order %v for you tomorrow 😙", menu[s]), true
	}
}

//ListWeeklyResultByUserId Returns the order records of a user in the current week
func ListWeeklyResultByUserId(id int64) string {
	var (
		res []*sea_dinner.OrderRecord
	)

	start, end := Processors.WeekStartEndDate(time.Now().Unix())

	if id <= 0 {
		log.Println("Id must be > 1.")
		return ""
	}

	if err := Processors.DB.Raw("SELECT * FROM order_log_tab WHERE user_id = ? AND order_time BETWEEN ? AND ?", id, start, end).Scan(&res).Error; err != nil {
		log.Printf("id : %v | Failed to retrieve record.", id)
		return "You have not ordered anything this week. 😕"
	}

	if res == nil {
		return "You have not ordered anything this week. 😕"
	}
	return GenerateWeeklyResultTable(res)
}

//GenerateWeeklyResultTable Outputs pre-formatted weekly order status.
func GenerateWeeklyResultTable(record []*sea_dinner.OrderRecord) string {
	start, end := Processors.WeekStartEndDate(time.Now().Unix())
	m := MakeMenuCodeMap()
	status := map[int64]string{int64(sea_dinner.OrderStatus_ORDER_STATUS_OK): "✅", int64(sea_dinner.OrderStatus_ORDER_STATUS_FAIL): "❌"}
	header := fmt.Sprintf("Your orders from %v to %v\n", Processors.ConvertTimeStampMonthDay(start), Processors.ConvertTimeStampMonthDay(end))
	table := "<pre>\n     Day    Code  Status\n"
	table += "-------------------------\n"
	for _, r := range record {
		table += fmt.Sprintf(" %v   %v     %v\n", Processors.ConvertTimeStampDayOfWeek(r.GetOrderTime()), m[r.GetFoodId()], status[r.GetStatus()])
	}
	table += "</pre>"
	return header + table
}

//BatchGetLatestResult Retrieves the most recent failed orders
func BatchGetLatestResult() []*sea_dinner.OrderRecord {
	var (
		res []*sea_dinner.OrderRecord
	)

	if err := Processors.DB.Raw("SELECT ol.* FROM order_log_tab ol INNER JOIN "+
		"(SELECT MAX(order_time) AS max_order_time FROM order_log_tab WHERE status <> ? AND order_time BETWEEN ? AND ? GROUP BY user_id) nestedQ "+
		"ON ol.order_time = nestedQ.max_order_time GROUP BY user_id",
		sea_dinner.OrderStatus_ORDER_STATUS_OK, Processors.GetLunchTime().Unix()-300, Processors.GetLunchTime().Unix()+300).
		Scan(&res).Error; err != nil {
		log.Println("Failed to retrieve record.")
		return nil
	}
	log.Println("BatchGetLatestResult:", len(res))
	return res
}

//SendNotifications Sends out notifications based on order status from BatchGetLatestResult
//Used to send failed orders only
func SendNotifications() {
	var (
		msg string
	)
	bot, err := tgbotapi.NewBotAPI(Common.GetTGToken())
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	res := BatchGetLatestResult()
	menu := MakeMenuNameMap()
	log.Println("SendNotifications | size:", len(res))

	for _, r := range res {
		if r.GetStatus() == int64(sea_dinner.OrderStatus_ORDER_STATUS_OK) {
			msg = fmt.Sprintf("Successfully ordered %v! 🥳", menu[r.GetFoodId()])
		} else {
			msg = fmt.Sprintf("Failed to order %v today. %v 😔", menu[r.GetFoodId()], r.GetErrorMsg())
		}

		if _, err := bot.Send(tgbotapi.NewMessage(r.GetUserId(), msg)); err != nil {
			log.Println(err)
		}
	}
}

//BatchGetUsersChoice Retrieves order_choice of all users
func BatchGetUsersChoice() []*sea_dinner.UserChoice {
	var (
		res    []*sea_dinner.UserChoice
		expiry = 7200 * time.Second
	)
	if err := Processors.DB.Raw("SELECT uc.* FROM user_choice_tab uc, user_key_tab uk WHERE uc.user_id = uk.user_id AND uk.is_mute <> ?", sea_dinner.MuteStatus_MUTE_STATUS_YES).Scan(&res).Error; err != nil {
		log.Println("BatchGetUsersChoice | Failed to retrieve record:", err.Error())
		return nil
	}

	//Save into cache
	//For Morning Reminder callback
	for _, r := range res {
		//Not neccesary to cache -1 orders because we never send reminder for those
		if r.GetUserChoice() != "-1" {
			key := fmt.Sprint(Common.USER_CHOICE_PREFIX, r.GetUserId())
			if err := Processors.RedisClient.Set(key, r.GetUserChoice(), expiry).Err(); err != nil {
				log.Printf("BatchGetUsersChoice | Error while writing to redis: %v", err.Error())
			} else {
				log.Printf("BatchGetUsersChoice | Successful | Written %v to redis", key)
			}
		}
	}
	log.Println("BatchGetUsersChoice | size:", len(res))
	log.Println(res)
	return res
}

//SendReminder Sends out daily reminder at 10.30 SGT on weekdays / working days
func SendReminder() {
	txn := Processors.App.StartTransaction("send_reminder")
	defer txn.End()

	bot, err := tgbotapi.NewBotAPI(Common.GetTGToken())
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	res := BatchGetUsersChoice()
	log.Println("SendReminder | size:", len(res))

	menu := MakeMenuNameMap()
	code := MakeMenuCodeMap()

	for _, r := range res {
		msg := tgbotapi.NewMessage(r.GetUserId(), "")

		var (
			mk     tgbotapi.InlineKeyboardMarkup
			out    [][]tgbotapi.InlineKeyboardButton
			rows   []tgbotapi.InlineKeyboardButton
			msgTxt string
		)

		if Processors.IsSOW(time.Now()) {
			//Everyone exceot "MUTE" will receive weekly reminders
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
				log.Printf("SendReminder | skip -1 records | %v", r.GetUserId())
				continue
			}

			_, ok := menu[r.GetUserChoice()]
			if !ok {
				msgTxt = fmt.Sprintf("Good Morning. Your previous order %v is not available today! I will not proceed to order. Choose another dish from /menu 😃 /mute to shut me up 🫢 ", r.GetUserChoice())
			} else {
				if r.GetUserChoice() != "-1" {
					//If choice was updated after yesterdays' lunch time
					if r.GetMtime() > Processors.GetPreviousDayLunchTime().Unix() {
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
			log.Println(err)
		}
	}
}

//MakeMenuNameMap Returns food_id:food_name mapping of current menu
func MakeMenuNameMap() map[string]string {
	var (
		key = os.Getenv("TOKEN")
	)
	menuMap := make(map[string]string)
	menu := Processors.GetMenu(Processors.Client, key)
	for _, m := range menu.GetFood() {
		menuMap[fmt.Sprint(m.GetId())] = m.GetName()
	}
	// Store -1 hash to menuMap
	menuMap["-1"] = "*NOTHING*" // to be renamed
	menuMap["RAND"] = "Random"
	return menuMap
}

//MakeMenuCodeMap Returns food_id:food_code mapping of current menu
func MakeMenuCodeMap() map[string]string {
	var (
		key = os.Getenv("TOKEN")
	)
	menuMap := make(map[string]string)
	menu := Processors.GetMenu(Processors.Client, key)
	for _, m := range menu.GetFood() {
		menuMap[fmt.Sprint(m.GetId())] = m.GetCode()
	}
	menuMap["RAND"] = "Random"
	return menuMap
}

//CallbackQueryHandler Handles the call back result of menu buttons
func CallbackQueryHandler(id int64, callBack *tgbotapi.CallbackQuery) (string, bool) {
	log.Printf("id: %v | CallbackQueryHandler | callback: %v", id, callBack.Data)

	if callBack.Data == "MUTE" || callBack.Data == "UNMUTE" {
		return UpdateMute(id, callBack.Data)
	}

	return GetChope(id, callBack.Data)
}

//MakeHelpResponse Prints out Introduction
func MakeHelpResponse() string {
	return "*Welcome to SeaHungerGamesBot!*\n\n" +
		"The goal of my existence is to help you snatch that dinner in milliseconds. And also we all know that you are too lazy to open up SeaTalk.\n\n" +
		"*Get started*\n" +
		"1. /key to tell me your Sea API key. This is important because without the key, I'm basically useless. When you refresh your key, remember to let me know in /newkey\n" +
		"2. /menu to browse through the dishes, and tap the button below to snatch. There are also options to choose a random dish or skip ordering. Do take note that if you choose to skip, I will remember that and stop ordering forever until you tell me to do so again.\n" +
		"3. /choice to check the current dish I'm tasked to order.\n" +
		"4. /status to see what you have ordered this week, and the order status.\n" +
		"5. /mute to stop receiving morning reminders. Not recommended tho!\n\n" +
		"*Features*\n" +
		"1. I will send you a daily reminder at 10.30am (If you never mute or skip order on that day). Order can be altered easily from the quick options:\n" +
		"🎲 to order a random dish\n" +
		"🙅 to stop ordering\n" +
		"2. At 12.29pm, I will no longer entertain your requests, because I have better things to do! Don't even think about last minute changes.\n" +
		"3. At 12.30pm sharp, I will begin to order your precious food.\n" +
		"4. It is almost guranteed that I can order it in less than 500ms. Will drop you a message too!\n\n" +
		"*Disclaimer*\n" +
		"By using my services, you agree to let me store your API key. However, not to worry! Your key is encrypted with AES-256, it's very unlikely that it will be stolen.\n\n" +
		"*Contribute*\n" +
		"If you see or encounter any bugs, or if there's any feature / improvement that you have in mind, feel free to open an Issue / Pull Request at https://github.com/aaronangxz/SeaDinner\n\n" +
		"Thank you and happy eating!😋"
}

func CheckMute(id int64) (string, []tgbotapi.InlineKeyboardMarkup) {
	var (
		res *sea_dinner.UserKey
		out []tgbotapi.InlineKeyboardMarkup
	)
	if err := Processors.DB.Raw("SELECT * FROM user_key_tab WHERE user_id = ?", id).Scan(&res).Error; err != nil {
		log.Println("CheckMute | Failed to retrieve record:", err.Error())
		return "", nil
	}

	if res == nil {
		log.Printf("CheckMute | Record not found | user_id:%v", id)
		return "Record not found.", nil
	}

	if res.GetIsMute() == int64(sea_dinner.MuteStatus_MUTE_STATUS_NO) {
		var rows []tgbotapi.InlineKeyboardButton
		muteBotton := tgbotapi.NewInlineKeyboardButtonData("Turn OFF 🔕", "MUTE")
		rows = append(rows, muteBotton)
		out = append(out, tgbotapi.NewInlineKeyboardMarkup(rows))
		return "Daily reminder notifications are *ON*.\nDo you want to turn it OFF?", out
	}
	var rows []tgbotapi.InlineKeyboardButton
	unmuteBotton := tgbotapi.NewInlineKeyboardButtonData("Turn ON 🔔", "UNMUTE")
	rows = append(rows, unmuteBotton)
	out = append(out, tgbotapi.NewInlineKeyboardMarkup(rows))
	return "Daily reminder notifications are *OFF*.\nDo you want to turn it ON?", out
}

func UpdateMute(id int64, callback string) (string, bool) {
	var (
		toUdate    = int64(sea_dinner.MuteStatus_MUTE_STATUS_YES)
		returnMsg  = "Daily reminder notifications are *OFF*.\nDo you want to turn it ON?"
		returnBool = true
	)

	if callback == "UNMUTE" {
		toUdate = int64(sea_dinner.MuteStatus_MUTE_STATUS_NO)
		returnMsg = "Daily reminder notifications are *ON*.\nDo you want to turn it OFF?"
		returnBool = false
	}

	if err := Processors.DB.Exec("UPDATE user_key_tab SET is_mute = ? WHERE user_id = ?", toUdate, id).Error; err != nil {
		log.Println("Failed to update DB")
		return err.Error(), false
	}

	return returnMsg, returnBool
}
