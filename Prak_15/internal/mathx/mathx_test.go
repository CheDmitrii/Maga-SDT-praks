package mathx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSum_Table(t *testing.T) {
	cases := []struct{ a, b, want int }{
		{2, 3, 5}, {10, -5, 5}, {0, 0, 0},
	}
	for _, c := range cases {
		got := Sum(c.a, c.b)
		if got != c.want {
			t.Fatalf("Sum(%d,%d)=%d; want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestDivide(t *testing.T) {
	got, err := Divide(10, 2)
	require.NoError(t, err)
	assert.Equal(t, 5, got)

	_, err = Divide(10, 0)
	assert.Error(t, err)
}

func BenchmarkSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Sum(123, 456)
	}
}

func TestMustDivide_Panic(t *testing.T) {
	require.Panics(t, func() { MustDivide(1, 0) })
	// normal case does not panic and returns correct result
	got := MustDivide(10, 2)
	assert.Equal(t, 5, got)
}
