package Bot

type UserRecord struct {
	UserID int64  `json:"user_id"`
	Choice int    `json:"choice"`
	Key    string `json:"key"`
}
