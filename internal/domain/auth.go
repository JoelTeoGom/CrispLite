package domain

type RegisterResponse struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

type RefreshResponse struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}
