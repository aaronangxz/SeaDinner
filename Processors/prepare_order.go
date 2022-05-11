package Processors

import (
	"fmt"
	"log"
)

func PrepOrder() []UserRecords {
	var (
		record []UserRecords
	)
	//check whole db
	if err := DB.Raw("SELECT * FROM user_records").Scan(&record).Error; err != nil {
		fmt.Println(err.Error())
		return nil
	}
	log.Println("Fetched user_records:", len(record))
	return record
}
