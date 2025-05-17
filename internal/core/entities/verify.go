package entities

type VerifyRequest struct {
	JwtToken string
}

type VerifyResponse struct {
	AccessToken  string
	RefreshToken string
}
