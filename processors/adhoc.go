package processors

import (
	"context"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"

	"github.com/aaronangxz/SeaDinner/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//SendAdHocNotification Sends a one-off message to a specific user
func SendAdHocNotification(ctx context.Context, id int64, text string) {

	bot, err := tgbotapi.NewBotAPI(common.GetTGToken(ctx))
	if err != nil {
		log.Error(ctx, err.Error())
	}
	bot.Debug = true
	log.Info(ctx, "Authorized on account %s", bot.Self.UserName)

	msg := tgbotapi.NewMessage(id, text)
	msg.Text = text

	if _, err := bot.Send(msg); err != nil {
		log.Error(ctx, err.Error())
	}
}

func SendAdHocNotificationToAll(ctx context.Context, text string) {
	var (
		record []*sea_dinner.UserKey
	)
	if err := DbInstance().Raw("SELECT * FROM user_key_tab WHERE is_mute <> ?", sea_dinner.MuteStatus_MUTE_STATUS_YES).Scan(&record).Error; err != nil {
		log.Error(ctx, err.Error())
	}

	bot, err := tgbotapi.NewBotAPI(common.GetTGToken(ctx))
	if err != nil {
		log.Error(ctx, err.Error())
	}
	bot.Debug = true
	log.Info(ctx, "Authorized on account %s", bot.Self.UserName)

	for _, r := range record {
		msg := tgbotapi.NewMessage(r.GetUserId(), text)
		msg.Text = text

		if _, err := bot.Send(msg); err != nil {
			log.Error(ctx, err.Error())
		}
	}

}
