package randomseq

import (
	"crypto/rand"
	"encoding/hex"
	mrand "math/rand"
	"strconv"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandomDigitNumber генерирует случайное n-значное число в виде строки.
func RandomDigitNumber(n int) string {
	rnd := mrand.New(mrand.NewSource(time.Now().UnixNano()))

	// генерация случайного n-значного числа
	number := ""
	for i := 0; i < n; i++ {
		number += strconv.Itoa(rnd.Intn(10))
	}

	return number
}

// RandomHexString генерирует случайную последовательность байт заданной длины и отдаёт их в качестве строки.
func RandomHexString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
