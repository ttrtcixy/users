package smtp

import (
	"crypto/tls"
	"fmt"
	"github.com/ttrtcixy/users/internal/config"
	"io"
	"net/smtp"
	"os"
)

const message = "From: %s\nTo: %s\nSubject: Hello\n\nToken: %s\n"

// todo обернуть ошибки

type SenderService struct {
	auth    smtp.Auth
	host    string
	addr    string
	from    string
	message string
}

func New(cfg *config.SmtpConfig) *SenderService {
	return &SenderService{
		auth:    smtp.PlainAuth("", cfg.Sender(), cfg.Password(), cfg.Host()),
		host:    cfg.Host(),
		addr:    cfg.Addr(),
		from:    cfg.Sender(),
		message: message,
	}
}

func (s *SenderService) DebugSend(to string, token string) error {
	const op = "SenderService.DebugSend"

	_, err := fmt.Fprintf(os.Stdout, s.message, s.from, to, token)
	if err != nil {
		return err
	}
	return nil
}

func (s *SenderService) Send(to string, token string) (err error) {
	const op = "SenderService.Send"
	client, err := s.newClient()
	if err != nil {
		return fmt.Errorf("%s: newClient failed: %w", op, err)
	}
	defer func() {
		if err != nil {
			_ = client.Close()
		}
	}()

	writer, err := s.prepareWriter(client, to)
	if err != nil {
		return fmt.Errorf("%s: prepareWriter failed: %w", op, err)
	}

	if err = s.writeMessage(writer, to, token); err != nil {
		return fmt.Errorf("%s: writeMessage failed: %w", op, err)
	}

	if err = client.Quit(); err != nil {
		return fmt.Errorf("%s: client.Quit failed: %w", op, err)
	}
	return nil
}

func (s *SenderService) newClient() (client *smtp.Client, err error) {
	tlsCfg := &tls.Config{ServerName: s.host}
	conn, err := tls.Dial("tcp", s.addr, tlsCfg)
	if err != nil {
		return nil, err
	}
	client, err = smtp.NewClient(conn, s.host)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return client, nil
}

func (s *SenderService) prepareWriter(client *smtp.Client, to string) (wc io.WriteCloser, err error) {
	if err = client.Auth(s.auth); err != nil {
		return nil, err
	}
	if err = client.Mail(s.from); err != nil {
		return nil, err
	}
	if err = client.Rcpt(to); err != nil {
		return nil, err
	}
	return client.Data()
}

func (s *SenderService) writeMessage(wc io.WriteCloser, to string, token string) error {
	msg := fmt.Sprintf(s.message, s.from, to, token)
	if _, err := wc.Write([]byte(msg)); err != nil {
		_ = wc.Close()
		return err
	}
	return wc.Close()
}
