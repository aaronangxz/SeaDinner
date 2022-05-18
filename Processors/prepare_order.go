package Processors

import (
	"fmt"
	"log"
	"strings"
)

func PrepOrder() ([]UserChoiceWithKeyAndStatus, bool) {
	var (
		record []UserChoiceWithKeyAndStatus
	)

	m := MakeMenuMap()
	inQuery := "("
	for e := range m {
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

	log.Println("PrepOrder | Fetched user_records:", len(record))
	return record, true
}
