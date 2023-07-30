package model

import "time"

type User struct {
	ID             ID
	Username       string
	Password       string
	TimeOfCreation time.Time
	City           string
}
