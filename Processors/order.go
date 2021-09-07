package Processors

import (
	"fmt"
	"strconv"

	"github.com/aaronangxz/SeaDinner/AuthToken"
	"github.com/go-resty/resty/v2"
)

func OrderDinnerQuery(client resty.Client, ID int) {
	var req OrderRequest

	//reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Selection: ")
	//selection, _ := reader.ReadString('\n')
	_, err := fmt.Scanf("%d", &req.FoodID)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Chosen: %d\n", req.FoodID)
	//Call orderdinner with request struct
	OrderDinner(client, ID, req)
}

func OrderDinner(client resty.Client, menuID int, choice OrderRequest) {
	//convert ID to string
	menuIDstr := strconv.Itoa(menuID)
	//url := "https://dinner.sea.com/menu/" + menuIDstr + "/make_order"
	url := "https://dinner.sea.com/api/order/" + menuIDstr

	var resp OrderResponse

	//fmt.Println("link:", url)

	_, err := client.R().
		SetHeader("Authorization", AuthToken.GetToken()).
		SetBody(OrderRequest{FoodID: choice.FoodID}).
		//SetBody({"food_id":1374}).
		SetResult(&resp).
		Post(url)

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(OrderRequest{FoodID: choice.FoodID})

	if resp.Error != "" {
		fmt.Printf("%s: %s\n", resp.Status, resp.Error)
	} else {
		fmt.Printf("Code: %s\n", resp.StatusCode)
		fmt.Printf("Selected: %d\n", resp.Selected)
	}

}
