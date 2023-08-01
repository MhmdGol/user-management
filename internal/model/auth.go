package model

type LoginRequest struct {
	Username string
	Password string
}

type JwtToken struct {
	Token string
}

type UpdatePassword struct {
	Username    string
	OldPassword string
	NewPassword string
}
