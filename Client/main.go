package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aaronangxz/SeaDinner/Bot"
	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/Log"
	"github.com/aaronangxz/SeaDinner/Processors"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	startListenKey   = false
	startListenChope = false
)

func main() {
	Processors.LoadEnv()
	Processors.Init()
	Processors.InitClient()
	Log.InitializeLogger()

	bot, err := tgbotapi.NewBotAPI(Common.GetTGToken(context.TODO()))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	Log.Info(context.TODO(), "Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 3600
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		ctx := Log.NewCtx()
		if update.CallbackQuery != nil {
			var muteType bool
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
			msg.ParseMode = "MARKDOWN"
			msg.Text, muteType = Bot.CallbackQueryHandler(ctx, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery)

			//Handle MUTE related callback separately
			if update.CallbackQuery.Data == "MUTE" || update.CallbackQuery.Data == "UNMUTE" {
				//Retrieve the previous chat ID after user calls /mute
				cacheKey := fmt.Sprint(Common.USER_MUTE_MSG_ID_PREFIX, update.CallbackQuery.Message.Chat.ID)
				val, redisErr := Processors.RedisClient.Get(cacheKey).Result()
				if redisErr != nil {
					if redisErr == redis.Nil {
						Log.Warn(ctx, "Callback Mute | No result of %v in Redis", cacheKey)
						//log.Printf("Callback Mute | No result of %v in Redis", cacheKey)
					} else {
						Log.Error(ctx, "Callback Mute | Error while reading from redis: %v", redisErr.Error())
						//log.Printf("Callback Mute | Error while reading from redis: %v", redisErr.Error())
					}
					//Return expired message if not found
					msg.Text = "Oops this selection had expired. Start over at /mute!"
					if _, err := bot.Send(msg); err != nil {
						Log.Error(ctx, err.Error())
						//log.Println(err)
					}
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
						muteBotton := tgbotapi.NewInlineKeyboardButtonData("Turn OFF 🔕", "MUTE")
						rows = append(rows, muteBotton)
						out = append(out, tgbotapi.NewInlineKeyboardMarkup(rows))
					}
					//Edit the previous message and buttons
					intVar, _ := strconv.Atoi(val)
					c := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, intVar, msg.Text)
					c.ParseMode = "MARKDOWN"
					c.ReplyMarkup = &out[0]
					bot.Send(c)
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
				bot.Send(c)
				continue
			}
			if _, err := bot.Send(msg); err != nil {
				Log.Error(ctx, err.Error())
				// log.Println(err)
			}
			continue
		}

		//Stop responding from 12.29pm to 12.31pm or until dinner order has started (For occasional weird order timings)
		if time.Now().Unix() >= Processors.GetLunchTime().Unix()-60 &&
			(time.Now().Unix() <= Processors.GetLunchTime().Unix()+60 && !Processors.IsPollStart()) {
			if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Omw to order, wait for my good news! 🏃")); err != nil {
				Log.Error(ctx, err.Error())
				// log.Println(err)
			}
			continue
		}

		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			if startListenKey {
				//Capture key
				msg, _ := Bot.UpdateKey(ctx, update.Message.Chat.ID, update.Message.Text)
				if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg)); err != nil {
					Log.Error(ctx, err.Error())
					// log.Println(err)
				}
				startListenKey = false
				continue
			} else if startListenChope {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				ok := false
				msg.Text, ok = Bot.GetChope(ctx, update.Message.Chat.ID, update.Message.Text)
				if !ok {
					if _, err := bot.Send(msg); err != nil {
						Log.Error(ctx, err.Error())
						// log.Println(err)
					}
					continue
				}
				//Capture chope
				msg.ParseMode = "MARKDOWN"
				if _, err := bot.Send(msg); err != nil {
					Log.Error(ctx, err.Error())
					// log.Println(err)
				}
				startListenChope = false
				continue
			} else {
				continue
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "start":
			s, ok := Bot.CheckKey(ctx, update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				msg.Text = "Hello! " + update.Message.Chat.UserName
			}
		case "menu":
			s, ok := Bot.CheckKey(ctx, update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				txt, mp := Processors.OutputMenuWithButton(ctx, Bot.GetKey(ctx, update.Message.Chat.ID), update.Message.Chat.ID)
				for i, r := range txt {
					msg.Text = r
					if len(mp) > 0 {
						msg.ReplyMarkup = mp[i]
					}
					if _, err := bot.Send(msg); err != nil {
						Log.Error(ctx, err.Error())
						// log.Panic(err)
					}
				}
				continue
			}
		case "help":
			msg.Text = Bot.MakeHelpResponse()
			msg.ParseMode = "MARKDOWN"
		case "key":
			msg.Text, _ = Bot.CheckKey(ctx, update.Message.Chat.ID)
		case "newkey":
			msg.Text = "What's your key? \nGo to https://dinner.sea.com/accounts/token, copy the Key under Generate Auth Token and paste it here:"
			startListenKey = true
		case "status":
			s, ok := Bot.CheckKey(ctx, update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				msg.Text = Bot.ListWeeklyResultByUserId(ctx, update.Message.Chat.ID)
				msg.ParseMode = "HTML"
			}
		case "chope":
			msg.Text = "This command is deprecated. Choose from /menu instead!😋"
		case "choice":
			s, ok := Bot.CheckKey(ctx, update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				msg.Text, _ = Bot.CheckChope(ctx, update.Message.Chat.ID)
			}
		case "reminder":
			//Backdoor for test env
			if os.Getenv("TEST_DEPLOY") == "TRUE" || Common.Config.Adhoc {
				Bot.SendReminder(ctx)
			}
		case "mute":
			s, ok := Bot.CheckKey(ctx, update.Message.Chat.ID)
			if !ok {
				msg.Text = s
			} else {
				txt, kb := Bot.CheckMute(ctx, update.Message.Chat.ID)
				msg.Text = txt
				if kb != nil {
					msg.ReplyMarkup = kb[0]
				}
				msg.ParseMode = "MARKDOWN"
				if msgTrace, err := bot.Send(msg); err != nil {
					Log.Error(ctx, err.Error())
					// log.Panic(err)
				} else {
					//save msg id into cache for msg update
					cacheKey := fmt.Sprint(Common.USER_MUTE_MSG_ID_PREFIX, update.Message.Chat.ID)
					if err := Processors.RedisClient.Set(cacheKey, msgTrace.MessageID, 1800*time.Second).Err(); err != nil {
						Log.Error(ctx, "Mute | Error while writing to redis: %v", err.Error())
						// log.Printf("Mute | Error while writing to redis: %v", err.Error())
					} else {
						Log.Info(ctx, "Mute | Successful | Written %v to redis", cacheKey)
						// log.Printf("Mute | Successful | Written %v to redis", cacheKey)
					}
					continue
				}
			}
		case "checkin":
			//Backdoor for test env
			if os.Getenv("TEST_DEPLOY") == "TRUE" || Common.Config.Adhoc {
				Bot.SendCheckInLink(ctx)
			}
		case "delete":
			//Backdoor for test env
			if os.Getenv("TEST_DEPLOY") == "TRUE" || Common.Config.Adhoc {
				Bot.DeleteCheckInLink(ctx)
			}
		default:
			msg.Text = "I don't understand this command :("
		}
		if msg.Text != "" {
			if _, err := bot.Send(msg); err != nil {
				Log.Error(ctx, err.Error())
				// log.Panic(err)
			}
		}
	}
}
