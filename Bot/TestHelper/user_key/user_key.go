package user_key

import (
	"log"
	"os"
	"time"

	"github.com/aaronangxz/SeaDinner/Bot/TestHelper"
	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/Processors"
	"google.golang.org/protobuf/proto"
)

var (
	defaultKey   = TestHelper.RandomString(10)
	defaultCtime = time.Now().Unix()
	defaultMtime = time.Now().Unix()
)

type UserKey struct {
	*TestHelper.UserKey
}

func New() *UserKey {
	TestHelper.InitTest()
	return &UserKey{
		UserKey: &TestHelper.UserKey{
			UserID:  new(int64),
			UserKey: new(string),
			Ctime:   new(int64),
			Mtime:   new(int64),
		},
	}
}

func (uk *UserKey) FillDefaults() *UserKey {
	if uk.UserKey.GetUserID() == 0 {
		uk.SetUserId(TestHelper.RandomInt(99999))
	}

	if uk.UserKey.GetUserKey() == "" {
		uk.SetKey(Processors.EncryptKey(defaultKey, os.Getenv("AES_KEY")))
	}

	if uk.UserKey.GetCtime() == 0 {
		uk.SetCtime(defaultCtime)
	}

	if uk.UserKey.GetMtime() == 0 {
		uk.SetMtime(defaultMtime)
	}
	return uk
}

func (uk *UserKey) Build() *UserKey {
	uk.FillDefaults()
	if err := Processors.DB.Table(Common.DB_USER_KEY_TAB).Create(&uk).Error; err != nil {
		log.Printf("Failed to insert to DB | user_id:%v | %v", uk.GetUserID(), err.Error())
		return nil
	}
	log.Printf("Successfully inserted to DB | user_id:%v", uk.GetUserID())
	return uk
}

func (uk *UserKey) SetUserId(userId int64) *UserKey {
	uk.UserKey.UserID = proto.Int64(userId)
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

func (uk *UserKey) TearDown() error {
	if err := Processors.DB.Exec("DELETE FROM user_key_tab WHERE user_id = ?", uk.GetUserID()).Error; err != nil {
		log.Printf("Failed to delete from DB | user_id:%v", uk.GetUserID())
		return err
	}
	log.Printf("Successfully deleted from DB | user_id:%v", uk.GetUserID())
	return nil
}

func DeleteUserKey(userId int64) error {
	if err := Processors.DB.Exec("DELETE FROM user_key_tab WHERE user_id = ?", userId).Error; err != nil {
		log.Printf("Failed to delete from DB | user_id:%v", userId)
		return err
	}

	log.Printf("Successfully deleted from DB | user_id:%v", userId)
	return nil
}
