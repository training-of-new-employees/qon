// Package randomseq - пакет для генерации случайных данных.
package randomseq

import (
	"crypto/rand"
	"encoding/hex"
	mrand "math/rand"
	"strconv"
	"time"
)

var rnd *mrand.Rand

// alphabetRandomString - набор для генерации случайной строки
var alphabetRandomString = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// alphabetName - набор для генерации имён/названий
var alphabetName = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	// lowerLetters - набор символов английского алфавита в нижнем регистре
	lowerLetters = "abcdefghijklmnopqrstuvwxyz"
	// upperLetters - набор символов английского алфавита в верхнем регистре
	upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// numbers - набор символов чисел
	numbers = "0123456789"
	// symbols - набор случайных символов, не являющиеся цифрами или буквами
	symbols = "!@#$%^&*()_+"
)

func init() {
	rnd = mrand.New(mrand.NewSource(time.Now().UnixNano()))
}

// RandomDigitNumber - генерирует случайное n-значное число в виде строки.
func RandomDigitNumber(n int) string {

	// генерация случайного n-значного числа
	number := ""
	for i := 0; i < n; i++ {
		number += strconv.Itoa(rnd.Intn(10))
	}

	return number
}

// RandomHexString - генерирует случайную последовательность байт заданной длины и отдаёт их в качестве строки.
func RandomHexString(n int) string {
	b := make([]byte, n)
	n, err := rand.Read(b)
	if err != nil {
		return RandomDigitNumber(n)
	}
	return hex.EncodeToString(b)
}

// RandomString - генерирует случайную строку заданной длины.
func RandomString(n int) string {
	seq := make([]byte, 0, n)
	for len(seq) < n {
		symb := alphabetRandomString[rnd.Intn(len(alphabetRandomString))]
		seq = append(seq, byte(symb))

	}
	return string(seq)
}

// RandomName - создает последовательность, удовлетворяющую требованиям валидатора имен.
func RandomName(n int) string {
	seq := make([]byte, 0, n)
	for len(seq) < n {
		symb := alphabetName[rnd.Intn(len(alphabetName))]
		seq = append(seq, byte(symb))

	}
	return string(seq)
}

// RandomBool - генерирует случайное булево значение.
func RandomBool() bool {
	return rnd.Intn(2) == 1
}

// RandomPassword - генерирует случайный пароль.
func RandomPassword() string {
	passwordLen := 30
	upperLettersCount := passwordLen / 2
	numbersCount := upperLettersCount / 2
	symbolsCount := numbersCount / 2
	password := make([]byte, passwordLen)

	for i := 0; i < len(password); i++ {
		password[i] = lowerLetters[rnd.Intn(len(lowerLetters))]
	}

	for i := 0; i < upperLettersCount; i++ {
		randomPosition := rnd.Intn(passwordLen)
		password[randomPosition] = upperLetters[rnd.Intn(len(upperLetters))]
	}

	for i := 0; i < numbersCount; i++ {
		randomPosition := rnd.Intn(passwordLen)
		password[randomPosition] = numbers[rnd.Intn(len(numbers))]
	}

	for i := 0; i < symbolsCount; i++ {
		randomPosition := rnd.Intn(passwordLen)
		password[randomPosition] = symbols[rnd.Intn(len(symbols))]
	}

	return string(password)
}

// RandomTestInt - генерирует случайное число от 100 до 356.
// Используется в тестах в основном для генерации идентификаторов объектов (также может использоваться и для других целей).
// Генерируются значения в диапазоне от 100 до 356, которые не включаются в тестируемые данные.
// При тестировании преимущественно используются первые цифры.
func RandomTestInt() int {
	return 100 + rnd.Intn(256)
}
