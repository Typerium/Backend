package password

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type bcryptProcessor struct {
	cost int
}

func NewBcryptProcessor(cost int) Processor {
	return &bcryptProcessor{cost: cost}
}

func (p *bcryptProcessor) Encode(password string) (string, error) {
	encodedPass, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(encodedPass), nil
}

func (p *bcryptProcessor) Equal(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}

		return false, errors.WithStack(err)
	}

	return true, nil
}
