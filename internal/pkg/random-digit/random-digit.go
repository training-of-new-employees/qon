package randomdigit

import (
	"math/rand"
	"strconv"
	"time"
)

// RandomDigitNumber генерирует случайное n-значное число в виде строки.
func RandomDigitNumber(n int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// генерация случайного n-значного числа
	number := ""
	for i := 0; i < n; i++ {
		number += strconv.Itoa(rnd.Intn(10))
	}

	return number
}
