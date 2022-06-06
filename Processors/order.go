package Processors

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/proto"
)

//OrderDinnerWithUpdate Orders dinner by calling Sea API, retry times determined in Config.
//Instantly sends notifications to user if order is successful.
func OrderDinnerWithUpdate(u *sea_dinner.UserChoiceWithKey) (int64, *sea_dinner.OrderRecord) {
	var (
		status  int64
		resp    sea_dinner.OrderResponse
		apiResp *resty.Response
		err     error
	)

	if os.Getenv("TEST_DEPLOY") == "TRUE" || Common.Config.Adhoc {
		log.Println("OrderDinnerWithUpdate | TEST | return dummy result.")
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
		log.Printf("id: %v | OrderDinner | Attempt %v", u.GetUserId(), i)
		apiResp, err = Client.R().
			SetHeader("Authorization", MakeToken(fmt.Sprint(u.GetUserKey()))).
			SetFormData(fData).
			SetResult(&resp).
			EnableTrace().
			Post(MakeURL(int(sea_dinner.URLType_URL_ORDER), proto.Int64(GetDayId())))

		if err != nil {
			log.Println(err)
			continue
		}

		if resp.Status != nil && resp.GetStatus() == "error" {
			log.Printf("id: %v | %v : %v : %v : %v", u.GetUserId(), resp.GetError(), resp.GetStatus(), resp.GetStatusCode(), resp.GetSelected())
		}

		if resp.GetSelected() != 0 {
			log.Printf("id: %v | Dinner Selected: %d. Successful in %v try.\n", u.GetUserId(), resp.GetSelected(), i)
			break
		}
		log.Printf("id: %v | Dinner Not Selected. Retrying.\n", u.GetUserId())
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
		record.Status = proto.Int64(int64(status))
	} else {
		status = int64(sea_dinner.OrderStatus_ORDER_STATUS_OK)
		record.Status = proto.Int64(int64(status))

		trace := apiResp.Request.TraceInfo()
		SendInstantNotification(u, trace.TotalTime.Milliseconds())
	}
	return status, record
}

//BatchOrderDinnerMultiThreaded Spawns multiple Order goroutines, and update order_log_tab with the respective results.
//Guranteed to execute goroutines for all users in the queue.
func BatchOrderDinnerMultiThreaded(userQueue []*sea_dinner.UserChoiceWithKey) {
	var (
		wg      sync.WaitGroup
		records []*sea_dinner.OrderRecord
	)
	txn := App.StartTransaction("batch_order_dinner_multi_threaded")
	defer txn.End()
	m := make(map[int64]int64)
	log.Printf("BatchOrderDinnerMultiThreaded | Begin | size: %v", len(userQueue))

	for _, user := range userQueue {
		if Common.IsInGrayScale(user.GetUserId()) {
			log.Printf("BatchOrderDinnerMultiThreaded | In grayscale, skipping | user_id:%v", user.GetUserId())
			continue
		}
		//Increment group
		wg.Add(1)
		go func(u *sea_dinner.UserChoiceWithKey) {
			//Release group
			defer wg.Done()
			var record *sea_dinner.OrderRecord
			m[u.GetUserId()], record = OrderDinnerWithUpdate(u)
			records = append(records, record)
		}(user)
	}

	//Wait for all groups to release
	wg.Wait()

	log.Printf("BatchOrderDinnerMultiThreaded | Done")
	UpdateOrderLog(records)
	OutputResults(m, "BatchOrderDinnerMultiThreaded")
}

func BatchOrderDinnerMultiThreadedWithWait(userQueue []*sea_dinner.UserChoiceWithKey) {
	var (
		wg      sync.WaitGroup
		records []*sea_dinner.OrderRecord
	)
	txn := App.StartTransaction("batch_order_dinner_multi_threaded_with_wait")
	defer txn.End()
	m := make(map[int64]int64)

	for _, user := range userQueue {
		if !Common.IsInGrayScale(user.GetUserId()) {
			log.Printf("BatchOrderDinnerMultiThreadedWithWait | Not in grayscale, skipping | user_id:%v", user.GetUserId())
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
					log.Printf("BatchOrderDinnerMultiThreadedWithWait | Begin | user_id: %v", u.GetUserId())
					m[u.GetUserId()], record = OrderDinnerWithUpdate(u)
					records = append(records, record)
					break
				}
			}
		}(user)
	}

	//Wait for all groups to release
	wg.Wait()

	log.Printf("BatchOrderDinnerMultiThreadedWithWait | Done")
	UpdateOrderLog(records)
	OutputResults(m, "BatchOrderDinnerMultiThreadedWithWait")
}

//UpdateOrderLog Batch insert new order records into order_log_tab
func UpdateOrderLog(records []*sea_dinner.OrderRecord) {
	if records == nil {
		log.Printf("UpdateOrderLog | No record to update.")
		return
	}
	if err := DB.Table(Common.DB_ORDER_LOG_TAB).Create(&records).Error; err != nil {
		log.Printf("UpdateOrderLog | Failed to update records.")
	}
}

//SendInstantNotification Spawns a one-time telegram bot instance and send notification to user
func SendInstantNotification(u *sea_dinner.UserChoiceWithKey, took int64) {
	var (
		msg string
	)
	bot, err := tgbotapi.NewBotAPI(Common.GetTGToken())
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	menu := MakeMenuMap()
	msg = fmt.Sprintf("Successfully ordered %v in %vms! 🥳", menu[u.GetUserChoice()], took)

	if _, err := bot.Send(tgbotapi.NewMessage(u.GetUserId(), msg)); err != nil {
		log.Println(err)
	}
	log.Printf("SendInstantNotification | user_id:%v | msg: %v", u.GetUserId(), msg)
}

//MakeMenuMap Returns food_id:food_name mapping of current menu
func MakeMenuMap() map[string]string {
	var (
		key = os.Getenv("TOKEN")
	)
	menuMap := make(map[string]string)
	menu := GetMenu(Client, key)
	for _, m := range menu.GetFood() {
		menuMap[fmt.Sprint(m.GetId())] = m.GetName()
	}
	// Store -1 hash to menuMap
	menuMap["-1"] = "*NOTHING*" // to be renamed
	menuMap["RAND"] = "RAND"
	return menuMap
}
