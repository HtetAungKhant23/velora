package user

import (
	"time"

	"github.com/google/uuid"
)

type UserID string

func newUserID() UserID {
	return UserID(uuid.New().String())
}

type User struct {
	id        UserID
	email     Email
	password  HashedPassword
	createdAt time.Time
	updatedAt time.Time
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) PasswordHash() []byte {
	return u.password.hash
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func NewUser(rawEmail, plainPassword string) (*User, error) {
	email, err := NewEmail(rawEmail)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := NewHashedPassword(plainPassword)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &User{
		id:        newUserID(),
		email:     email,
		password:  hashedPassword,
		createdAt: now,
		updatedAt: now,
	}

	return user, nil
}

func (u *User) VerifyPassword(plainText string) error {
	if err := u.password.Verify(plainText); err != nil {
		return err
	}
	return nil
}

func ReconstitueUser(
	id UserID,
	email Email,
	passwordHash []byte,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:        id,
		email:     email,
		password:  ReconstitueHashedPassword(passwordHash),
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
