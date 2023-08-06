package model

import "time"

type User struct {
	ID             ID
	Username       Username
	Password       HashedPass
	Role           Role
	TimeOfCreation time.Time
	City           string
	Version        int32
}

type UserInfo struct {
	ID             ID
	Username       Username
	Password       Password
	Role           Role
	TimeOfCreation time.Time
	City           string
	Version        int32
}
