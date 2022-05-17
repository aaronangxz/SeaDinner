package Processors

const (
	// UrlPrefix   = "https://dinner.sea.com"
	// TokenPrefix = "Token "

	URL_CURRENT = 0
	URL_MENU    = 1
	URL_ORDER   = 2

	ORDER_STATUS_OK   = 0
	ORDER_STATUS_FAIL = 1

	DB_ORDER_LOG_TAB   = "order_log_tab"
	DB_USER_CHOICE_TAB = "user_choice_tab"
	DB_USER_KEY_TAB    = "user_key_tab"
)

var Constant_URL_type = map[int32]string{
	0: "URL_CURRENT",
	1: "URL_MENU",
	2: "URL_ORDER",
}

func Int(v int) *int          { return &v }
func Int64(v int64) *int64    { return &v }
func String(s string) *string { return &s }
func Bool(b bool) *bool       { return &b }

type DinnerMenu struct {
	Status string `json:"status"`
	Dishes Food   `json:"food"`
}

type DinnerMenuArr struct {
	Status    *string `json:"status"`
	DinnerArr []Food  `json:"food"`
}

func (d *DinnerMenuArr) GetStatus() string {
	if d != nil && d.Status != nil {
		return *d.Status
	}
	return ""
}

type Current struct {
	Status *string `json:"status"`
	Menu   Details `json:"menu"`
}
type OrderRequest struct {
	FoodID int `json:"food_id"`
}

type OrderResponse struct {
	Status     *string `json:"status"`
	StatusCode *int    `json:"status_code"`
	Selected   *int    `json:"selected"`
	Error      *string `json:"error"`
}

func (o *OrderResponse) GetStatus() string {
	if o != nil && o.Status != nil {
		return *o.Status
	}
	return ""
}

func (o *OrderResponse) GetStatusCode() int {
	if o != nil && o.StatusCode != nil {
		return *o.StatusCode
	}
	return 0
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

type Details struct {
	Id          *int    `json:"id"`
	Name        *string `json:"name"`
	Comment     *string `json:"comment"`
	PollStart   *string `json:"poll_start"`
	PollEnd     *string `json:"poll_end"`
	ServingTime *string `json:"serving_time"`
	Active      *bool   `json:"active"`
	VenueId     *int    `json:"venue_id"`
}

func (m *Details) GetId() int {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Details) GetPollStart() string {
	if m != nil && m.Id != nil {
		return *m.PollStart
	}
	return ""
}

type UserChoice struct {
	UserID     int64 `json:"user_id"`
	UserChoice int64 `json:"user_choice"`
	Ctime      int64 `json:"ctime"`
	Mtime      int64 `json:"mtime"`
}

type UserChoiceWithKeyAndStatus struct {
	UserID     *int64  `json:"user_id"`
	UserKey    *string `json:"user_key"`
	UserChoice *string `json:"user_choice"`
	Ctime      *int64  `json:"ctime"`
	Mtime      *int64  `json:"mtime"`
	IsSuccess  *bool   `json:"is_success"`
}

func (u *UserChoiceWithKeyAndStatus) GetUserID() int64 {
	if u != nil && u.UserID != nil {
		return *u.UserID
	}
	return 0
}

func (u *UserChoiceWithKeyAndStatus) GetUserKey() string {
	if u != nil && u.UserKey != nil {
		return *u.UserKey
	}
	return ""
}

func (u *UserChoiceWithKeyAndStatus) SetUserKey(key *string) {
	u.UserKey = key
}

func (u *UserChoiceWithKeyAndStatus) GetUserChoice() string {
	if u != nil && u.UserChoice != nil {
		return *u.UserChoice
	}
	return ""
}

func (u *UserChoiceWithKeyAndStatus) GetIsSuccess() bool {
	if u != nil && u.IsSuccess != nil {
		return *u.IsSuccess
	}
	return false
}

type OrderRecord struct {
	ID        *int64  `json:"id"`
	UserID    *int64  `json:"user_id"`
	FoodID    *string `json:"food_id"`
	OrderTime *int64  `json:"order_time"`
	Status    *int64  `json:"status"`
	ErrorMsg  *string `json:"error_msg"`
}

func (o *OrderRecord) GetUserID() int64 {
	if o != nil && o.UserID != nil {
		return *o.UserID
	}
	return 0
}

func (o *OrderRecord) GetFoodID() string {
	if o != nil && o.FoodID != nil {
		return *o.FoodID
	}
	return ""
}

func (o *OrderRecord) GetOrderTime() int64 {
	if o != nil && o.OrderTime != nil {
		return *o.OrderTime
	}
	return 0
}

func (o *OrderRecord) GetStatus() int64 {
	if o != nil && o.Status != nil {
		return *o.Status
	}
	return 0
}

func (o *OrderRecord) GetErrorMsg() string {
	if o != nil && o.ErrorMsg != nil {
		return *o.ErrorMsg
	}
	return ""
}
