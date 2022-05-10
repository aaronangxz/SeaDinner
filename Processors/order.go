package Processors

import (
	"fmt"

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

	_, err := client.R().
		SetHeader("Authorization", MakeToken(key)).
		SetFormData(fData).
		SetResult(&resp).
		Post(MakeURL(URL_CURRENT, &menuID))

	if err != nil {
		fmt.Println(err)
		return false
	}

	if resp.Error != "" {
		fmt.Printf("%s: %s\n", resp.Status, resp.Error)
		return false
	}
	fmt.Printf("Dinner Selected: %d\n", resp.Selected)
	return true

}

func BatchOrderDinner(u []UserRecords) {
	var (
		m = make(map[int64]bool)
	)

	for _, r := range u {
		m[r.UserID] = OrderDinner(Client, GetDayId(r.Key), r.Choice, r.Key)
	}
}
