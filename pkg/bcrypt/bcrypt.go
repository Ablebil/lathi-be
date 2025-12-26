package bcrypt

import (
	"sync"

	b "golang.org/x/crypto/bcrypt"
)

type BcryptItf interface {
	Hash(plain string) (string, error)
	Compare(password, hashed string) bool
}

type bcrypt struct{}

var (
	instance BcryptItf
	once     sync.Once
)

func NewBcrypt() BcryptItf {
	once.Do(func() {
		instance = &bcrypt{}
	})
	return instance
}

func (bc *bcrypt) Hash(plain string) (string, error) {
	bytes, err := b.GenerateFromPassword([]byte(plain), b.DefaultCost)
	return string(bytes), err
}

func (bc *bcrypt) Compare(password, hashed string) bool {
	err := b.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
