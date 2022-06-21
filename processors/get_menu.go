package processors

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"google.golang.org/protobuf/proto"
)

//GetMenu Calls Sea API, retrieves the current day's menu in realtime
func GetMenu(ctx context.Context, key string) *sea_dinner.DinnerMenuArray {
	var (
		currentArr *sea_dinner.DinnerMenuArray
	)
	txn := App.StartTransaction("get_menu")
	defer txn.End()

	_, err := Client.R().
		SetHeader("Authorization", MakeToken(ctx, key)).
		SetResult(&currentArr).
		Get(MakeURL(int(sea_dinner.URLType_URL_MENU), proto.Int64(GetDayID(ctx))))
	if err != nil {
		log.Error(ctx, err.Error())
	}
	log.Info(ctx, "GetMenu | Query status of today's menu: %v", currentArr.GetStatus())
	return currentArr
}
