package Processors

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func GetMenu(client resty.Client, ID int, key string) DinnerMenuArr {
	var currentarr DinnerMenuArr

	_, err := client.R().
		SetHeader("Authorization", MakeToken(key)).
		SetResult(&currentarr).
		Get(MakeURL(URL_CURRENT, &ID))

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Query status of today's menu: %v\n\n", currentarr.Status)

	return currentarr
}

func OutputMenu(key string) string {
	menu := GetMenu(Client, GetDayId(key), key)
	output := ""

	for i := range menu.DinnerArr {
		output += fmt.Sprintf(UrlPrefix+"%v\nFood ID: %v\nName: %v\nQuota: %v\n\n",
			menu.DinnerArr[i].ImageURL, menu.DinnerArr[i].Id, menu.DinnerArr[i].Name, menu.DinnerArr[i].Quota)
	}
	return output
}
