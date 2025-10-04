package utils

import (
	"crypto/rand"
	"encoding/base64"
	mathrand "math/rand"
	"strconv"
)

// GenerateRandomNumber generates a random 6-digit(a random number between 100000 and 999999 ) number as a string.
func GenerateRandomNumber() string {
	num := mathrand.Intn(900000) + 100000
	return strconv.Itoa(num)
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}
