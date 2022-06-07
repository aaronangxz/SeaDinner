package Processors

import (
	"fmt"
	"log"
	"strings"

	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"google.golang.org/protobuf/proto"
)

//PrepOrder Retrieves all the user's key and choice, where the user's choice is in the current menu / RAND.
//If choice is RAND, generates a random food id.
//Returns UserChoiceWithKey
func PrepOrder() ([]*sea_dinner.UserChoiceWithKey, bool) {
	var (
		record []*sea_dinner.UserChoiceWithKey
	)

	m := MakeMenuMap()
	inQuery := "("
	for e := range m {
		// Skip menu id: -1
		if e == "-1" {
			continue
		}
		if e == "RAND" {
			inQuery += "'RAND', "
			continue
		}
		inQuery += e + ", "
	}
	inQuery += ")"
	inQuery = strings.ReplaceAll(inQuery, ", )", ")")
	query := fmt.Sprintf("SELECT c.*, k.user_key FROM user_choice_tab c, user_key_tab k WHERE user_choice IN %v AND c.user_id = k.user_id", inQuery)
	log.Println(query)

	//check whole db
	if err := DB.Raw(query).Scan(&record).Error; err != nil {
		fmt.Println(err.Error())
		return nil, false
	}

	for _, r := range record {
		if r.GetUserChoice() == "RAND" {
			r.UserChoice = proto.String(RandomFood(m))
			log.Printf("PrepOrder | id:%v | random choice:%v", r.GetUserId(), r.GetUserChoice())
		}
	}

	log.Println("PrepOrder | Fetched user_records:", len(record))
	return record, true
}

func GetOrderByUserId(user_id int64) (string, bool) {
	var (
		record *sea_dinner.UserChoice
	)

	if err := DB.Raw("SELECT * FROM user_choice_tab WHERE user_id = ?", user_id).Scan(&record).Error; err != nil {
		fmt.Printf("GetOrderByUserId | failed to retrieve record: %v", err.Error())
		return "I can't find your order 😥 Try to cancel from SeaTalk instead!", false
	}

	if record == nil {
		return "I can't find your order 😥 Try to cancel from SeaTalk instead!", false
	}

	return record.GetUserChoice(), true
}
