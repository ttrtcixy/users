package entities

type SignupRequest struct {
	Email    string
	Username string
	Password string
}

// todo refactor
type SignupResponse struct {
}

type CheckLoginResponse struct {
	Status         bool
	UsernameExists bool
	EmailExists    bool
}

type CreateUserRequest struct {
	Username        string
	Email           string
	PasswordHash    string
	PasswordSalt    string
	ActivationToken string
}
