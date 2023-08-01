package model

import "time"

type User struct {
	ID             ID
	Username       string
	Password       string
	Role           Role
	TimeOfCreation time.Time
	City           string
	Version        int
}
