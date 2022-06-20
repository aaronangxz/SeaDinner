package main

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	handlers "github.com/aaronangxz/SeaDinner/handlers"
	"os"
	"strconv"
	"time"

	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot              *tgbotapi.BotAPI
	err              error
	startListenKey   = false
	startListenChope = false
)

func resumeLog() {
	bot.Debug = true
}

func skipLog() {
	bot.Debug = false
}

func main() {
	log.InitializeLogger()
	processors.LoadEnv()
	processors.Init()
	processors.InitClient()

	bot, err = tgbotapi.NewBotAPI(common.GetTGToken(context.TODO()))
	if err != nil {
		log.Error(context.TODO(), err.Error())
	}

	bot.Debug = true
	log.Info(context.TODO(), "Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 3600
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		ctx := log.NewCtx()
		resumeLog()
		if update.CallbackQuery != nil {
			var muteType bool
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
			msg.ParseMode = "MARKDOWN"
			msg.Text, muteType = handlers.CallbackQueryHandler(ctx, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery)

			//Handle MUTE related callback separately
			if update.CallbackQuery.Data == "MUTE" || update.CallbackQuery.Data == "UNMUTE" {
				//Retrieve the previous chat ID after user calls /mute
				cacheKey := fmt.Sprint(common.USER_MUTE_MSG_ID_PREFIX, update.CallbackQuery.Message.Chat.ID)
				val, redisErr := processors.RedisClient.Get(cacheKey).Result()
				if redisErr != nil {
					if redisErr == redis.Nil {
						log.Warn(ctx, "Callback Mute | No result of %v in Redis", cacheKey)
					} else {
						log.Error(ctx, "Callback Mute | Error while reading from redis: %v", redisErr.Error())
					}
					//Return expired message if not found
					msg.Text = "Oops this selection had expired. Start over at /mute!"
					skipLog()
					if _, err := bot.Send(msg); err != nil {
						log.Error(ctx, err.Error())
					}
					resumeLog()
					continue
				} else {
					var out []tgbotapi.InlineKeyboardMarkup
					if muteType {
						var rows []tgbotapi.InlineKeyboardButton
						unmuteBotton := tgbotapi.NewInlineKeyboardButtonData("Turn ON 🔔", "UNMUTE")
						rows = append(rows, unmuteBotton)
						out = append(out, tgbotapi.NewInlineKeyboardMarkup(rows))
					} else {
						var rows []tgbotapi.InlineKeyboardButton
						muteButton := tgbotapi.NewInlineKeyboardButtonData("Turn OFF 🔕", "MUTE")
						rows = append(rows, muteButton)
						out = append(out, tgbotapi.NewInlineKeyboardMarkup(rows))
					}
					//Edit the previous message and buttons
					intVar, _ := strconv.Atoi(val)
					c := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, intVar, msg.Text)
					c.ParseMode = "MARKDOWN"
					c.ReplyMarkup = &out[0]
					skipLog()
					if _, err := bot.Send(c); err != nil {
						log.Error(ctx, err.Error())
					}
					resumeLog()
					continue
				}
			} else if update.CallbackQuery.Data == "ATTEMPTCANCEL" {
				var mk tgbotapi.InlineKeyboardMarkup
				var out [][]tgbotapi.InlineKeyboardButton
				var rows []tgbotapi.InlineKeyboardButton

				cancelButton := tgbotapi.NewInlineKeyboardButtonData("⚠️ DO IT ⚠️", "CANCEL")
				rows = append(rows, cancelButton)
				skipButton := tgbotapi.NewInlineKeyboardButtonData("OK NVM", "SKIP")
				rows = append(rows, skipButton)

				c := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
				c.Text = "Are you sure you want to cancel? 😦 You have to re-select it in SeaTalk."
				c.ParseMode = "MARKDOWN"
				out = append(out, rows)
				mk.InlineKeyboard = out
				c.ReplyMarkup = mk
				skipLog()
				if _, err := bot.Send(c); err != nil {
					log.Error(ctx, err.Error())
				}
				resumeLog()
				continue
			}
			skipLog()
			if _, err := bot.Send(msg); err != nil {
				log.Error(ctx, err.Error())
			}
			resumeLog()
			continue
		}

		//Stop responding from 12.29pm to 12.31pm or until dinner order has started (For occasional weird order timings)
		if time.Now().Unix() >= processors.GetLunchTime().Unix()-60 &&
			(time.Now().Unix() <= processors.GetLunchTime().Unix()+60 && !processors.IsPollStart()) {
			if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Omw to order, wait for my good news! 🏃")); err != nil {
				log.Error(ctx, err.Error())
			}
			continue
		}

		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			if startListenKey {
				//Capture key
				msg, _ := handlers.UpdateKey(ctx, update.Message.Chat.ID, update.Message.Text)
				skipLog()
				if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg)); err != nil {
					log.Error(ctx, err.Error())
				}
				resumeLog()
				startListenKey = false
				continue
			} else if startListenChope {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				ok := false
				msg.Text, ok = handlers.UpdateChope(ctx, update.Message.Chat.ID, update.Message.Text)
				if !ok {
					if _, err := bot.Send(msg); err != nil {
						log.Error(ctx, err.Error())
					}
					continue
				}
				//Capture chope
				msg.ParseMode = "MARKDOWN"
				skipLog()
				if _, err := bot.Send(msg); err != nil {
					log.Error(ctx, err.Error())
				}
				resumeLog()
				startListenChope = false
				continue
			} else {
				continue
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "start":
			s, ok := handlers.CheckKey(ctx, update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				msg.Text = "Hello! " + update.Message.Chat.UserName
			}
		case "menu":
			s, ok := handlers.CheckKey(ctx, update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				skipLog()
				txt, mp := processors.OutputMenuWithButton(ctx, handlers.GetKey(ctx, update.Message.Chat.ID))
				for i, r := range txt {
					msg.Text = r
					if len(mp) > 0 {
						msg.ReplyMarkup = mp[i]
					}
					if _, err := bot.Send(msg); err != nil {
						log.Error(ctx, err.Error())
					}
				}
				resumeLog()
				continue
			}
		case "help":
			skipLog()
			msg.Text = handlers.MakeHelpResponse()
			msg.ParseMode = "MARKDOWN"
		case "key":
			skipLog()
			msg.Text, _ = handlers.CheckKey(ctx, update.Message.Chat.ID)
		case "newkey":
			skipLog()
			msg.Text = "What's your key? \nGo to https://dinner.sea.com/accounts/token, copy the Key under Generate Auth Token and paste it here:"
			startListenKey = true
		case "status":
			s, ok := handlers.CheckKey(ctx, update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				skipLog()
				msg.Text = handlers.ListWeeklyResultByUserID(ctx, update.Message.Chat.ID)
				msg.ParseMode = "HTML"
			}
		case "chope":
			msg.Text = "This command is deprecated. Choose from /menu instead!😋"
		case "choice":
			s, ok := handlers.CheckKey(ctx, update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				skipLog()
				msg.Text, _ = handlers.CheckChope(ctx, update.Message.Chat.ID)
			}
		case "reminder":
			//Backdoor for test env
			if os.Getenv("TEST_DEPLOY") == "TRUE" || common.Config.Adhoc {
				handlers.SendReminder(ctx)
			}
		case "mute":
			s, ok := handlers.CheckKey(ctx, update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				txt, kb := handlers.CheckMute(ctx, update.Message.Chat.ID)
				msg.Text = txt
				if kb != nil {
					msg.ReplyMarkup = kb[0]
				}
				msg.ParseMode = "MARKDOWN"
				skipLog()
				if msgTrace, err := bot.Send(msg); err != nil {
					log.Error(ctx, err.Error())
				} else {
					//save msg id into cache for msg update
					cacheKey := fmt.Sprint(common.USER_MUTE_MSG_ID_PREFIX, update.Message.Chat.ID)
					if err := processors.RedisClient.Set(cacheKey, msgTrace.MessageID, 1800*time.Second).Err(); err != nil {
						log.Error(ctx, "CheckMute | Error while writing to redis: %v", err.Error())
					} else {
						log.Info(ctx, "CheckMute | Successful | Written %v to redis", cacheKey)
					}
					resumeLog()
					continue
				}
			}
		case "checkin":
			//Backdoor for test env
			if os.Getenv("TEST_DEPLOY") == "TRUE" || common.Config.Adhoc {
				handlers.SendCheckInLink(ctx)
			}
		case "delete":
			//Backdoor for test env
			if os.Getenv("TEST_DEPLOY") == "TRUE" || common.Config.Adhoc {
				handlers.DeleteCheckInLink(ctx)
			}
		default:
			msg.Text = "I don't understand this command :("
		}
		if msg.Text != "" {
			if _, err := bot.Send(msg); err != nil {
				log.Error(ctx, err.Error())
			}
		}
		resumeLog()
	}
}
