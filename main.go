package main

import (
	"fmt"

	"github.com/aaronangxz/SeaDinner/auth"
	"github.com/go-resty/resty/v2"
)

type DinnerMenu struct {
	status string
	dishes Food
}

type Current struct {
	Status  string `json:"status"`
	Details Menu   `json:"menu"`
}

type Food struct {
	code        string
	id          string
	name        string
	description string
	image       string
	ordered     int
	quota       int
	disabled    bool
}

type Menu struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	PollStart   string `json:"pollstart"`
	PollEnd     string `json:"pollend"`
	ServingTime string `json:"servingtime"`
	Active      bool   `json:"active"`
}

func main() {
	// Create a Resty Client
	client := resty.New()

	//get today's dinner info
	GetCurrent(*client)

}

func GetCurrent(client resty.Client) {
	var currentmenu Current

	_, err := client.R().
		SetHeader("Authorization", auth.Tokenauth()).
		SetResult(&currentmenu). // or SetResult(AuthSuccess{}).
		//SetAuthToken("Token e8c2f78d9a09bd8b59f83ef2ab6c0b22649798a9").
		//ForceContentType("application/json").
		Get("https://dinner.sea.com/api/current")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Query status of today's menu: %v\n\n", currentmenu.Status)
	fmt.Printf("ID: %v\n", currentmenu.Details.Id)
	fmt.Printf("%v\n", currentmenu.Details.Name)
	//fmt.Printf("%v\n", currentmenu.Details.Comment)
	fmt.Printf("Start: %v\n", currentmenu.Details.PollStart)
	fmt.Printf("End: %v\n", currentmenu.Details.PollEnd)
	fmt.Printf("Serving Time: %v\n", currentmenu.Details.ServingTime)
}
