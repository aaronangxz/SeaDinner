package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"google.golang.org/protobuf/proto"
)

//CancelOrder Cancels the user's order after it is processed
func CancelOrder(ctx context.Context, id int64) (string, bool) {
	var (
		resp *sea_dinner.OrderResponse
	)
	txn := processors.App.StartTransaction("cancel_order")
	defer txn.End()

	//Get currently ordered food id
	currOrder, ok := processors.GetOrderByUserID(ctx, id)
	if !ok {
		return currOrder, false
	}

	fData := make(map[string]string)
	fData["food_id"] = currOrder

	_, err := processors.Client.R().
		SetHeader("Authorization", processors.MakeToken(ctx, fmt.Sprint(GetKey(ctx, id)))).
		SetFormData(fData).
		SetResult(&resp).
		EnableTrace().
		Delete(processors.MakeURL(int(sea_dinner.URLType_URL_ORDER), proto.Int64(processors.GetDayID(ctx))))

	if err != nil {
		log.Error(ctx, "CancelOrder | error: %v", err.Error())
		return "There were some issues 😥 Try to cancel from SeaTalk instead!", false
	}

	if resp.GetStatus() == "error" {
		log.Error(ctx, "CancelOrder | status error: %v", resp.GetError())
		return fmt.Sprintf("I can't cancel this order: %v 😥 Try to cancel from SeaTalk instead!", resp.GetError()), false
	}

	if resp.Selected != nil {
		log.Error(ctx, "CancelOrder | failed to cancel order")
		return "It seems like you ordered something else 😥 Try to cancel from SeaTalk instead!", false
	}
	log.Info(ctx, "CancelOrder | Success | user_id:%v", id)
	return "I have cancelled your order!😀", true
}
