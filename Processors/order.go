package Processors

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func OrderDinner(client resty.Client, menuID int, u UserChoiceWithKey) bool {
	var resp OrderResponse
	fData := make(map[string]string)
	fData["food_id"] = fmt.Sprint(u.Choice)

	for i := 1; i < Config.Runtime.RetryTimes; i++ {
		log.Printf("id: %v | Attempt %v", u.UserID, i)

		_, err := client.R().
			SetHeader("Authorization", MakeToken(u.Key)).
			SetFormData(fData).
			SetResult(&resp).
			Post(MakeURL(URL_ORDER, &menuID))

		if err != nil {
			log.Println(err)
			continue
		}

		if resp.Error != nil && resp.GetError() != "success" {
			log.Printf("id: %v : | %v: %v\n", u.UserID, resp.GetStatus(), resp.GetError())
			continue
		}

		if resp.GetSelected() != 0 {
			log.Printf("id: %v | Dinner Selected: %d. Successful in %v try.\n", u.UserID, i, resp.GetSelected())
			break
		}
		log.Printf("id: %v | Dinner Not Selected. Retrying.\n", u.UserID)
	}
	return resp.GetSelected() != 0
}

func BatchOrderDinner(u []UserChoiceWithKey) {
	var (
		m = make(map[int64]bool)
	)

	for _, r := range u {
		log.Printf("id: %v | Ordering\n", r.UserID)
		m[r.UserID] = OrderDinner(Client, GetDayId(r.Key), r)
	}

	OutputResults(m)
}
