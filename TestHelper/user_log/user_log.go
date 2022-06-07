package user_log

import (
	"fmt"
	"log"
	"time"

	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/Processors"
	"github.com/aaronangxz/SeaDinner/TestHelper"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"google.golang.org/protobuf/proto"
)

var (
	defaultOrderTime    = time.Now().Unix()
	defaultStatus       = int64(sea_dinner.OrderStatus_ORDER_STATUS_OK)
	defaultErrorMessage = "Default Error Message"
)

type UserLog struct {
	*sea_dinner.OrderRecord
}

func New() *UserLog {
	TestHelper.InitTest()
	return &UserLog{
		OrderRecord: &sea_dinner.OrderRecord{
			Id:        new(int64),
			UserId:    new(int64),
			FoodId:    new(string),
			OrderTime: new(int64),
			Status:    new(int64),
			ErrorMsg:  new(string),
		},
	}
}

func (ul *UserLog) FillDefaults() *UserLog {
	if ul.OrderRecord.UserId == nil {
		ul.SetUserId(TestHelper.RandomInt(99999))
	}

	if ul.OrderRecord.FoodId == nil {
		ul.SetFoodId(fmt.Sprint(TestHelper.RandomInt(99999)))
	}

	if ul.OrderRecord.OrderTime == nil {
		ul.SetOrderTime(defaultOrderTime)
	}

	if ul.OrderRecord.Status == nil {
		ul.SetStatus(defaultStatus)
	}

	if ul.OrderRecord.GetStatus() == int64(sea_dinner.OrderStatus_ORDER_STATUS_FAIL) && ul.OrderRecord.ErrorMsg == nil {
		ul.SetErrorMsg(defaultErrorMessage)
	}

	return ul
}

func (ul *UserLog) Build() *UserLog {
	ul.FillDefaults()
	if err := Processors.DB.Table(Common.DB_ORDER_LOG_TAB).Create(&ul).Error; err != nil {
		log.Printf("Failed to insert to DB | user_id:%v | %v", ul.GetUserId(), err.Error())
		return nil
	}
	log.Printf("Sulcessfully inserted to DB | user_id:%v", ul.GetUserId())
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

func (ul *UserLog) SetStatus(status int64) *UserLog {
	ul.OrderRecord.Status = proto.Int64(status)
	return ul
}

func (ul *UserLog) SetErrorMsg(errorMsg string) *UserLog {
	ul.OrderRecord.ErrorMsg = proto.String(errorMsg)
	return ul
}

func (ul *UserLog) TearDown() error {
	if err := Processors.DB.Exec("DELETE FROM user_log_tab WHERE user_id = ?", ul.GetUserId()).Error; err != nil {
		log.Printf("Failed to delete from DB | user_id:%v", ul.GetUserId())
		return err
	}
	log.Printf("Sulcessfully deleted from DB | user_id:%v", ul.GetUserId())
	return nil
}

func DeleteUserKey(userId int64) error {
	if err := Processors.DB.Exec("DELETE FROM user_log_tab WHERE user_id = ?", userId).Error; err != nil {
		log.Printf("Failed to delete from DB | user_id:%v", userId)
		return err
	}
	log.Printf("Sulcessfully deleted from DB | user_id:%v", userId)
	return nil
}
