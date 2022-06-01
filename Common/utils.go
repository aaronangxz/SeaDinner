package Common

import (
	"log"
	"os"
)

func GetTGToken() string {
	if os.Getenv("TEST_DEPLOY") == "TRUE" || Config.Adhoc {
		log.Println("Running Test Telegram Bot Instance")
		return os.Getenv("TELEGRAM_TEST_APITOKEN")
	}
	return os.Getenv("TELEGRAM_APITOKEN")
}

func IsInGrayScale(userId int64) bool {
	return userId%100 >= Config.GrayScale.Percentage
}
