package main

import (
	"log"
	"time"

	"github.com/aaronangxz/SeaDinner/Bot"
	"github.com/aaronangxz/SeaDinner/Processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
)

var (
	donePrep = false
	r        []*sea_dinner.UserChoiceWithKey
	start    int64
	elapsed  int64
)

func main() {
	Processors.Init()
	Processors.InitClient()
	for {
		if Processors.IsSendReminderTime() {
			Bot.SendReminder()
		}

		if Processors.IsPrepOrderTime() && !donePrep {
			r, donePrep = Processors.PrepOrder()
		}

		if Processors.IsOrderTime() {
			for {
				if !Processors.IsPollStart() {
					start = time.Now().UnixMilli()
					Processors.BatchOrderDinnerMultiThreaded(r)
					elapsed = time.Now().UnixMilli() - start
					break
				}
				log.Println("Poll has not started, retrying.")
			}
			Bot.SendNotifications()
			log.Printf("Finished run | %v at %v in %vms",
				Processors.ConvertTimeStamp(time.Now().Unix()),
				Processors.ConvertTimeStampTime(time.Now().Unix()), elapsed)
		}
	}
}
