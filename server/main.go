package main

import (
	handlers "github.com/aaronangxz/SeaDinner/handlers"
	"os"
	"time"

	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
)

var (
	donePrep = false
	r        []*sea_dinner.UserChoiceWithKey
	start    int64
	elapsed  int64
)

func main() {
	processors.Init()
	processors.InitClient()
	go processors.MenuRefresher(log.NewCtx())

	//For adhoc use only
	//processors.SendAdHocNotification(0,"")

	for {
		ctx := log.NewCtx()
		if processors.IsSendReminderTime() {
			if processors.IsSOW(time.Now()) {
				processors.StoreFoodMappings(ctx)
			}
			handlers.SendReminder(ctx)
			handlers.SendPotentialUsers(ctx)
		}

		if processors.IsPrepOrderTime() && !donePrep {
			r, donePrep = processors.PrepOrder(ctx)
			//Test
			go processors.BatchOrderDinnerMultiThreadedWithWait(ctx, r)
		}

		if processors.IsOrderTime() {
			for {
				if processors.IsPollStart() {
					start = time.Now().UnixMilli()
					processors.BatchOrderDinnerMultiThreaded(ctx, r)
					elapsed = time.Now().UnixMilli() - start
					break
				}
				log.Warn(ctx, "Poll has not started, retrying.")
			}
			handlers.SendNotifications(ctx)
			log.Info(ctx, "Finished run | %v at %v in %vms",
				processors.ConvertTimeStamp(time.Now().Unix()),
				processors.ConvertTimeStampTime(time.Now().Unix()), elapsed)
		}

		if os.Getenv("SEND_CHECKIN") == "TRUE" {
			if processors.IsSendCheckInTime() {
				handlers.SendCheckInLink(ctx)
			}

			if processors.IsDeleteCheckInTime() {
				handlers.DeleteCheckInLink(ctx)
			}
		}
	}
}
