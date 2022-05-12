package Processors

import (
	"fmt"
	"log"
)

func PrepOrder() ([]UserChoiceWithKey, bool) {
	var (
		record []UserChoiceWithKey
	)
	//check whole db
	if err := DB.Raw("SELECT c.*, k.key FROM user_choice c, user_key k WHERE c.user_id = k.user_id").Scan(&record).Error; err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	log.Println("Fetched user_records:", len(record))
	return record, true
}
