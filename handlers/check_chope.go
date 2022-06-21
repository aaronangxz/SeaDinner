package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"time"
)

//CheckChope Retrieves the current food choice made by user.
func CheckChope(ctx context.Context, id int64) (string, bool) {
	var (
		existingRecord sea_dinner.UserChoice
		dayText        = "today"
	)

	txn := processors.App.StartTransaction("check_chope")
	defer txn.End()

	if id <= 0 {
		log.Error(ctx, "Id must be > 1.")
		return "", false
	}
	tz, _ := time.LoadLocation(processors.TimeZone)
	if !processors.IsWeekDay() || !processors.IsNotEOW(time.Now().In(tz)) {
		return "We are done for this week! You can tell me your order again next week 😀", false
	}

	if err := processors.DbInstance().Raw("SELECT * FROM user_choice_tab WHERE user_id = ?", id).Scan(&existingRecord).Error; err != nil {
		return "I have yet to receive your order 🥲 You can choose from /menu", false
	}
	if existingRecord.UserChoice == nil {
		return "I have yet to receive your order 🥲 You can choose from /menu", false
	} else if existingRecord.GetUserChoice() == "-1" {
		//Dynamic text based on time - shows tomorrow if current time is past lunch
		tz, _ := time.LoadLocation(processors.TimeZone)
		if time.Now().In(tz).Unix() > processors.GetLunchTime().Unix() {
			if processors.IsNotEOW(time.Now().In(tz)) {
				dayText = "tomorrow"
			}
		}
		return fmt.Sprintf("Not placing dinner order for you %v 🙅 Changed your mind? You can choose from /menu", dayText), false
	}
	menu := MakeMenuNameMap(ctx)

	_, ok := menu[existingRecord.GetUserChoice()]

	if !ok {
		return fmt.Sprintf("Your choice %v is not available this week, so I will not order anything 🥲 Choose a new dish from /menu", existingRecord.GetUserChoice()), true
	}

	if existingRecord.GetUserChoice() == "RAND" {
		return "I'm tasked to snatch a random dish for you 😀 Changed your mind? You can choose from /menu", true
	}
	return fmt.Sprintf("I'm tasked to snatch %v for you 😀 Changed your mind? You can choose from /menu", menu[existingRecord.GetUserChoice()]), true
}
