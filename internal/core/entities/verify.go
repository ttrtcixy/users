package entities

type VerifyRequest struct {
	email string
	token string
}

type VerifyResponse struct {
	AccessToken  string
	RefreshToken string
}
