package cache

import (
	"math/rand"
	"time"
)

// TTLWithJitter возвращает base + случайное значение от 0 до jitter.
// Это предотвращает одновременное истечение большого числа ключей.
func TTLWithJitter(base time.Duration, jitter time.Duration) time.Duration {
	if jitter <= 0 {
		return base
	}
	extra := time.Duration(rand.Int63n(int64(jitter) + 1))
	return base + extra
}
