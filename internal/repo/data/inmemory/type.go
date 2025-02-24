package inmemory

import "time"

type User struct {
	UID      string `json:"uid"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Token struct {
	Value     string    `json:"value"`
	UserUID   string    `json:"userUid"`
	IsRefresh bool      `json:"isRefresh"`
	Expired   time.Time `json:"expired"`
}
