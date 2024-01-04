package randomseq

import (
	"crypto/rand"
	"encoding/hex"
	mrand "math/rand"
	"strconv"
	"time"
)

var rnd *mrand.Rand
var alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rnd = mrand.New(mrand.NewSource(time.Now().UnixNano()))
}

// RandomDigitNumber генерирует случайное n-значное число в виде строки.
func RandomDigitNumber(n int) string {

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
	n, err := rand.Read(b)
	if err != nil {
		return RandomDigitNumber(n)
	}
	return hex.EncodeToString(b)
}

// RandomString генерирует случайную строку заданной длины
func RandomString(n int) string {
	seq := make([]byte, 0, n)
	for len(seq) < n {
		symb := alphabet[rnd.Intn(len(alphabet))]
		seq = append(seq, byte(symb))

	}
	return string(seq)
}

// RandomTestInt - генерирует случайное число от 100 до 356.
func RandomTestInt() int {
	return 100 + rnd.Intn(256)
}
