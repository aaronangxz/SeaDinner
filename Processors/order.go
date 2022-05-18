package Processors

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func OrderDinner(client resty.Client, menuID int, u UserChoiceWithKeyAndStatus) OrderResponse {
	var resp OrderResponse
	fData := make(map[string]string)
	fData["food_id"] = fmt.Sprint(u.GetUserChoice())

	for i := 1; i <= Config.Runtime.RetryTimes; i++ {
		log.Printf("id: %v | OrderDinner | Attempt %v", u.GetUserID(), i)

		_, err := client.R().
			SetHeader("Authorization", MakeToken(u.GetUserKey())).
			SetFormData(fData).
			SetResult(&resp).
			Post(MakeURL(URL_ORDER, &menuID))

		if err != nil {
			log.Println(err)
			continue
		}

		if resp.Status != nil && resp.GetStatus() == "error" {
			log.Printf("id: %v | %v : %v : %v : %v", u.GetUserID(), resp.GetError(), resp.GetStatus(), resp.GetStatusCode(), resp.GetSelected())
		}

		if resp.GetSelected() != 0 {
			log.Printf("id: %v | Dinner Selected: %d. Successful in %v try.\n", u.GetUserID(), resp.GetSelected(), i)
			break
		}
		log.Printf("id: %v | Dinner Not Selected. Retrying.\n", u.GetUserID())
	}
	return resp
}

func BatchOrderDinner(u *[]UserChoiceWithKeyAndStatus) []OrderRecord {
	var (
		records []OrderRecord
		m       = make(map[int64]int)
	)

	for i := 0; i < len(*u); i++ {
		r := (*u)[i]
		log.Printf("id: %v | BatchOrderDinner | Ordering\n", r.GetUserID())
		start := time.Now().UnixMilli()
		resp := OrderDinner(Client, GetDayId(), r)

		record := OrderRecord{
			UserID:    Int64(r.GetUserID()),
			FoodID:    String(r.GetUserChoice()),
			OrderTime: Int64(time.Now().Unix()),
		}

		if resp.GetSelected() == 0 {
			m[r.GetUserID()] = ORDER_STATUS_FAIL
			if resp.Error == nil {
				record.ErrorMsg = String("Unknown Error")
			} else {
				record.ErrorMsg = String(resp.GetError())
			}
			record.Status = Int64(int64(m[r.GetUserID()]))
		} else {
			elapsed := time.Now().UnixMilli() - start
			m[r.GetUserID()] = ORDER_STATUS_OK
			record.Status = Int64(int64(m[r.GetUserID()]))
			SendInstantNotification(r, elapsed)
			//Truncate successful orders so it won't be repeated next round
			*u = PopSuccessfulOrder(*u, i)
			i--
		}
		records = append(records, record)
	}
	log.Println("BatchOrderDinner | Records:", len(records))
	return records
}

func UpdateOrderLog(records []OrderRecord) {
	for _, r := range records {
		log.Printf("id : %v | UpdateOrderLog | status: %v", r.GetUserID(), r.GetStatus())
		if err := DB.Exec("INSERT INTO order_log_tab (user_id, food_id, order_time, status, error_msg) VALUES (?,?,?,?,?)",
			r.GetUserID(), r.GetFoodID(), r.GetOrderTime(), r.GetStatus(), r.GetErrorMsg()).Error; err != nil {
			log.Printf("id : %v | Failed to update record.", r.GetUserID())
		}
	}
}

func SendInstantNotification(u UserChoiceWithKeyAndStatus, took int64) {
	var (
		msg string
	)
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	menu := MakeMenuMap()
	msg = fmt.Sprintf("Successfully ordered %v in %vms! 🥳", menu[u.GetUserChoice()], took)

	if _, err := bot.Send(tgbotapi.NewMessage(u.GetUserID(), msg)); err != nil {
		log.Println(err)
	}
	log.Printf("SendInstantNotification | user_id:%v | msg: %v", u.GetUserID(), msg)
}

func MakeMenuMap() map[string]string {
	key := os.Getenv("TOKEN")
	menuMap := make(map[string]string)

	menu := GetMenu(Client, GetDayId(), key)

	for _, m := range menu.DinnerArr {
		menuMap[fmt.Sprint(m.Id)] = m.Name
	}
	return menuMap
}
