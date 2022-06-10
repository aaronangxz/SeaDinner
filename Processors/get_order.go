package Processors

import (
	"github.com/aaronangxz/SeaDinner/Log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"google.golang.org/protobuf/proto"
)

//GetSuccessfulOrder Calls Sea API, check the actual order of a user. Returns food = nil if there is no order.
func GetSuccessfulOrder(key string) bool {
	var (
		order *sea_dinner.DinnerMenu
	)

	_, err := Client.R().
		SetHeader("Authorization", MakeToken(Ctx, key)).
		SetResult(&order).
		Get(MakeURL(int(sea_dinner.URLType_URL_ORDER), proto.Int64(GetDayId(Ctx))))

	if err != nil {
		Log.Error(Ctx, err.Error())
		// fmt.Println(err)
	}

	if order.GetStatus() != "success" {
		Log.Error(Ctx, "GetSuccessfulOrder | Error: %v", order.GetStatus())
		// log.Println("GetSuccessfulOrder | Error:", order.GetStatus())
		return false
	}
	return order.Food != nil
}
