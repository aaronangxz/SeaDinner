package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type DinnerMenu struct {
	status string
	dishes Food
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
	id          string
	name        string
	comment     string
	pollStart   string
	pollEnd     string
	servingTime string
	active      bool
}

func main() {
	fmt.Println("Hello, world.")

	// Create a Resty Client
	client := resty.New()
	//var AuthSuccess Menu
	resp, err := client.R().
		//SetQueryParams(map[string]string{
		//	"page_no": "1",
		//	"limit": "20",
		//	"sort":"name",
		//	"order": "asc",
		//	"random":strconv.FormatInt(time.Now().Unix(), 10),
		//}).
		SetHeader("Authorization", "Token e8c2f78d9a09bd8b59f83ef2ab6c0b22649798a9").
		SetHeader("Content-Type", "application/json").
		SetResult(). // or SetResult(AuthSuccess{}).
		//SetAuthToken("Token e8c2f78d9a09bd8b59f83ef2ab6c0b22649798a9").
		ForceContentType("application/json").
		Get("https://dinner.sea.com/api/current")

	if err != nil {
		fmt.Println(err)
	}

	//json.Unmarshal([]byte(resp), &bird)

	fmt.Println(resp.Result())

	//
	//// Sample of using Request.SetQueryString method
	//resp, err := client.R().
	//	SetQueryString("productId=232&template=fresh-sample&cat=resty&source=google&kw=buy a lot more").
	//	SetHeader("Accept", "application/json").
	//	SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F").
	//	Get("/show_product")

	// If necessary, you can force response content type to tell Resty to parse a JSON response into your struct
	//resp, err := client.R().
	//	SetResult(result).

	//	Get("v2/alpine/manifests/latest")
}
