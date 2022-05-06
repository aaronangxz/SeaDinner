package main

import (
	"fmt"
	"time"

	"github.com/aaronangxz/SeaDinner/Processors"
	"github.com/go-resty/resty/v2"
)

func main() {
	Processors.LoadEnv()
	// Create a Resty Client
	client := resty.New()

	//get today's dinner info and retrieve today's ID
	ID := Processors.GetCurrent(*client)

	//get today's menu
	Processors.GetMenu(*client, ID)
	lunchTime := Processors.GetLunchTime()
	var req Processors.OrderRequest
	req.FoodID = 1194

	fmt.Println("Waiting..")

	// Processors.OrderDinner(*client, ID, req)
	for {
		if time.Now().Unix() == lunchTime.Unix() {
			for i := 0; i < 1; i++ {
				fmt.Println("Attempt ", i)
				Processors.OrderDinner(*client, ID, req)
			}
		}
	}
}
