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
	//Processors.ConnectDataBase()
	Processors.ConnectMySQL()

	//For testing only, update in config.toml
	if Processors.Config.Adhoc {
		r, donePrep = Processors.PrepOrder()
		time.Sleep(1 * time.Second)
		Processors.BatchOrderDinner(r)
		time.Sleep(1 * time.Second)
		Bot.SendNotifications()
		return
	}

	for {
		if Processors.IsWeekDay() && time.Now().Unix() == Processors.GetLunchTime().Unix()-7200 {
			Bot.SendNotifications()
		}

		if (Processors.IsWeekDay() && time.Now().Unix() >= Processors.GetLunchTime().Unix()-60 &&
			time.Now().Unix() <= Processors.GetLunchTime().Unix()-15) &&
			!donePrep {
			//get key and choice
			r, donePrep = Processors.PrepOrder()
		}

		if Processors.IsWeekDay() && time.Now().Unix() == Processors.GetLunchTime().Unix() {
			Processors.BatchOrderDinner(r)
			time.Sleep(1 * time.Second)
			//send notifications
			Bot.SendNotifications()
		}
	}
}
