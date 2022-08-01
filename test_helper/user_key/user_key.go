package user_key

import (
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/test_helper"
	"log"
	"os"
	"time"

	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"google.golang.org/protobuf/proto"
)

var (
	defaultKey    = test_helper.RandomString(10)
	defaultCtime  = time.Now().Unix()
	defaultMtime  = time.Now().Unix()
	defaultIsMute = int64(sea_dinner.MuteStatus_MUTE_STATUS_NO)
	defaultStatus = int64(sea_dinner.UserStatus_USER_STATUS_ACTIVE)
)

type UserKey struct {
	*sea_dinner.UserKey
}

func New() *UserKey {
	test_helper.InitTest()
	return &UserKey{
		UserKey: &sea_dinner.UserKey{
			UserId:  new(int64),
			UserKey: new(string),
			Ctime:   new(int64),
			Mtime:   new(int64),
			IsMute:  new(int64),
			Status:  new(int64),
		},
	}
}

func (uk *UserKey) FillDefaults() *UserKey {
	if uk.UserKey.GetUserId() == 0 {
		uk.SetUserId(test_helper.RandomInt(99999))
	}

	if uk.UserKey.GetUserKey() == "" {
		uk.SetKey(processors.EncryptKey(defaultKey, os.Getenv("AES_KEY")))
	}

	if uk.UserKey.GetCtime() == 0 {
		uk.SetCtime(defaultCtime)
	}

	if uk.UserKey.GetMtime() == 0 {
		uk.SetMtime(defaultMtime)
	}

	if uk.UserKey.GetIsMute() == 0 {
		uk.SetIsMute(defaultIsMute)
	}

	if uk.UserKey.GetStatus() == 0 {
		uk.SetStatus(defaultStatus)
	}
	return uk
}

func (uk *UserKey) Build() *UserKey {
	uk.FillDefaults()
	if err := processors.DbInstance().Table(common.DB_USER_KEY_TAB).Create(&uk).Error; err != nil {
		log.Printf("Failed to insert to DB | user_id:%v | %v", uk.GetUserId(), err.Error())
		return nil
	}
	log.Printf("Successfully inserted to DB | user_id:%v", uk.GetUserId())
	return uk
}

func (uk *UserKey) SetUserId(userId int64) *UserKey {
	uk.UserKey.UserId = proto.Int64(userId)
	return uk
}

func (uk *UserKey) SetKey(userKey string) *UserKey {
	uk.UserKey.UserKey = proto.String(userKey)
	return uk
}

func (uk *UserKey) SetCtime(ctime int64) *UserKey {
	uk.UserKey.Ctime = proto.Int64(ctime)
	return uk
}

func (uk *UserKey) SetMtime(mtime int64) *UserKey {
	uk.UserKey.Mtime = proto.Int64(mtime)
	return uk
}

func (uk *UserKey) SetIsMute(mute int64) *UserKey {
	uk.UserKey.IsMute = proto.Int64(mute)
	return uk
}

func (uk *UserKey) SetStatus(status int64) *UserKey {
	uk.UserKey.Status = proto.Int64(status)
	return uk
}

func (uk *UserKey) TearDown() error {
	if err := processors.DbInstance().Exec("DELETE FROM user_key_tab WHERE user_id = ?", uk.GetUserId()).Error; err != nil {
		log.Printf("Failed to delete from DB | user_id:%v", uk.GetUserId())
		return err
	}
	log.Printf("Successfully deleted from DB | user_id:%v", uk.GetUserId())
	return nil
}

func DeleteUserKey(userId int64) error {
	if err := processors.DbInstance().Exec("DELETE FROM user_key_tab WHERE user_id = ?", userId).Error; err != nil {
		log.Printf("Failed to delete from DB | user_id:%v", userId)
		return err
	}

	log.Printf("Successfully deleted from DB | user_id:%v", userId)
	return nil
}

func CheckUserKey(userId int64) *sea_dinner.UserKey {
	var (
		row *sea_dinner.UserKey
	)
	if err := processors.DbInstance().Raw("SELECT * FROM user_key_tab WHERE user_id = ?", userId).Scan(&row).Error; err != nil {
		log.Printf("Failed to read from DB | user_id:%v", userId)
		return nil
	}

	log.Printf("Successfully read from DB | user_id:%v", userId)
	return row
}
