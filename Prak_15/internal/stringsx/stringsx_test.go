package stringsx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClip(t *testing.T) {
	cases := []struct {
		name string
		s    string
		max  int
		want string
	}{
		{"empty", "", 5, ""},
		{"max=0", "hello", 0, ""},
		{"max<0", "hello", -1, ""},
		{"max==len", "hello", 5, "hello"},
		{"max>len", "hi", 5, "hi"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Clip(tc.s, tc.max)
			assert.Equal(t, tc.want, got)
		})
	}
}
