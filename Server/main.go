package main

import (
	"log"
	"time"

	"github.com/aaronangxz/SeaDinner/Bot"
	"github.com/aaronangxz/SeaDinner/Processors"
)

var (
	donePrep    = false
	r           []Processors.UserChoiceWithKeyAndStatus
	records     []Processors.OrderRecord
	start       int64
	elapsed     int64
	totalOrders int
)

func main() {
	Processors.LoadEnv()
	Processors.Init()
	//For testing only, update in config.toml
	if Processors.Config.Adhoc {
		Processors.ConnectTestMySQL()
	} else {
		Processors.ConnectMySQL()
	}

	for {
		if Processors.IsWeekDay(time.Now()) && time.Now().Unix() == Processors.GetLunchTime().Unix()-7200 {
			Bot.SendReminder()
		}

		if (Processors.IsWeekDay(time.Now()) && time.Now().Unix() >= Processors.GetLunchTime().Unix()-60 &&
			time.Now().Unix() <= Processors.GetLunchTime().Unix()-15) &&
			!donePrep {
			//get key and choice
			r, donePrep = Processors.PrepOrder()
			totalOrders = len(r)
		}

		if Processors.IsWeekDay(time.Now()) && time.Now().Unix() == Processors.GetLunchTime().Unix() {
			for {
				if time.Now().Unix() <= Processors.GetLunchTime().Unix()+int64(Processors.Config.Runtime.RetryOffsetSeconds) {
					start = time.Now().UnixMilli()
					records = Processors.BatchOrderDinner(&r)
					if len(r) == 0 {
						elapsed = time.Now().UnixMilli() - start
						log.Println("Successfully processed all orders.")
						break
					}
					time.Sleep(time.Duration(Processors.Config.Runtime.BatchRetryCooldownSeconds) * time.Second)
					continue
				}
				if elapsed == 0 {
					elapsed = time.Now().UnixMilli() - start
				}
				Processors.UpdateOrderLog(records)
				Processors.OutputResultsCount(totalOrders, len(records))
				break
			}
			Bot.SendNotifications()
			log.Printf("Finished run | %v at %v in %vms", Processors.ConvertTimeStamp(time.Now().Unix()), Processors.ConvertTimeStampTime(time.Now().Unix()), elapsed)
		}
	}
}
