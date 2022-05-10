package Processors

import "fmt"

func PrepOrder() []UserRecords {
	var (
		record []UserRecords
	)
	//check whole db
	if err := DB.Raw("SELECT * FROM user_records").Scan(&record).Error; err != nil {
		fmt.Println(err.Error())
		return nil
	}

	for _, r := range record {
		fmt.Printf("%v:%v:%v\n", r.UserID, r.Key, r.Choice)
	}

	return record
}
