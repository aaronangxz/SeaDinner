package processors

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
)

//GetOrderByUserID Retrieves user_choice of a single user
func GetOrderByUserID(ctx context.Context, userID int64) (string, bool) {
	var (
		record *sea_dinner.UserChoice
	)

	if err := DB.Raw("SELECT * FROM user_choice_tab WHERE user_id = ?", userID).Scan(&record).Error; err != nil {
		log.Error(ctx, "GetOrderByUserId | failed to retrieve record: %v", err.Error())
		return "I can't find your order 😥 Try to cancel from SeaTalk instead!", false
	}

	if record == nil {
		return "I can't find your order 😥 Try to cancel from SeaTalk instead!", false
	}
	log.Info(ctx, "GetOrderByUserId | Success")
	return record.GetUserChoice(), true
}
