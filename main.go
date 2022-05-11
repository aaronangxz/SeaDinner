package main

import (
	"fmt"
	"time"

	"github.com/aaronangxz/SeaDinner/Processors"
)

var (
	execPrep = false
	donePrep = false
	r        []Processors.UserRecords
)

func main() {
	Processors.LoadEnv()
	Processors.Init()
	Processors.ConnectDataBase()

	for {
		if !execPrep && (time.Now().Unix() < Processors.GetLunchTime().Unix()-300 || time.Now().Unix() > Processors.GetLunchTime().Unix()+300) {
			//Bot.InitBot()
			fmt.Println("exited")
			execPrep = true
		}

		if execPrep && !donePrep {
			//get key and choice
			r = Processors.PrepOrder()
			execPrep = false
			donePrep = true
		}

		if time.Now().Unix() == Processors.GetLunchTime().Unix() {
			Processors.BatchOrderDinner(r)
			return
		}
		//fmt.Println("looping")
	}
}
