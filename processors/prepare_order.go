package processors

import (
	"context"
	"fmt"
	"strings"

	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"google.golang.org/protobuf/proto"
)

//PrepOrder Retrieves all the user's key and choice, where the user's choice is in the current menu / RAND.
//If choice is RAND, generates a random food id.
//Returns UserChoiceWithKey
func PrepOrder(ctx context.Context) ([]*sea_dinner.UserChoiceWithKey, bool) {
	var (
		record []*sea_dinner.UserChoiceWithKey
	)

	m := MakeMenuMap(ctx)
	inQuery := "("
	for e := range m {
		// Skip menu id: -1
		if e == "-1" {
			continue
		}
		if e == "RAND" {
			inQuery += "'RAND', "
			continue
		}
		inQuery += e + ", "
	}
	inQuery += ")"
	inQuery = strings.ReplaceAll(inQuery, ", )", ")")
	query := fmt.Sprintf("SELECT c.*, k.user_key FROM user_choice_tab c, user_key_tab k WHERE user_choice IN %v AND c.user_id = k.user_id AND k,status = %v", inQuery, sea_dinner.UserStatus_USER_STATUS_ACTIVE)
	log.Info(ctx, query)

	//check whole db
	if err := DbInstance().Raw(query).Scan(&record).Error; err != nil {
		log.Error(ctx, err.Error())
		return nil, false
	}

	for _, r := range record {
		if r.GetUserChoice() == "RAND" {
			r.UserChoice = proto.String(RandomFood(ctx, m))
			log.Info(ctx, "PrepOrder | id:%v | random choice:%v", r.GetUserId(), r.GetUserChoice())
		}
	}
	log.Info(ctx, "PrepOrder | Fetched user_records: %v", len(record))
	return record, true
}
