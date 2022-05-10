package main

import (
	"fmt"
	"time"

	"github.com/aaronangxz/SeaDinner/Bot"
	"github.com/aaronangxz/SeaDinner/Processors"
)

var (
	execPrep = false
	r        []Processors.UserRecords
)

func main() {
	Processors.LoadEnv()
	Processors.Init()
	Processors.ConnectDataBase()

	for {
		if !execPrep && (time.Now().Unix() < Processors.GetLunchTime().Unix()-300 || time.Now().Unix() > Processors.GetLunchTime().Unix()+300) {
			Bot.InitBot()
			fmt.Println("exited")
			execPrep = true
		}

		if execPrep {
			//get key and choice
			r = Processors.PrepOrder()
			execPrep = false
		}

		if time.Now().Unix() == Processors.GetLunchTime().Unix() {
			for i := 0; i < 1; i++ {
				fmt.Println("Attempt ", i)
				Processors.BatchOrderDinner(r)
			}
		}
		fmt.Println("looping")
	}
}
