package data

import (
	"math/rand"
	"time"
)

const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func IsUnixZero(t time.Time) bool {
	return t.Unix() == 0
}

// GetTimeUTCPointer gets utc time pointer
func GetTimeUTCPointer() *time.Time {
	t := time.Now().UTC()
	return &t
}

func TimeToPointer(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func UnixToTimePointer(v int64) *time.Time {
	if v == 0 {
		return nil
	}
	return TimeToPointer(time.Unix(v, 0))
}

// RandString generates random sting from characters in letters
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}
