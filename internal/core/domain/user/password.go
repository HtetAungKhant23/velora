package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort = errors.New("password: must have minimum 8 characters")
	ErrPasswordMismatch = errors.New("password: incorrect")
)

const bcryptCost = 12

type HashedPassword struct {
	hash []byte
}

func NewHashedPassword(plainText string) (HashedPassword, error) {
	if len(plainText) < 8 {
		return HashedPassword{}, ErrPasswordTooShort
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcryptCost)
	if err != nil {
		return HashedPassword{}, err
	}

	return HashedPassword{hash: hash}, nil
}

func (p HashedPassword) Verify(plainText string) error {
	if err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainText)); err != nil {
		return ErrPasswordMismatch
	}

	return nil
}

func ReconstitueHashedPassword(hash []byte) HashedPassword {
	return HashedPassword{hash: hash}
}
