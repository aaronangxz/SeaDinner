package main

import (
	"log"
	"time"

	"github.com/aaronangxz/SeaDinner/Bot"
	"github.com/aaronangxz/SeaDinner/Processors"
)

var (
	donePrep = false
	r        []Processors.UserChoiceWithKeyAndStatus
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
			//Test
			Processors.BatchOrderDinnerMultiThreadedWithWait(r)
		}

		if Processors.IsOrderTime() {
			for {
				if Processors.IsPollStart() {
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
