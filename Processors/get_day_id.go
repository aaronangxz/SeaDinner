package Processors

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

func GetDayId(client resty.Client) (ID int) {
	var currentmenu Current

	_, err := client.R().
		SetHeader("Authorization", "Token "+os.Getenv("Token")).
		SetResult(&currentmenu).
		Get(UrlPrefix + "/api/current")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Query status of today's menu: %v\n\n", currentmenu.Status)
	fmt.Printf("Day ID: %v\n", currentmenu.Details.Id)
	fmt.Printf("%v\n", currentmenu.Details.Name)
	return currentmenu.Details.Id
}
