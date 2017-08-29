package utils

import (
	"math/rand"
	"time"
)

// GenerateRandomString returns a random string.
func GenerateRandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]rune, n)
	for i := range b {
		ran := r.Intn(len(letters))
		b[i] = letters[ran]
	}
	return string(b)
}
