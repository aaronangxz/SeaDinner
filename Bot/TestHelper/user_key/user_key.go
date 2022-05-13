package user_key

import (
	"log"
	"time"

	"github.com/aaronangxz/SeaDinner/Bot/TestHelper"
	"github.com/aaronangxz/SeaDinner/Processors"
)

var (
	defaultUserID = TestHelper.RandomInt(99999)
	defaultKey    = TestHelper.RandomString(10)
	defaultCtime  = time.Now().Unix()
	defaultMtime  = time.Now().Unix()
)

type UserKey struct {
	*TestHelper.UserKey
}

func New() *UserKey {
	return &UserKey{
		UserKey: &TestHelper.UserKey{
			UserID: new(int64),
			Key:    new(string),
			Ctime:  new(int64),
			Mtime:  new(int64),
		},
	}
}

func (uk *UserKey) FillDefaults() *UserKey {
	if uk.UserKey.GetUserID() == 0 {
		uk.SetUserId(defaultUserID)
	}

	if uk.UserKey.GetKey() == "" {
		uk.SetKey(defaultKey)
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
	TestHelper.InitTest()
	if err := Processors.DB.Table("user_key").Create(&uk).Error; err != nil {
		log.Printf("Failed to insert to DB | user_id:%v | %v", uk.GetUserID(), err.Error())
		return nil
	}
	log.Printf("Successfully inserted to DB | user_id:%v", uk.GetUserID())
	return uk
}

func (uk *UserKey) SetUserId(userId int64) *UserKey {
	uk.UserKey.UserID = Processors.Int64(userId)
	return uk
}

func (uk *UserKey) SetKey(key string) *UserKey {
	uk.UserKey.Key = Processors.String(key)
	return uk
}

func (uk *UserKey) SetCtime(ctime int64) *UserKey {
	uk.UserKey.Ctime = Processors.Int64(ctime)
	return uk
}

func (uk *UserKey) SetMtime(mtime int64) *UserKey {
	uk.UserKey.Mtime = Processors.Int64(mtime)
	return uk
}

func (uk *UserKey) TearDown() error {
	if err := Processors.DB.Exec("DELETE FROM user_key WHERE user_id = ?", uk.GetUserID()).Error; err != nil {
		log.Printf("Failed to delete from DB | user_id:%v", uk.GetUserID())
		return err
	}
	log.Printf("Successfully deleted from DB | user_id:%v", uk.GetUserID())
	return nil
}
