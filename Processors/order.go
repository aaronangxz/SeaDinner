package Processors

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

//Deprecated
func OrderDinnerQuery(client resty.Client, ID int) {
	var req OrderRequest

	//reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Selection: ")
	//selection, _ := reader.ReadString('\n')
	_, err := fmt.Scanf("%d", &req.FoodID)

	if err != nil {
		fmt.Println(err)
	}
	//Call orderdinner with request struct
	OrderDinner(client, ID, req.FoodID, "")
}

func OrderDinner(client resty.Client, menuID int, choice int, key string) bool {
	var resp OrderResponse

	fData := make(map[string]string)
	fData["food_id"] = fmt.Sprint(choice)

	for i := 0; i < 3; i++ {
		log.Println("Attempt", i)

		_, err := client.R().
			SetHeader("Authorization", MakeToken(key)).
			SetFormData(fData).
			SetResult(&resp).
			Post(MakeURL(URL_CURRENT, &menuID))

		if err != nil {
			log.Println(err)
			continue
		}

		if resp.Status != "" {
			log.Printf("%s: %v\n", resp.Status, resp.Error)
			continue
		}

		if resp.Selected != 0 {
			log.Printf("Dinner Selected: %d\n", resp.Selected)
			break
		}
		log.Printf("Dinner Not Selected. Retrying.\n")
	}
	return resp.Selected != 0
}

func BatchOrderDinner(u []UserRecords) {
	var (
		m = make(map[int64]bool)
	)

	for _, r := range u {
		log.Println("Ordering for user_id:", r.UserID)
		m[r.UserID] = OrderDinner(Client, GetDayId(r.Key), r.Choice, r.Key)
	}

	OutputResults(m)
}
