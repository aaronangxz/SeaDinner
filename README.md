<h2>Now</h2>
GET 
current Returns the menu details of the 'current menu'.<br>
This would be the same menu as the 'Order' link on the website.<br>

``` Go
type Menu struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	PollStart   string `json:"pollstart"`
	PollEnd     string `json:"pollend"`
	ServingTime string `json:"servingtime"`
	Active      bool   `json:"active"`
}
```

<h2>List</h2>
GET list Returns a list of menus between the stipulated timing.


<h2>Menu</h2>
GET menu/:menu_id List the food items belonging to the menu. Also includes the details of the menu.

``` Go
type Food struct {
	Code        string `json:"code"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Ordered     int    `json:"ordered"`
	Quota       int    `json:"quota"`
	Disabled    bool   `json:"disabled"`
}
```

<h2>Ordering</h2>
GET order/:menu_id Returns the food item you have ordered for the menu.<br>
POST order/:menu_id Makes an order for the food item you have specified.<br>
DELETE order/:menu_id Delete the food item you specified from the menu.