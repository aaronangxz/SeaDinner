package main

import (
	"github.com/aaronangxz/SeaDinner/Processors"
	"github.com/go-resty/resty/v2"
)

func main() {
	// Create a Resty Client
	client := resty.New()

	//get today's dinner info
	// Processors.GetCurrent(*client)

	//get today's menu
	Processors.GetMenu(*client)

}
