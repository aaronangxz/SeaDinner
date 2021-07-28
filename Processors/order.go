package Processors

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/aaronangxz/SeaDinner/AuthToken"
	"github.com/go-resty/resty/v2"
)

func OrderDinnerQuery(client resty.Client, ID int) {
	var req OrderRequest

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Selection: ")
	req.FoodID, _ = reader.ReadString('\n')

	//Call orderdinner with request struct
	OrderDinner(client, ID, req)
}

func OrderDinner(client resty.Client, ID int, choice OrderRequest) {
	//convert ID to string
	IDstr := strconv.Itoa(ID)

	var resp OrderResponse

	_, err := client.R().
		SetHeader("Authorization", AuthToken.GetToken()).
		SetBody(OrderRequest{FoodID: choice.FoodID}).
		SetResult(&resp).
		Get("https://dinner.sea.com/api/order/" + IDstr)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Status: %s\n", resp.Status)

	if resp.Error != "" {
		fmt.Printf("Error: %s\n", resp.Error)

	}

	fmt.Printf("Selected: %d\n", resp.Selected)
}
