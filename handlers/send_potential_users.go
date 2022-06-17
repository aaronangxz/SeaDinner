package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"strconv"
	"strings"
	"time"
)

func SendPotentialUsers(ctx context.Context) {
	var (
		msgText = "Hey, I realise you chatted with me before but did not explore my features! /menu to try it out now ☺️"
	)

	txn := processors.App.StartTransaction("send_potential_user")
	defer txn.End()

	//Retrieve the whole set
	s, err := processors.RedisClient.SMembers(common.POTENTIAL_USER_SET).Result()
	if err != nil {
		log.Error(ctx, "SendPotentialUsers | Error while reading from redis: %v", err.Error())
		return
	}
	if s == nil {
		log.Warn(ctx, "SendPotentialUsers | Set is empty.")
		return
	}

	for _, pair := range s {
		//split <user_id>:<time> by ':'
		split := strings.Split(pair, ":")
		userID, _ := strconv.ParseInt(split[0], 10, 64)
		firstLoginTime, _ := strconv.ParseInt(split[1], 10, 64)

		if (time.Now().Unix() - firstLoginTime) < common.ONE_MONTH {
			log.Info(ctx, "SendPotentialUsers | Skip | Previous login time is not within range, must be < %v | user_id:%v", time.Now().Unix()-common.ONE_MONTH, userID)
			continue
		}

		if os.Getenv("TEST_DEPLOY") == "TRUE" {
			log.Info(ctx, "SendPotentialUsers | Test Flow | Skip Bot init")
		} else {
			bot, err := tgbotapi.NewBotAPI(common.GetTGToken(ctx))
			if err != nil {
				log.Error(ctx, err.Error())
			}
			bot.Debug = true
			log.Info(ctx, "Authorized on account %s", bot.Self.UserName)
			if _, err := bot.Send(tgbotapi.NewMessage(userID, msgText)); err != nil {
				log.Error(ctx, err.Error())
			}
		}

		log.Info(ctx, "SendPotentialUsers | Success | user_id:%v", userID)

		//Remove the old key and update with the new time in Set
		//As long as users do not give us the key, they will always be in the pool
		//We continuously update the time after each cold message to avoid annoyance
		if err := processors.RedisClient.SRem(common.POTENTIAL_USER_SET, pair).Err(); err != nil {
			log.Error(ctx, "SendPotentialUsers | Error while writing to redis: %v", err.Error())
		} else {
			log.Info(ctx, "SendPotentialUsers | Successful | Removed %v from potential_user set", pair)
		}

		toWrite := fmt.Sprint(userID, ":", time.Now().Unix())
		if err := processors.RedisClient.SAdd(common.POTENTIAL_USER_SET, toWrite).Err(); err != nil {
			log.Error(ctx, "SendPotentialUsers | Error while writing to redis: %v", err.Error())
		} else {
			log.Info(ctx, "SendPotentialUsers | Successful | Written %v to potential_user set", toWrite)
		}
	}
	log.Info(ctx, "SendPotentialUsers | Done")
}
