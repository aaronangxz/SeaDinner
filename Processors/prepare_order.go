package Processors

import (
	"fmt"
	"log"
)

func PrepOrder() ([]UserChoiceWithKeyAndStatus, bool) {
	var (
		record []UserChoiceWithKeyAndStatus
	)
	//check whole db
	if err := DB.Raw("SELECT c.*, k.user_key FROM user_choice_tab c, user_key_tab k WHERE c.user_id = k.user_id").Scan(&record).Error; err != nil {
		fmt.Println(err.Error())
		return nil, false
	}

	for _, r := range record {
		r.IsSuccess = Bool(false)
	}

	log.Println("Fetched user_records:", len(record))
	return record, true
}
