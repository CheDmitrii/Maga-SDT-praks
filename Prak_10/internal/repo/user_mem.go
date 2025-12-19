package repo

import (
	"errors"
	"golang.org/x/crypto/bcrypt"

	"Prak_10/internal/core"
)

type UserRecord struct {
	ID    int64
	Email string
	Role  string
	Hash  []byte
}

type UserMem struct{ users map[string]UserRecord } // key = email

func NewUserMem() *UserMem {
	// заранее захэшированные пароли (пример: "secret123")
	hash := func(s string) []byte { h, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost); return h }
	return &UserMem{users: map[string]UserRecord{
		"admin@example.com": {ID: 1, Email: "admin@example.com", Role: "admin", Hash: hash("secret123")},
		"user@example.com":  {ID: 2, Email: "user@example.com", Role: "user", Hash: hash("secret123")},
	}}
}

var ErrNotFound = errors.New("user not found")
var ErrBadCreds = errors.New("bad credentials")

func (r *UserMem) ByEmail(email string) (UserRecord, error) {
	u, ok := r.users[email]
	if !ok {
		return UserRecord{}, ErrNotFound
	}
	return u, nil
}

func (r *UserMem) CheckPassword(email, pass string) (core.User, error) {
	rec, err := r.ByEmail(email)
	if err != nil {
		return core.User{}, ErrNotFound
	}

	if bcrypt.CompareHashAndPassword(rec.Hash, []byte(pass)) != nil {
		return core.User{}, ErrBadCreds
	}

	return core.User{
		ID:    rec.ID,
		Email: rec.Email,
		Role:  rec.Role,
	}, nil
}
