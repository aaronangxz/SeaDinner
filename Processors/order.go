package Processors

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/Log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	if os.Getenv("TEST_DEPLOY") == "TRUE" || Common.Config.Adhoc {
		Log.Info(ctx, "OrderDinnerWithUpdate | TEST | return dummy result.")
		// log.Println("OrderDinnerWithUpdate | TEST | return dummy result.")
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

	for i := 1; i <= Common.Config.Runtime.RetryTimes; i++ {
		Log.Info(ctx, "id: %v | OrderDinner | Attempt %v", u.GetUserId(), i)
		// log.Printf("id: %v | OrderDinner | Attempt %v", u.GetUserId(), i)
		apiResp, err = Client.R().
			SetHeader("Authorization", MakeToken(ctx, fmt.Sprint(u.GetUserKey()))).
			SetFormData(fData).
			SetResult(&resp).
			EnableTrace().
			Post(MakeURL(int(sea_dinner.URLType_URL_ORDER), proto.Int64(GetDayId(ctx))))

		if err != nil {
			Log.Error(ctx, err.Error())
			// log.Println(err)
			continue
		}

		if resp.Status != nil && resp.GetStatus() == "error" {
			Log.Error(ctx, "id: %v | %v : %v : %v : %v", u.GetUserId(), resp.GetError(), resp.GetStatus(), resp.GetStatusCode(), resp.GetSelected())
			// log.Printf("id: %v | %v : %v : %v : %v", u.GetUserId(), resp.GetError(), resp.GetStatus(), resp.GetStatusCode(), resp.GetSelected())
		}

		if resp.GetSelected() != 0 {
			Log.Info(ctx, "id: %v | Dinner Selected: %d. Successful in %v try.\n", u.GetUserId(), resp.GetSelected(), i)
			// log.Printf("id: %v | Dinner Selected: %d. Successful in %v try.\n", u.GetUserId(), resp.GetSelected(), i)
			break
		}
		Log.Error(ctx, "id: %v | Dinner Not Selected. Retrying.\n", u.GetUserId())
		// log.Printf("id: %v | Dinner Not Selected. Retrying.\n", u.GetUserId())
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

//BatchOrderDinnerMultiThreaded Spawns multiple Order goroutines, and update order_log_tab with the respective results.
//Guranteed to execute goroutines for all users in the queue.
func BatchOrderDinnerMultiThreaded(ctx context.Context, userQueue []*sea_dinner.UserChoiceWithKey) {
	var (
		wg      sync.WaitGroup
		records []*sea_dinner.OrderRecord
	)
	txn := App.StartTransaction("batch_order_dinner_multi_threaded")
	defer txn.End()

	m := make(map[int64]int64)
	Log.Info(ctx, "BatchOrderDinnerMultiThreaded | Begin | size: %v", len(userQueue))
	// log.Printf("BatchOrderDinnerMultiThreaded | Begin | size: %v", len(userQueue))

	for _, user := range userQueue {
		if Common.IsInGrayScale(user.GetUserId()) {
			Log.Info(ctx, "BatchOrderDinnerMultiThreaded | In grayscale, skipping | user_id:%v", user.GetUserId())
			// log.Printf("BatchOrderDinnerMultiThreaded | In grayscale, skipping | user_id:%v", user.GetUserId())
			continue
		}
		//Increment group
		wg.Add(1)
		go func(u *sea_dinner.UserChoiceWithKey) {
			//Release group
			defer wg.Done()
			var record *sea_dinner.OrderRecord
			m[u.GetUserId()], record = OrderDinnerWithUpdate(ctx, u)
			records = append(records, record)
		}(user)
	}

	//Wait for all groups to release
	wg.Wait()

	Log.Info(ctx, "BatchOrderDinnerMultiThreaded | Done")
	// log.Printf("BatchOrderDinnerMultiThreaded | Done")
	BatchInsertOrderLogs(ctx, records)
	OutputResults(ctx, m, "BatchOrderDinnerMultiThreaded")
}

func BatchOrderDinnerMultiThreadedWithWait(ctx context.Context, userQueue []*sea_dinner.UserChoiceWithKey) {
	var (
		wg      sync.WaitGroup
		records []*sea_dinner.OrderRecord
	)
	txn := App.StartTransaction("batch_order_dinner_multi_threaded_with_wait")
	defer txn.End()

	m := make(map[int64]int64)

	for _, user := range userQueue {
		if !Common.IsInGrayScale(user.GetUserId()) {
			Log.Info(ctx, "BatchOrderDinnerMultiThreadedWithWait | Not in grayscale, skipping | user_id:%v", user.GetUserId())
			//log.Printf("BatchOrderDinnerMultiThreadedWithWait | Not in grayscale, skipping | user_id:%v", user.GetUserId())
			continue
		}
		//Increment group
		wg.Add(1)
		go func(u *sea_dinner.UserChoiceWithKey) {
			//Release group
			defer wg.Done()
			var record *sea_dinner.OrderRecord
			for {
				if IsOrderTime() && IsPollStart() {
					Log.Info(ctx, "BatchOrderDinnerMultiThreadedWithWait | Begin | user_id: %v", u.GetUserId())
					//log.Printf("BatchOrderDinnerMultiThreadedWithWait | Begin | user_id: %v", u.GetUserId())
					m[u.GetUserId()], record = OrderDinnerWithUpdate(ctx, u)
					records = append(records, record)
					break
				}
			}
		}(user)
	}

	//Wait for all groups to release
	wg.Wait()
	Log.Info(ctx, "BatchOrderDinnerMultiThreadedWithWait | Done")
	//log.Printf("BatchOrderDinnerMultiThreadedWithWait | Done")
	BatchInsertOrderLogs(ctx, records)
	OutputResults(ctx, m, "BatchOrderDinnerMultiThreadedWithWait")
}

//BatchInsertOrderLogs Batch insert new order records into order_log_tab
func BatchInsertOrderLogs(ctx context.Context, records []*sea_dinner.OrderRecord) {
	txn := App.StartTransaction("batch_insert_order_logs")
	defer txn.End()

	if records == nil {
		Log.Warn(ctx, "BatchInsertOrderLogs | No record to update.")
		//log.Printf("BatchInsertOrderLogs | No record to update.")
		return
	}
	if err := DB.Table(Common.DB_ORDER_LOG_TAB).Create(&records).Error; err != nil {
		Log.Error(ctx, fmt.Sprintf("BatchInsertOrderLogs | Failed to update records | %v", err.Error()))
		//log.Printf("BatchInsertOrderLogs | Failed to update records | %v", err.Error())
		return
	}
	Log.Info(ctx, fmt.Sprintf("BatchInsertOrderLogs | Successfully updated records | size: %v", len(records)))
	//log.Printf("BatchInsertOrderLogs | Successfully updated records | size: %v", len(records))
}

//UpdateOrderLog Update a single record in order_log_tab
func UpdateOrderLog(record *sea_dinner.OrderRecord) {
	txn := App.StartTransaction("update_order_log")
	defer txn.End()

	if record == nil {
		log.Printf("UpdateOrderLog | No record to update.")
		return
	}

	if err := DB.Exec("UPDATE user_log_tab SET status = ? WHERE user_id = ?", sea_dinner.OrderStatus_ORDER_STATUS_CANCEL, record.GetUserId()).Error; err != nil {
		log.Printf("UpdateOrderLog | Failed to update records | %v", err.Error())
		return
	}
	log.Printf("UpdateOrderLog | Successfully updated record | user_id: %v", record.GetUserId())
}

//SendInstantNotification Spawns a one-time telegram bot instance and send notification to user
func SendInstantNotification(ctx context.Context, u *sea_dinner.UserChoiceWithKey, took int64) {
	var (
		mk   tgbotapi.InlineKeyboardMarkup
		out  [][]tgbotapi.InlineKeyboardButton
		rows []tgbotapi.InlineKeyboardButton
	)
	txn := App.StartTransaction("send_instant_notifications")
	defer txn.End()

	bot, err := tgbotapi.NewBotAPI(Common.GetTGToken(ctx))
	if err != nil {
		Log.Error(ctx, err.Error())
		// log.Panic(err)
	}
	bot.Debug = true
	Log.Info(ctx, "Authorized on account %s", bot.Self.UserName)
	// log.Printf("Authorized on account %s", bot.Self.UserName)

	menu := MakeMenuMap(ctx)
	msg := tgbotapi.NewMessage(u.GetUserId(), "")
	msg.Text = fmt.Sprintf("Successfully ordered %v in %vms! 🥳", menu[u.GetUserChoice()], took)

	skipBotton := tgbotapi.NewInlineKeyboardButtonData("I DON'T NEED IT 🙅", "ATTEMPTCANCEL")
	rows = append(rows, skipBotton)
	out = append(out, rows)
	mk.InlineKeyboard = out
	msg.ReplyMarkup = mk
	if _, err := bot.Send(msg); err != nil {
		Log.Error(ctx, err.Error())
		// log.Println(err)
	}
	Log.Info(ctx, "SendInstantNotification | user_id:%v | msg: %v", u.GetUserId(), msg)
	// log.Printf("SendInstantNotification | user_id:%v | msg: %v", u.GetUserId(), msg)
}

//MakeMenuMap Returns food_id:food_name mapping of current menu
func MakeMenuMap(ctx context.Context) map[string]string {
	var (
		key = os.Getenv("TOKEN")
	)
	txn := App.StartTransaction("make_menu_map")
	defer txn.End()

	menuMap := make(map[string]string)
	menu := GetMenuUsingCache(ctx, Client, key)
	for _, m := range menu.GetFood() {
		menuMap[fmt.Sprint(m.GetId())] = m.GetName()
	}
	// Store -1 hash to menuMap
	menuMap["-1"] = "*NOTHING*" // to be renamed
	menuMap["RAND"] = "RAND"
	return menuMap
}
