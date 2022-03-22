package Processors

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
)

func GetMenu(client resty.Client, ID int) {
	IDstr := strconv.Itoa(ID)
	var currentarr DinnerMenuArr

	_, err := client.R().
		SetHeader("Authorization", "Token "+os.Getenv("Token")).
		SetResult(&currentarr).
		Get("https://dinner.sea.com/api/menu/" + IDstr)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Query status of today's menu: %v\n\n", currentarr.Status)

	for i := range currentarr.DinnerArr {
		fmt.Printf("Food ID: %v\n", currentarr.DinnerArr[i].Id)
		fmt.Printf("Name: %v\n", currentarr.DinnerArr[i].Name)
		fmt.Printf("Ordered: %v\n", currentarr.DinnerArr[i].Ordered)
		fmt.Printf("Quota: %v\n", currentarr.DinnerArr[i].Quota)
		fmt.Printf("\n")
	}
}
