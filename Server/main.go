package main

import (
	"time"

	"github.com/aaronangxz/SeaDinner/Bot"
	"github.com/aaronangxz/SeaDinner/Processors"
)

var (
	donePrep = false
	r        []Processors.UserChoiceWithKey
	// jobMutex sync.Mutex
)

func main() {
	Processors.LoadEnv()
	Processors.Init()
	Processors.ConnectDataBase()

	for {
		if (time.Now().Unix() >= Processors.GetLunchTime().Unix()-60 &&
			time.Now().Unix() <= Processors.GetLunchTime().Unix()-15) &&
			!donePrep {
			//get key and choice
			r, donePrep = Processors.PrepOrder()
		}

		if time.Now().Unix() == Processors.GetLunchTime().Unix() {
			Processors.BatchOrderDinner(r)
			time.Sleep(1 * time.Second)
			//send notifications
			Bot.SendNotifications()
		}

		// if !execPrep && (time.Now().Unix() < Processors.GetLunchTime().Unix()-30 || time.Now().Unix() > Processors.GetLunchTime().Unix()+30) {
		// 	fmt.Println("starting")
		// 	jobMutex.Lock()
		// 	Bot.InitBot()
		// 	fmt.Println("exited")
		// 	jobMutex.Unlock()
		// 	execPrep = true
		// }
		// time.Sleep(1 * time.Second)
	}
}
