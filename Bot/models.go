package Bot

type UserKey struct {
	UserID int64  `json:"user_id"`
	Key    string `json:"key"`
	Ctime  int64  `json:"ctime"`
	Mtime  int64  `json:"mtime"`
}

type UserChoice struct {
	UserID int64 `json:"user_id"`
	Choice int64 `json:"choice"`
	Ctime  int64 `json:"ctime"`
	Mtime  int64 `json:"mtime"`
}
