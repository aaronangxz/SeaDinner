package Processors

import (
	"fmt"
	"time"
)

func GetLunchTime() time.Time {
	sg, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		fmt.Println(err)
	}
	now := time.Now().In(sg)

	year, month, day := now.Date()
	return time.Date(year, month, day, 12, 30, 0, 0, now.Location())
}
