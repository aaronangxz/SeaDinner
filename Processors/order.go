package Processors

import (
	"fmt"
	"os"
	"strconv"

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
	//Call orderdinner with request struct
	OrderDinner(client, ID, req)
}

func OrderDinner(client resty.Client, menuID int, choice OrderRequest) {
	//convert ID to string
	menuIDstr := strconv.Itoa(menuID)
	//url := "https://dinner.sea.com/api/order/" + menuIDstr
	url := "https://dinner.sea.com/menu/" + menuIDstr + "/make_order"
	fmt.Println("url:", url)
	fmt.Println("choice:", choice.FoodID)

	var req OrderRequest
	var resp OrderResponse

	req.FoodID = choice.FoodID
	fmt.Println(req)

	_, err := client.R().
		SetHeader("Authorization", "Token "+os.Getenv("Token")).
		SetBody(fmt.Sprintf("food_id=%s", choice.FoodID)).
		SetResult(&resp).
		Post(url)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Error: %s\n", resp.Error)
	//fmt.Println(OrderRequest{FoodID: choice.FoodID})

	if resp.Error != "" {
		fmt.Printf("%s: %s\n", resp.Status, resp.Error)
	} else {
		fmt.Printf("Code: %s\n", resp.StatusCode)
		fmt.Printf("Selected: %d\n", resp.Selected)
	}
}
