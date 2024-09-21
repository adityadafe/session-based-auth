package db

import "time"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	SessionId string
	Ttx       time.Time
}

func (s Session) IsExpired() bool {
	return s.Ttx.Before(time.Now())
}
