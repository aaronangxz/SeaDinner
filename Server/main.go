package main

import (
	"os"
	"time"

	"github.com/aaronangxz/SeaDinner/Bot"
	"github.com/aaronangxz/SeaDinner/Log"
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
	go Processors.MenuRefresher(Log.NewCtx())

	//For adhoc use only
	//Processors.SendAdHocNotification(0,"")

	for {
		ctx := Log.NewCtx()
		if Processors.IsSendReminderTime() {
			Bot.SendReminder(ctx)
		}

		if Processors.IsPrepOrderTime() && !donePrep {
			r, donePrep = Processors.PrepOrder(ctx)
			//Test
			go Processors.BatchOrderDinnerMultiThreadedWithWait(ctx, r)
		}

		if Processors.IsOrderTime() {
			for {
				if Processors.IsPollStart() {
					start = time.Now().UnixMilli()
					Processors.BatchOrderDinnerMultiThreaded(ctx, r)
					elapsed = time.Now().UnixMilli() - start
					break
				}
				Log.Warn(ctx, "Poll has not started, retrying.")
			}
			Bot.SendNotifications(ctx)
			Log.Info(ctx, "Finished run | %v at %v in %vms",
				Processors.ConvertTimeStamp(time.Now().Unix()),
				Processors.ConvertTimeStampTime(time.Now().Unix()), elapsed)
		}

		if os.Getenv("SEND_CHECKIN") == "TRUE" {
			if Processors.IsSendCheckInTime() {
				Bot.SendCheckInLink(ctx)
			}

			if Processors.IsDeleteCheckInTime() {
				Bot.DeleteCheckInLink(ctx)
			}
		}
	}
}
