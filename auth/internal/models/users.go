package models

import "time"

type Users struct {
	Id       string
	Email    string
	Password string
	Date     time.Time
}
