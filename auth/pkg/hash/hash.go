package hash

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHandler interface {
	HashPass(password string) (string, error)
	HashPassCheck(hashedPassword, password string) bool
}

type SGA1Hasher struct {
	salt string
}

func NewSGA1Haser(salt string) *SGA1Hasher {
	return &SGA1Hasher{salt: salt}
}

// HashPass Хеширование пароля.
func (h *SGA1Hasher) HashPass(password string) (string, error) {
	pass := password + h.salt
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("ошибка при хешировании пароля")
	}
	return string(bytes), nil
}

// HashPassCheck Проверка пароль с хеш
func (h *SGA1Hasher) HashPassCheck(hashedPassword, password string) bool {
	pass := password + h.salt
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(pass))
	return err == nil
}
