package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"strings"
)

//BatchGetUsersChoiceWithKey Retrieves the user's choice and key. Only return those that has valid choices in the current week.
func BatchGetUsersChoiceWithKey(ctx context.Context) ([]*sea_dinner.UserChoiceWithKey, error) {
	var (
		record []*sea_dinner.UserChoiceWithKey
	)
	txn := processors.App.StartTransaction("batch_get_users_choice_with_key")
	defer txn.End()

	m := MakeMenuNameMap(ctx)
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
	query := fmt.Sprintf("SELECT c.*, k.user_key FROM user_choice_tab c, user_key_tab k WHERE user_choice IN %v AND c.user_id = k.user_id", inQuery)
	log.Info(ctx, query)
	//check whole db
	if err := processors.DB.Raw(query).Scan(&record).Error; err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	log.Info(ctx, "BatchGetUsersChoiceWithKey | Success | size: %v", len(record))
	return record, nil
}
