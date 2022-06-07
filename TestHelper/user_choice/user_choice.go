package user_choice

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
	defaultCtime = time.Now().Unix()
	defaultMtime = time.Now().Unix()
)

type UserChoice struct {
	*sea_dinner.UserChoice
}

func New() *UserChoice {
	TestHelper.InitTest()
	return &UserChoice{
		UserChoice: &sea_dinner.UserChoice{
			UserId:     new(int64),
			UserChoice: new(string),
			Ctime:      new(int64),
			Mtime:      new(int64),
		},
	}
}

func (uk *UserChoice) FillDefaults() *UserChoice {
	if uk.UserChoice.GetUserId() == 0 {
		uk.SetUserId(TestHelper.RandomInt(99999))
	}

	if uk.UserChoice.GetUserChoice() == "" {
		uk.SetUserChoice(fmt.Sprint(TestHelper.RandomInt(9999)))
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
		log.Printf("Failed to insert to DB | user_id:%v | %v", uc.GetUserId(), err.Error())
		return nil
	}
	log.Printf("Successfully inserted to DB | user_id:%v", uc.GetUserId())
	return uc
}

func (uc *UserChoice) SetUserId(userId int64) *UserChoice {
	uc.UserChoice.UserId = proto.Int64(userId)
	return uc
}

func (uc *UserChoice) SetUserChoice(userChoice string) *UserChoice {
	uc.UserChoice.UserChoice = proto.String(userChoice)
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
	if err := Processors.DB.Exec("DELETE FROM user_choice_tab WHERE user_id = ?", uc.GetUserId()).Error; err != nil {
		log.Printf("Failed to delete from DB | user_id:%v", uc.GetUserId())
		return err
	}
	log.Printf("Successfully deleted from DB | user_id:%v", uc.GetUserId())
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
