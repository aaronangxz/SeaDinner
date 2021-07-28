package Processors

import (
	"fmt"
	"strconv"

	"github.com/aaronangxz/SeaDinner/AuthToken"
	"github.com/go-resty/resty/v2"
)

func GetMenu(client resty.Client, ID int) {
	var currentmenu DinnerMenu
	IDstr := strconv.Itoa(ID)
	//var currentarr []DinnerMenu

	_, err := client.R().
		SetHeader("Authorization", AuthToken.GetToken()).
		SetResult(&currentmenu).
		Get("https://dinner.sea.com/api/menu/" + IDstr)
	//currentarr = append(currentarr, currentmenu)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Query status of today's menu: %v\n\n", currentmenu.Status)
	fmt.Printf("ID: %v\n", currentmenu.Dishes.Code)
	// fmt.Printf("%v\n", currentmenu.Details.Name)
	// //fmt.Printf("%v\n", currentmenu.Details.Comment)
	// fmt.Printf("Start: %v\n", currentmenu.Details.PollStart)
	// fmt.Printf("End: %v\n", currentmenu.Details.PollEnd)
	// fmt.Printf("Serving Time: %v\n", currentmenu.Details.ServingTime)
}
