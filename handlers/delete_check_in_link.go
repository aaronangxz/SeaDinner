package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

//DeleteCheckInLink Deletes the supposingly expired check-in link
func DeleteCheckInLink(ctx context.Context) {
	txn := processors.App.StartTransaction("delete_check_in_link")
	defer txn.End()

	//Retrieve the whole set
	s := processors.CacheInstance().SMembers(common.CHECK_IN_LINK_SET)
	if s == nil {
		log.Error(ctx, "DeleteCheckInLink | Set is empty.")
		return
	}

	for _, pair := range s.Val() {
		//split <user_id>:<message_id> by ':'
		split := strings.Split(pair, ":")
		userID, _ := strconv.Atoi(split[0])
		msgID, _ := strconv.Atoi(split[1])

		bot, err := tgbotapi.NewBotAPI(common.GetTGToken(ctx))
		if err != nil {
			log.Error(ctx, err.Error())
		}
		bot.Debug = true
		log.Info(ctx, "Authorized on account %s", bot.Self.UserName)
		c := tgbotapi.NewDeleteMessage(int64(userID), msgID)
		if _, err := bot.Send(c); err != nil {
			log.Error(ctx, err.Error())
		}
	}
	log.Info(ctx, "DeleteCheckInLink | Successfully deleted check in links.")

	//Clear set
	if err := processors.CacheInstance().Del(common.CHECK_IN_LINK_SET).Err(); err != nil {
		log.Error(ctx, "DeleteCheckInLink | Error while erasing from redis: %v", err.Error())
	} else {
		log.Info(ctx, "DeleteCheckInLink | Successful | Deleted checkin_link set")
	}
}
