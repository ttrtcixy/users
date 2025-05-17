package entities

type User struct {
	ID           int64
	Username     string
	Email        string
	Password     string
	IsActive     bool
	RoleId       int
	PasswordHash string
	PasswordSalt string
}
