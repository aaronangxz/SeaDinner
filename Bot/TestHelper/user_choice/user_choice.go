package user_choice

import (
	"log"
	"time"

	"github.com/aaronangxz/SeaDinner/Bot/TestHelper"
	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/Processors"
	"google.golang.org/protobuf/proto"
)

var (
	defaultCtime = time.Now().Unix()
	defaultMtime = time.Now().Unix()
)

type UserChoice struct {
	*TestHelper.UserChoice
}

func New() *UserChoice {
	TestHelper.InitTest()
	return &UserChoice{
		UserChoice: &TestHelper.UserChoice{
			UserID:     new(int64),
			UserChoice: new(int64),
			Ctime:      new(int64),
			Mtime:      new(int64),
		},
	}
}

func (uk *UserChoice) FillDefaults() *UserChoice {
	if uk.UserChoice.GetUserID() == 0 {
		uk.SetUserId(TestHelper.RandomInt(99999))
	}

	if uk.UserChoice.GetUserChoice() == 0 {
		uk.SetUserChoice(TestHelper.RandomInt(9999))
	}

	if uk.UserChoice.GetCtime() == 0 {
		uk.SetCtime(defaultCtime)
	}

	if uk.UserChoice.GetMtime() == 0 {
		uk.SetMtime(defaultMtime)
	}
	return uk
}

func (uc *UserChoice) Build() *UserChoice {
	uc.FillDefaults()
	if err := Processors.DB.Table(Common.DB_USER_CHOICE_TAB).Create(&uc).Error; err != nil {
		log.Printf("Failed to insert to DB | user_id:%v | %v", uc.GetUserID(), err.Error())
		return nil
	}
	log.Printf("Successfully inserted to DB | user_id:%v", uc.GetUserID())
	return uc
}

func (uc *UserChoice) SetUserId(userId int64) *UserChoice {
	uc.UserChoice.UserID = proto.Int64(userId)
	return uc
}

func (uc *UserChoice) SetUserChoice(userChoice int64) *UserChoice {
	uc.UserChoice.UserChoice = proto.Int64(userChoice)
	return uc
}

func (uc *UserChoice) SetCtime(ctime int64) *UserChoice {
	uc.UserChoice.Ctime = proto.Int64(ctime)
	return uc
}

func (uc *UserChoice) SetMtime(mtime int64) *UserChoice {
	uc.UserChoice.Mtime = proto.Int64(mtime)
	return uc
}

func (uc *UserChoice) TearDown() error {
	if err := Processors.DB.Exec("DELETE FROM user_choice_tab WHERE user_id = ?", uc.GetUserID()).Error; err != nil {
		log.Printf("Failed to delete from DB | user_id:%v", uc.GetUserID())
		return err
	}
	log.Printf("Successfully deleted from DB | user_id:%v", uc.GetUserID())
	return nil
}

func DeleteUserKey(userId int64) error {
	if err := Processors.DB.Exec("DELETE FROM user_choice_tab WHERE user_id = ?", userId).Error; err != nil {
		log.Printf("Failed to delete from DB | user_id:%v", userId)
		return err
	}
	log.Printf("Successfully deleted from DB | user_id:%v", userId)
	return nil
}
