package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/aaronangxz/SeaDinner/Bot"
	"github.com/aaronangxz/SeaDinner/Processors"
)

var (
	execPrep = false
	donePrep = false
	r        []Processors.UserChoiceWithKey
	jobMutex sync.Mutex
)

func main() {
	Processors.LoadEnv()
	Processors.Init()
	Processors.ConnectDataBase()

	for {
		fmt.Println("wait")
		if execPrep && !donePrep {
			//get key and choice
			r, donePrep = Processors.PrepOrder()
			execPrep = !donePrep
		}

		if time.Now().Unix() == Processors.GetLunchTime().Unix() {
			Processors.BatchOrderDinner(r)
			time.Sleep(1 * time.Second)
			//send notifications
			Bot.SendNotifications()
		}

		if !execPrep && (time.Now().Unix() < Processors.GetLunchTime().Unix()-30 || time.Now().Unix() > Processors.GetLunchTime().Unix()+30) {
			fmt.Println("starting")
			jobMutex.Lock()
			Bot.InitBot()
			fmt.Println("exited")
			jobMutex.Unlock()
			execPrep = true
		}
		time.Sleep(1 * time.Second)
	}
}
