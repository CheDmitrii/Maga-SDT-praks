package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubRepo struct {
	users map[string]User
}

func (r stubRepo) ByEmail(email string) (User, error) {
	u, ok := r.users[email]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

func TestFindIDByEmail(t *testing.T) {
	repo := stubRepo{
		users: map[string]User{
			"a@example.com": {ID: 1, Email: "a@example.com"},
		},
	}
	service := New(repo)

	t.Run("found", func(t *testing.T) {
		id, err := service.FindIDByEmail("a@example.com")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})

	t.Run("not_found", func(t *testing.T) {
		id, err := service.FindIDByEmail("unknown@example.com")
		assert.Equal(t, int64(0), id)
		assert.Equal(t, ErrNotFound, err)
	})
}
