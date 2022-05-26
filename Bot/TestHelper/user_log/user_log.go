package user_log

import (
	"fmt"
	"log"
	"time"

	"github.com/aaronangxz/SeaDinner/Bot/TestHelper"
	"github.com/aaronangxz/SeaDinner/Processors"
)

var (
	defaultOrderTime    = time.Now().Unix()
	defaultStatus       = int64(Processors.ORDER_STATUS_OK)
	defaultErrorMessage = "Default Error Message"
)

type UserLog struct {
	*TestHelper.UserLog
}

func New() *UserLog {
	TestHelper.InitTest()
	return &UserLog{
		UserLog: &TestHelper.UserLog{
			ID:        new(int64),
			UserID:    new(int64),
			FoodID:    new(string),
			OrderTime: new(int64),
			Status:    new(int64),
			ErrorMsg:  new(string),
		},
	}
}

func (ul *UserLog) FillDefaults() *UserLog {
	if ul.UserLog.UserID == nil {
		ul.SetUserId(TestHelper.RandomInt(99999))
	}

	if ul.UserLog.FoodID == nil {
		ul.SetFoodId(fmt.Sprint(TestHelper.RandomInt(99999)))
	}

	if ul.UserLog.OrderTime == nil {
		ul.SetOrderTime(defaultOrderTime)
	}

	if ul.UserLog.Status == nil {
		ul.SetStatus(defaultStatus)
	}

	if ul.UserLog.GetStatus() == Processors.ORDER_STATUS_FAIL && ul.UserLog.ErrorMsg == nil {
		ul.SetErrorMsg(defaultErrorMessage)
	}

	return ul
}

func (ul *UserLog) Build() *UserLog {
	ul.FillDefaults()
	if err := Processors.DB.Table(Processors.DB_ORDER_LOG_TAB).Create(&ul).Error; err != nil {
		log.Printf("Failed to insert to DB | user_id:%v | %v", ul.GetUserID(), err.Error())
		return nil
	}
	log.Printf("Sulcessfully inserted to DB | user_id:%v", ul.GetUserID())
	return ul
}

func (ul *UserLog) SetUserId(userId int64) *UserLog {
	ul.UserLog.UserID = Processors.Int64(userId)
	return ul
}

func (ul *UserLog) SetFoodId(foodId string) *UserLog {
	ul.UserLog.FoodID = Processors.String(foodId)
	return ul
}

func (ul *UserLog) SetOrderTime(orderTime int64) *UserLog {
	ul.UserLog.OrderTime = Processors.Int64(orderTime)
	return ul
}

func (ul *UserLog) SetStatus(status int64) *UserLog {
	ul.UserLog.Status = Processors.Int64(status)
	return ul
}

func (ul *UserLog) SetErrorMsg(errorMsg string) *UserLog {
	ul.UserLog.ErrorMsg = Processors.String(errorMsg)
	return ul
}

func (ul *UserLog) TearDown() error {
	if err := Processors.DB.Exec("DELETE FROM user_log_tab WHERE user_id = ?", ul.GetUserID()).Error; err != nil {
		log.Printf("Failed to delete from DB | user_id:%v", ul.GetUserID())
		return err
	}
	log.Printf("Sulcessfully deleted from DB | user_id:%v", ul.GetUserID())
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
