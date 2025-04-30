package entities

type SignupRequest struct {
	Email    string
	Username string
	Password string
}

type SignupResponse struct {
	AccessToken  string
	RefreshToken string
}

type CheckLoginResponse struct {
	Status         bool
	UsernameExists bool
	EmailExists    bool
}
