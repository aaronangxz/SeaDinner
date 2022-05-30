package Common

import "os"

func GetTGToken() string {
	if os.Getenv("TEST_DEPLOY") == "TRUE" {
		return os.Getenv("TELEGRAM_TEST_APITOKEN")
	}
	return os.Getenv("TELEGRAM_APITOKEN")
}
