package TestHelper

type UserKey struct {
	UserID  *int64  `json:"user_id"`
	UserKey *string `json:"user_key"`
	Ctime   *int64  `json:"ctime"`
	Mtime   *int64  `json:"mtime"`
}

func (u *UserKey) GetUserID() int64 {
	if u != nil && u.UserID != nil {
		return *u.UserID
	}
	return 0
}

func (u *UserKey) GetUserKey() string {
	if u != nil && u.UserKey != nil {
		return *u.UserKey
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
	UserID     int64 `json:"user_id"`
	UserChoice int64 `json:"user_choice"`
	Ctime      int64 `json:"ctime"`
	Mtime      int64 `json:"mtime"`
}
