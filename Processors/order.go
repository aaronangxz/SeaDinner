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
	fData["food_id"] = fmt.Sprint(u.Choice)

	log.Println(u)
	for i := 1; i < Config.Runtime.RetryTimes; i++ {
		log.Printf("id: %v | Attempt %v", u.UserID, i)

		start := time.Now().Unix()

		_, err := client.R().
			SetHeader("Authorization", MakeToken(u.Key)).
			SetFormData(fData).
			SetResult(&resp).
			Post(MakeURL(URL_ORDER, &menuID))

		elapsed := time.Now().Unix() - start

		if err != nil {
			log.Println(err)
			continue
		}

		if resp.Status != nil && resp.GetStatus() == "error" {
			log.Printf("id: %v | %v : %v : %v : %v", u.UserID, resp.GetError(), resp.GetStatus(), resp.GetStatusCode(), resp.GetSelected())
		}

		if resp.GetSelected() != 0 {
			log.Printf("id: %v | Dinner Selected: %d. Successful in %v try. Time: %vms\n", u.UserID, resp.GetSelected(), i, elapsed)
			break
		}
		log.Printf("id: %v | Dinner Not Selected. Retrying.\n", u.UserID)
	}
	return resp
}

func BatchOrderDinner(u []UserChoiceWithKey) {
	var (
		records []OrderRecord
		m       = make(map[int64]int)
	)

	for _, r := range u {
		log.Printf("id: %v | Ordering\n", r.UserID)
		resp := OrderDinner(Client, GetDayId(r.Key), r)

		if resp.GetSelected() == 0 {
			m[r.UserID] = ORDER_STATUS_FAIL
		} else {
			m[r.UserID] = ORDER_STATUS_OK
		}

		record := OrderRecord{
			UserID:    Int64(r.UserID),
			FoodID:    Int64(r.Choice),
			OrderTime: Int64(time.Now().Unix()),
			Status:    Int64(int64(m[r.UserID])),
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
		if err := DB.Exec("INSERT INTO order_log (user_id, food_id, order_time, status, error_msg) VALUES (?,?,?,?,?)",
			r.GetUserID(), r.GetFoodID(), r.GetOrderTime(), r.GetStatus(), r.GetErrorMsg()).Error; err != nil {
			log.Printf("id : %v | Failed to update record.", r.GetUserID())
		}
	}
}
