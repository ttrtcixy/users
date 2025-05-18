package ports

type HasherService interface {
	Hash(str string) (hash string, err error)
	Salt(length int) ([]byte, error)
	HashWithSalt(str string, salt []byte) (hash string, err error)
	ComparePasswords(storedHash, password, salt string) (bool, error)
}
