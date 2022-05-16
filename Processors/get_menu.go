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

	if ID == 0 {
		log.Println("GetMenu | Invalid id:", ID)
		return currentarr
	}

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

	m := GetMenu(Client, GetDayId(key), key)

	if m.Status == nil {
		return "There is no dinner order today!"
	}

	for _, d := range m.DinnerArr {
		output += fmt.Sprintf(Config.Prefix.UrlPrefix+"%v\nFood ID: %v\nName: %v\nQuota: %v\n\n",
			d.ImageURL, d.Id, d.Name, d.Quota)
	}
	return output
}
