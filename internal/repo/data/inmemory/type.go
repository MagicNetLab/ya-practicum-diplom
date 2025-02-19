package inmemory

import "time"

type User struct {
	Uid      string `json:"uid"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Token struct {
	Value     string    `json:"value"`
	UserUid   string    `json:"userUid"`
	IsRefresh bool      `json:"isRefresh"`
	Expired   time.Time `json:"expired"`
}
