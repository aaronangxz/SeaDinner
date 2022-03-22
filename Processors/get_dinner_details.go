package Processors

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

func GetCurrent(client resty.Client) (ID int) {
	var currentmenu Current

	_, err := client.R().
		SetHeader("Authorization", "Token "+os.Getenv("Token")).
		SetResult(&currentmenu).
		Get("https://dinner.sea.com/api/current")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Query status of today's menu: %v\n\n", currentmenu.Status)
	fmt.Printf("Day ID: %v\n", currentmenu.Details.Id)
	fmt.Printf("%v\n", currentmenu.Details.Name)
	//fmt.Printf("%v\n", currentmenu.Details.Comment)
	//fmt.Printf("Start: %v\n", currentmenu.Details.PollStart)
	//fmt.Printf("End: %v\n", currentmenu.Details.PollEnd)
	//fmt.Printf("Serving Time: %v\n", currentmenu.Details.ServingTime)
	return currentmenu.Details.Id
}
