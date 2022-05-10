package Processors

import (
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

func GetMenu(client resty.Client, ID int, key string) DinnerMenuArr {
	IDstr := strconv.Itoa(ID)
	var currentarr DinnerMenuArr

	_, err := client.R().
		SetHeader("Authorization", "Token "+key).
		SetResult(&currentarr).
		Get(UrlPrefix + "/api/menu/" + IDstr)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Query status of today's menu: %v\n\n", currentarr.Status)

	return currentarr
}

func OutputMenu(key string) string {
	menu := GetMenu(Client, GetDayId(Client), key)
	output := ""

	for i := range menu.DinnerArr {
		output += fmt.Sprintf(UrlPrefix+"%v\nFood ID: %v\nName: %v\nQuota: %v\n\n",
			menu.DinnerArr[i].ImageURL, menu.DinnerArr[i].Id, menu.DinnerArr[i].Name, menu.DinnerArr[i].Quota)
	}
	return output
}
