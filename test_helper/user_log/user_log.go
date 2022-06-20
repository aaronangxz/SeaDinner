package user_log

import (
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"log"
	"time"

	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/aaronangxz/SeaDinner/test_helper"
	"google.golang.org/protobuf/proto"
)

var (
	defaultOrderTime    = time.Now().Unix()
	defaultTimeTaken    = int64(100)
	defaultStatus       = int64(sea_dinner.OrderStatus_ORDER_STATUS_OK)
	defaultErrorMessage = "Default Error Message"
)

type UserLog struct {
	*sea_dinner.OrderRecord
}

func New() *UserLog {
	test_helper.InitTest()
	return &UserLog{
		OrderRecord: &sea_dinner.OrderRecord{
			UserId:    new(int64),
			FoodId:    new(string),
			OrderTime: new(int64),
			TimeTaken: new(int64),
			Status:    new(int64),
			ErrorMsg:  new(string),
		},
	}
}

func (ul *UserLog) FillDefaults() *UserLog {
	if ul.OrderRecord.GetUserId() == 0 {
		ul.SetUserId(test_helper.RandomInt(99999))
	}

	if ul.OrderRecord.GetFoodId() == "" {
		ul.SetFoodId(fmt.Sprint(test_helper.RandomInt(99999)))
	}

	if ul.OrderRecord.GetOrderTime() == 0 {
		ul.SetOrderTime(defaultOrderTime)
	}

	if ul.OrderRecord.GetTimeTaken() == 0 {
		ul.SetTimeTaken(defaultTimeTaken)
	}
	if ul.OrderRecord.GetStatus() == 0 {
		ul.SetStatus(defaultStatus)
	}

	if ul.OrderRecord.GetStatus() == int64(sea_dinner.OrderStatus_ORDER_STATUS_FAIL) && ul.OrderRecord.GetErrorMsg() == "" {
		ul.SetErrorMsg(defaultErrorMessage)
	}

	return ul
}

func (ul *UserLog) Build() *UserLog {
	ul.FillDefaults()
	if err := processors.DB.Table(common.DB_ORDER_LOG_TAB).Create(&ul).Error; err != nil {
		log.Printf("Failed to insert to DB | user_id:%v | %v", ul.GetUserId(), err.Error())
		return nil
	}
	log.Printf("Sulcessfully inserted to DB | user_id:%v", ul.GetUserId())
	return ul
}

func (ul *UserLog) SetId(id int64) *UserLog {
	ul.OrderRecord.Id = proto.Int64(id)
	return ul
}

func (ul *UserLog) SetUserId(userId int64) *UserLog {
	ul.OrderRecord.UserId = proto.Int64(userId)
	return ul
}

func (ul *UserLog) SetFoodId(foodId string) *UserLog {
	ul.OrderRecord.FoodId = proto.String(foodId)
	return ul
}

func (ul *UserLog) SetOrderTime(orderTime int64) *UserLog {
	ul.OrderRecord.OrderTime = proto.Int64(orderTime)
	return ul
}

func (ul *UserLog) SetTimeTaken(TimeTaken int64) *UserLog {
	ul.OrderRecord.TimeTaken = proto.Int64(TimeTaken)
	return ul
}

func (ul *UserLog) SetStatus(status int64) *UserLog {
	ul.OrderRecord.Status = proto.Int64(status)
	return ul
}

func (ul *UserLog) SetErrorMsg(errorMsg string) *UserLog {
	ul.OrderRecord.ErrorMsg = proto.String(errorMsg)
	return ul
}

func (ul *UserLog) TearDown() error {
	if err := processors.DB.Exec("DELETE FROM order_log_tab WHERE user_id = ?", ul.GetUserId()).Error; err != nil {
		log.Printf("Failed to delete from DB | user_id:%v", ul.GetUserId())
		return err
	}
	log.Printf("Sulcessfully deleted from DB | user_id:%v", ul.GetUserId())
	return nil
}

func DeleteUserKey(userId int64) error {
	if err := processors.DB.Exec("DELETE FROM order_log_tab WHERE user_id = ?", userId).Error; err != nil {
		log.Printf("Failed to delete from DB | user_id:%v", userId)
		return err
	}
	log.Printf("Sulcessfully deleted from DB | user_id:%v", userId)
	return nil
}

func ConvertUserLogToOrderRecord(original *UserLog) *sea_dinner.OrderRecord {
	return &sea_dinner.OrderRecord{
		Id:        proto.Int64(original.GetId()),
		UserId:    proto.Int64(original.GetUserId()),
		FoodId:    proto.String(original.GetFoodId()),
		OrderTime: proto.Int64(original.GetOrderTime()),
		TimeTaken: proto.Int64(original.GetTimeTaken()),
		Status:    proto.Int64(original.GetStatus()),
		ErrorMsg:  proto.String(original.GetErrorMsg()),
	}
}

func CheckOrderLog(userId int64) *sea_dinner.OrderRecord {
	var (
		row *sea_dinner.OrderRecord
	)
	if err := processors.DB.Raw("SELECT * FROM user_log_tab WHERE user_id = ?", userId).Scan(&row).Error; err != nil {
		log.Printf("Failed to read from DB | user_id:%v", userId)
		return nil
	}

	log.Printf("Successfully read from DB | user_id:%v", userId)
	return row
}
