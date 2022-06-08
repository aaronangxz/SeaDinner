package Processors

import (
	"fmt"
	"log"

	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"google.golang.org/protobuf/proto"
)

func GetSuccessfulOrder(key string) bool {
	var (
		order *sea_dinner.DinnerMenu
	)

	_, err := Client.R().
		SetHeader("Authorization", MakeToken(key)).
		SetResult(&order).
		Get(MakeURL(int(sea_dinner.URLType_URL_ORDER), proto.Int64(GetDayId())))

	if err != nil {
		fmt.Println(err)
	}

	if order.GetStatus() != "success" {
		log.Println("GetSuccessfulOrder | Error:", order.GetStatus())
		return false
	}
	return order.Food != nil
}
