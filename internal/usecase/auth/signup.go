package authusecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/entities"
	"github.com/ttrtcixy/users/internal/errors"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/usecase/ports"
	"io"
	"net/smtp"
	"os"
)

type SignupUseCase struct {
	log  logger.Logger
	cfg  *config.UsecaseConfig
	repo ports.SignupRepository

	smtp *Smtp
}

func NewSignup(ctx context.Context, log logger.Logger, cfg *config.UsecaseConfig, repo ports.SignupRepository) *SignupUseCase {
	return &SignupUseCase{
		log:  log,
		repo: repo,
		cfg:  cfg,
		smtp: NewSmtp(cfg),
	}
}

func (u *SignupUseCase) Signup(ctx context.Context, payload *entities.SignupRequest) (*entities.SignupResponse, error) {
	exists, err := u.repo.CheckLoginExist(ctx, payload)
	if err != nil {
		return nil, err
	}

	if exists.Status {
		var err = &apperrors.ErrLoginExists{}
		if exists.UsernameExists {
			err.Username = payload.Username
		}
		if exists.EmailExists {
			err.Email = payload.Email
		}
		return nil, err
	}

	hash, err := u.hash(payload.Password)
	if err != nil {
		return nil, err
	}
	_ = hash

	err = u.smtp.DebugSend(payload.Email)
	if err != nil {
		u.log.Error(err.Error())
		return nil, err
	}
	// todo add new user
	// todo check user email
	// todo generate and send tokens

	return nil, nil
}

// hash generates a hash of the password using a salt
func (u *SignupUseCase) hash(password string) (string, error) {
	salt, err := u.salt()
	if err != nil {
		return "", err
	}

	hasher := sha256.New()

	bytePassword := append([]byte(password), salt...)
	hasher.Write(bytePassword)

	hash := hasher.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash), nil
}

// salt generate random salt for password
func (u *SignupUseCase) salt() ([]byte, error) {
	salt := make([]byte, u.cfg.PasswordSaltLength())
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

const message = "From: %s\nTo: %s\nSubject: Hello\n\nThis is a test email.\n"

type Smtp struct {
	auth    smtp.Auth
	host    string
	addr    string
	from    string
	message string
}

func NewSmtp(cfg *config.UsecaseConfig) *Smtp {
	return &Smtp{
		auth:    smtp.PlainAuth("", cfg.SMTPLogin(), cfg.SMTPPassword(), cfg.SMTPHost()),
		host:    cfg.SMTPHost(),
		addr:    cfg.SMTPAddr(),
		from:    cfg.SMTPSender(),
		message: message,
	}
}

func (s *Smtp) DebugSend(to string) error {
	_, err := fmt.Fprintf(os.Stdout, s.message, s.from, to)
	if err != nil {
		return err
	}
	return nil
}

func (s *Smtp) Send(to string) error {
	const op = "usecase.Signup.Smtp"

	client, err := s.newClient()
	if err != nil {
		return fmt.Errorf("%s: newClient failed: %w", op, err)
	}

	writer, err := s.prepareWriter(client, to)
	if err != nil {
		_ = client.Close()
		return fmt.Errorf("%s: prepareWriter failed: %w", op, err)
	}

	if err := s.writeMessage(writer, to); err != nil {
		_ = client.Close()
		return fmt.Errorf("%s: writeMessage failed: %w", op, err)
	}

	if err := client.Quit(); err != nil {
		_ = client.Close()
		return fmt.Errorf("%s: client.Quit failed: %w", op, err)
	}
	return nil
}

func (s *Smtp) newClient() (*smtp.Client, error) {
	tlsCfg := &tls.Config{ServerName: s.host}
	conn, err := tls.Dial("tcp", s.addr, tlsCfg)
	if err != nil {
		return nil, err
	}
	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return client, nil
}

func (s *Smtp) prepareWriter(client *smtp.Client, to string) (io.WriteCloser, error) {
	if err := client.Auth(s.auth); err != nil {
		return nil, err
	}
	if err := client.Mail(s.from); err != nil {
		return nil, err
	}
	if err := client.Rcpt(to); err != nil {
		return nil, err
	}
	return client.Data()
}

func (s *Smtp) writeMessage(wc io.WriteCloser, to string) error {
	msg := fmt.Sprintf(s.message, s.from, to)
	if _, err := wc.Write([]byte(msg)); err != nil {
		_ = wc.Close()
		return err
	}
	return wc.Close()
}
