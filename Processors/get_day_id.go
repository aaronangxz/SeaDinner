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

	fmt.Printf("Query status of today's menu: %v\n\n", currentmenu.Status)
	fmt.Printf("Day ID: %v\n", currentmenu.Details.Id)
	fmt.Printf("%v\n", currentmenu.Details.Name)
	return currentmenu.Details.Id
}
