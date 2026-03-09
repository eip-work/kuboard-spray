package common

import (
	"crypto/rand"
	"math/big"
)

const (
	letters        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers        = "0123456789"
	specialChars   = "!@#$%^&*()-_=+[]{}|;:,.<>?/~`"
	allChars       = letters + numbers + specialChars
	passwordLength = 12 // 你可以根据需要调整密码长度
)

func GeneratePassword(length int) string {
	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "0123456789"
		}
		password[i] = allChars[num.Int64()]
	}
	return string(password)
}
