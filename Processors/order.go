package Processors

import (
	"fmt"
	"strconv"

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
	//convert ID to string
	menuIDstr := strconv.Itoa(menuID)
	url := "https://dinner.sea.com/api/order/" + menuIDstr
	//url := "https://dinner.sea.com/menu/" + menuIDstr + "/make_order"

	fmt.Println("url:", url)
	fmt.Println("choice:", choice)

	var resp OrderResponse

	fData := make(map[string]string)
	fData["food_id"] = fmt.Sprint(choice)

	_, err := client.R().
		SetHeader("Authorization", "Token "+key).
		SetFormData(fData).
		SetResult(&resp).
		Post(url)

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
		m[r.UserID] = OrderDinner(Client, GetDayId(Client), r.Choice, r.Key)
	}
}
