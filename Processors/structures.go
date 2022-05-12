package Processors

const (
	// UrlPrefix   = "https://dinner.sea.com"
	// TokenPrefix = "Token "

	URL_CURRENT = 0
	URL_MENU    = 1
	URL_ORDER   = 2
)

var Constant_URL_type = map[int32]string{
	0: "URL_CURRENT",
	1: "URL_MENU",
	2: "URL_ORDER",
}

func Int(v int) *int { return &v }

type DinnerMenu struct {
	Status string `json:"status"`
	Dishes Food   `json:"food"`
}

type DinnerMenuArr struct {
	Status    string `json:"status"`
	DinnerArr []Food `json:"food"`
}

type Current struct {
	Status  string `json:"status"`
	Details Menu   `json:"menu"`
}

type OrderRequest struct {
	FoodID int `json:"food_id"`
}

type OrderResponse struct {
	Status   *string `json:"status"`
	Selected *int    `json:"selected"`
	Error    *string `json:"error"`
}

func (o *OrderResponse) GetStatus() string {
	if o != nil && o.Status != nil {
		return *o.Status
	}
	return ""
}

func (o *OrderResponse) GetSelected() int {
	if o != nil && o.Selected != nil {
		return *o.Selected
	}
	return 0
}

func (o *OrderResponse) GetError() string {
	if o != nil && o.Error != nil {
		return *o.Error
	}
	return ""
}

type Food struct {
	Code        string `json:"code"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Ordered     int    `json:"ordered"`
	Quota       int    `json:"quota"`
	Disabled    bool   `json:"disabled"`
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

type UserChoice struct {
	UserID int64 `json:"user_id"`
	Choice int64 `json:"choice"`
	Ctime  int64 `json:"ctime"`
	Mtime  int64 `json:"mtime"`
}

type UserChoiceWithKey struct {
	UserID int64  `json:"user_id"`
	Key    string `json:"string"`
	Choice int64  `json:"choice"`
	Ctime  int64  `json:"ctime"`
	Mtime  int64  `json:"mtime"`
}
