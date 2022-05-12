package Processors

import (
	"fmt"
	"log"
)

func GetDayId(key string) (ID int) {
	var currentmenu Current
	log.Println("header:", MakeToken(key))
	log.Println("url:", MakeURL(URL_CURRENT, nil))

	_, err := Client.R().
		SetHeader("Authorization", MakeToken(key)).
		SetResult(&currentmenu).
		Get(MakeURL(URL_CURRENT, nil))

	if err != nil {
		fmt.Println(err)
	}
	return currentmenu.Details.GetId()
}
