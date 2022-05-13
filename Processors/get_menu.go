package Processors

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func GetMenu(client resty.Client, ID int, key string) DinnerMenuArr {
	var currentarr DinnerMenuArr

	log.Println("key:", key)
	log.Println("header:", MakeToken(key))
	log.Println("url:", MakeURL(URL_MENU, &ID))

	_, err := client.R().
		SetHeader("Authorization", MakeToken(key)).
		SetResult(&currentarr).
		Get(MakeURL(URL_MENU, &ID))

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Query status of today's menu: %v\n\n", currentarr.Status)

	return currentarr
}

func OutputMenu(key string) string {
	var (
		output string
	)

	for _, d := range GetMenu(Client, GetDayId(key), key).DinnerArr {
		output += fmt.Sprintf(Config.Prefix.UrlPrefix+"%v\nFood ID: %v\nName: %v\nQuota: %v\n\n",
			d.ImageURL, d.Id, d.Name, d.Quota)
	}
	return output
}
