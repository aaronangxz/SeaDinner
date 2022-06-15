package processors

import (
	"context"

	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"google.golang.org/protobuf/proto"
)

//GetSuccessfulOrder Calls Sea API, check the actual order of a user. Returns food = nil if there is no order.
func GetSuccessfulOrder(ctx context.Context, key string) bool {
	var (
		order *sea_dinner.DinnerMenu
	)

	_, err := Client.R().
		SetHeader("Authorization", MakeToken(ctx, key)).
		SetResult(&order).
		Get(MakeURL(int(sea_dinner.URLType_URL_ORDER), proto.Int64(GetDayID(ctx))))

	if err != nil {
		log.Error(ctx, err.Error())
	}

	if order.GetStatus() != "success" {
		log.Error(ctx, "GetSuccessfulOrder | Error: %v", order.GetStatus())
		return false
	}
	return order.Food != nil
}
