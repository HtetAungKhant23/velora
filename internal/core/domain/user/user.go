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
