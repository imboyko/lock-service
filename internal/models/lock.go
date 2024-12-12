package models

import "time"

type Lock struct {
	Id        string    `json:"id" redis:"id"`
	Timestamp time.Time `json:"timestamp" redis:"timestamp"`
	Username  string    `json:"username" redis:"username"`
}
