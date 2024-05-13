package util

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var r *rand.Rand

func init() {
	s := rand.NewSource(time.Now().UnixNano())
	r = rand.New(s)
}

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}

	return string(b)
}

func GenerateRandomInt(min, max int) int {
	return r.Intn(max-min+1) + min
}
