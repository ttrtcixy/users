package ports

type SmtpService interface {
	Send(to string, token string) error
	DebugSend(to string, token string) error
}
