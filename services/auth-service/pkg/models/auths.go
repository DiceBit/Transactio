package models

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type UserAuthRequest struct {
	Username string
	Email    string
	Password string
	Role     []string
}