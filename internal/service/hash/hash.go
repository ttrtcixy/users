package hash

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	apperrors "github.com/ttrtcixy/users/internal/errors"
)

type HasherService struct {
}

func New() *HasherService {
	return &HasherService{}
}

// Hash generates a hash using sha256 and return base64 string
func (h *HasherService) Hash(str string) (hash string, err error) {
	const op = "HasherService.Hash"

	hasher := sha256.New()
	if _, err = hasher.Write([]byte(str)); err != nil {
		return "", apperrors.Wrap(op, err)
	}
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil)), nil
}

// Salt generate random salt
func (h *HasherService) Salt(length int) ([]byte, error) {
	const op = "HasherService.Salt"

	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return nil, apperrors.Wrap(op, err)
	}
	return salt, nil
}

// HashWithSalt generates a hash with salt using sha256 and return base64 string
func (h *HasherService) HashWithSalt(str string, salt []byte) (hash string, err error) {
	const op = "HasherService.HashWithSalt"

	hasher := sha256.New()
	data := append([]byte(str), salt...)
	if _, err = hasher.Write(data); err != nil {
		return "", apperrors.Wrap(op, err)
	}

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil)), nil
}

// ComparePasswords -
func (h *HasherService) ComparePasswords(storedHash, password string, salt string) (bool, error) {
	const op = "HasherService.ComparePasswords"
	byteSalt, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return false, apperrors.Wrap(op, err)
	}

	computedHash, err := h.HashWithSalt(password, byteSalt)
	if err != nil {
		return false, apperrors.Wrap(op, err)
	}

	return storedHash == computedHash, nil
}
