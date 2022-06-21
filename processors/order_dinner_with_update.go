package processors

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"os"
	"time"

	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/go-resty/resty/v2"
	"google.golang.org/protobuf/proto"
)

//OrderDinnerWithUpdate Orders dinner by calling Sea API, retry times determined in Config.
//Instantly sends notifications to user if order is successful.
func OrderDinnerWithUpdate(ctx context.Context, u *sea_dinner.UserChoiceWithKey) (int64, *sea_dinner.OrderRecord) {
	var (
		status  int64
		resp    sea_dinner.OrderResponse
		apiResp *resty.Response
		err     error
	)
	txn := App.StartTransaction("order_dinner_with_update")
	defer txn.End()

	if os.Getenv("TEST_DEPLOY") == "TRUE" || common.Config.Adhoc {
		log.Info(ctx, "OrderDinnerWithUpdate | TEST | return dummy result.")
		return int64(sea_dinner.OrderStatus_ORDER_STATUS_OK), &sea_dinner.OrderRecord{
			UserId:    proto.Int64(u.GetUserId()),
			FoodId:    proto.String(u.GetUserChoice()),
			OrderTime: proto.Int64(time.Now().Unix()),
			Status:    proto.Int64(int64(sea_dinner.OrderStatus_ORDER_STATUS_OK)),
			ErrorMsg:  proto.String("TEST"),
		}
	}

	fData := make(map[string]string)
	fData["food_id"] = fmt.Sprint(u.GetUserChoice())

	for i := 1; i <= common.Config.Runtime.RetryTimes; i++ {
		log.Info(ctx, "id: %v | OrderDinner | Attempt %v", u.GetUserId(), i)
		apiResp, err = Client.R().
			SetHeader("Authorization", MakeToken(ctx, fmt.Sprint(u.GetUserKey()))).
			SetFormData(fData).
			SetResult(&resp).
			EnableTrace().
			Post(MakeURL(int(sea_dinner.URLType_URL_ORDER), proto.Int64(GetDayID(ctx))))

		if err != nil {
			log.Error(ctx, err.Error())
			continue
		}

		if resp.Status != nil && resp.GetStatus() == "error" {
			log.Error(ctx, "id: %v | %v : %v : %v : %v", u.GetUserId(), resp.GetError(), resp.GetStatus(), resp.GetStatusCode(), resp.GetSelected())
		}

		if resp.GetSelected() != 0 {
			log.Info(ctx, "id: %v | Dinner Selected: %d. Successful in %v try.\n", u.GetUserId(), resp.GetSelected(), i)
			break
		}
		log.Error(ctx, "id: %v | Dinner Not Selected. Retrying.\n", u.GetUserId())
	}

	record := &sea_dinner.OrderRecord{
		UserId:    proto.Int64(u.GetUserId()),
		FoodId:    proto.String(u.GetUserChoice()),
		OrderTime: proto.Int64(time.Now().Unix()),
	}

	if resp.GetSelected() == 0 {
		status = int64(sea_dinner.OrderStatus_ORDER_STATUS_FAIL)
		if resp.Error == nil {
			record.ErrorMsg = proto.String("Unknown Error")
		} else {
			record.ErrorMsg = proto.String(resp.GetError())
		}
	} else {
		status = int64(sea_dinner.OrderStatus_ORDER_STATUS_OK)
		timeTaken := apiResp.Request.TraceInfo().TotalTime.Milliseconds()
		SendInstantNotification(ctx, u, timeTaken)
		record.TimeTaken = proto.Int64(timeTaken)
	}
	record.Status = proto.Int64(status)
	return status, record
}
