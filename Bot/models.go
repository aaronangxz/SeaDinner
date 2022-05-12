package Bot

type UserKey struct {
	UserID *int64  `json:"user_id"`
	Key    *string `json:"key"`
	Ctime  *int64  `json:"ctime"`
	Mtime  *int64  `json:"mtime"`
}

func (u *UserKey) GetUserID() int64 {
	if u != nil && u.UserID != nil {
		return *u.UserID
	}
	return 0
}

func (u *UserKey) GetKey() string {
	if u != nil && u.Key != nil {
		return *u.Key
	}
	return ""
}

func (u *UserKey) GetCtime() int64 {
	if u != nil && u.Ctime != nil {
		return *u.Ctime
	}
	return 0
}
func (u *UserKey) GetMtime() int64 {
	if u != nil && u.Mtime != nil {
		return *u.Mtime
	}
	return 0
}

type UserChoice struct {
	UserID int64 `json:"user_id"`
	Choice int64 `json:"choice"`
	Ctime  int64 `json:"ctime"`
	Mtime  int64 `json:"mtime"`
}
