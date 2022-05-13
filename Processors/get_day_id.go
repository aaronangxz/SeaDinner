package Processors

import (
	"fmt"
)

func GetDayId(key string) (ID int) {
	var currentmenu Current

	_, err := Client.R().
		SetHeader("Authorization", MakeToken(key)).
		SetResult(&currentmenu).
		Get(MakeURL(URL_CURRENT, nil))

	if err != nil {
		fmt.Println(err)
	}
	return currentmenu.Menu.GetId()
}
