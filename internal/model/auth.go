package model

import "time"

type JwtToken string

type UpdatePassword struct {
	Username    Username
	OldPassword Password
	NewPassword Password
}

type TokenClaim struct {
	Username       Username
	Role           Role
	ExpirationTime time.Time
}
