package common

import "time"

//goland:noinspection ALL
const (
	DB_ORDER_LOG_TAB   = "order_log_tab"
	DB_USER_CHOICE_TAB = "user_choice_tab"
	DB_USER_KEY_TAB    = "user_key_tab"

	MENU_CACHE_KEY_PREFIX   = "current_menu:"
	DAY_ID_KEY_PREFIX       = "day_id:"
	USER_KEY_PREFIX         = "user_key:"
	USER_CHOICE_PREFIX      = "user_choice:"
	USER_MUTE_MSG_ID_PREFIX = "user_mute:"

	POTENTIAL_USER_SET = "potential_user"
	CHECK_IN_LINK_SET  = "checkin_link"
)

//goland:noinspection ALL
const (
	ONE_HOUR  = int64(3600 * time.Second)
	ONE_DAY   = 24 * ONE_HOUR
	ONE_WEEK  = 7 * ONE_DAY
	ONE_MONTH = 30 * ONE_WEEK
)
