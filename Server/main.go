package main

import (
	"time"

	"github.com/aaronangxz/SeaDinner/Bot"
	"github.com/aaronangxz/SeaDinner/Processors"
)

var (
	donePrep = false
	r        []Processors.UserChoiceWithKey
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
		}

		if Processors.IsWeekDay(time.Now()) && time.Now().Unix() == Processors.GetLunchTime().Unix() {
			for {
				if time.Now().Unix() <= Processors.GetLunchTime().Unix()+180 {
					Processors.BatchOrderDinner(r)
					time.Sleep(time.Duration(Processors.Config.Runtime.BatchRetryCooldownSeconds) * time.Second)
					continue
				}
				break
			}
			Bot.SendNotifications()
		}
	}
}
