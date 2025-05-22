package entities

type VerifyRequest struct {
	JwtEmailToken string
}

type VerifyResponse struct {
	AccessToken  string
	RefreshToken string
	ClientUUID   string
}
