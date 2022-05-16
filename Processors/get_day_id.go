package Processors

import (
	"fmt"
	"log"
	"time"
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

	if currentmenu.Menu.GetPollStart() != fmt.Sprint(ConvertTimeStamp(time.Now().Unix()), "T04:30:00Z") {
		log.Println("GetDayId | Today's ID not found:", currentmenu.Menu.GetPollStart())
		return 0
	}

	return currentmenu.Menu.GetId()
}
