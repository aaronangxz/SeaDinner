package Processors

import (
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

func OrderDinner(client resty.Client, menuID int, u UserChoiceWithKey) OrderResponse {
	var resp OrderResponse
	fData := make(map[string]string)
	fData["food_id"] = fmt.Sprint(u.GetUserChoice())

	for i := 1; i < Config.Runtime.RetryTimes; i++ {
		log.Printf("id: %v | Attempt %v", u.GetUserID(), i)

		start := time.Now().UnixMilli()

		_, err := client.R().
			SetHeader("Authorization", MakeToken(u.GetUserKey())).
			SetFormData(fData).
			SetResult(&resp).
			Post(MakeURL(URL_ORDER, &menuID))

		elapsed := time.Now().UnixMilli() - start

		if err != nil {
			log.Println(err)
			continue
		}

		if resp.Status != nil && resp.GetStatus() == "error" {
			log.Printf("id: %v | %v : %v : %v : %v", u.GetUserID(), resp.GetError(), resp.GetStatus(), resp.GetStatusCode(), resp.GetSelected())
		}

		if resp.GetSelected() != 0 {
			log.Printf("id: %v | Dinner Selected: %d. Successful in %v try. Time: %vms\n", u.GetUserID(), resp.GetSelected(), i, elapsed)
			break
		}
		log.Printf("id: %v | Dinner Not Selected. Retrying.\n", u.GetUserID())
	}
	return resp
}

func BatchOrderDinner(u []UserChoiceWithKey) {
	var (
		records []OrderRecord
		m       = make(map[int64]int)
	)

	for _, r := range u {
		log.Printf("id: %v | Ordering\n", r.GetUserID())
		resp := OrderDinner(Client, GetDayId(r.GetUserKey()), r)

		if resp.GetSelected() == 0 {
			m[r.GetUserID()] = ORDER_STATUS_FAIL
		} else {
			m[r.GetUserID()] = ORDER_STATUS_OK
		}

		record := OrderRecord{
			UserID:    Int64(r.GetUserID()),
			FoodID:    String(r.GetUserChoice()),
			OrderTime: Int64(time.Now().Unix()),
			Status:    Int64(int64(m[r.GetUserID()])),
			ErrorMsg:  String(resp.GetError()),
		}
		records = append(records, record)
	}
	log.Println("records:", len(records))

	UpdateOrderLog(records)
	OutputResults(m)
}

func UpdateOrderLog(records []OrderRecord) {
	for _, r := range records {
		if err := DB.Exec("INSERT INTO order_log_tab (user_id, food_id, order_time, status, error_msg) VALUES (?,?,?,?,?)",
			r.GetUserID(), r.GetFoodID(), r.GetOrderTime(), r.GetStatus(), r.GetErrorMsg()).Error; err != nil {
			log.Printf("id : %v | Failed to update record.", r.GetUserID())
		}
	}
}
