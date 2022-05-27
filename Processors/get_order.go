package Processors

import (
	"fmt"
	"log"
)

func GetSuccessfulOrder(key string) bool {
	var (
		order UserOrder
	)

	_, err := Client.R().
		SetHeader("Authorization", MakeToken(key)).
		SetResult(&order).
		Get(MakeURL(URL_ORDER, Int(GetDayId())))

	if err != nil {
		fmt.Println(err)
	}

	if order.GetStatus() != "success" {
		log.Println("GetSuccessfulOrder | Error:", order.GetError())
		return false
	}
	return order.Food != nil
}
