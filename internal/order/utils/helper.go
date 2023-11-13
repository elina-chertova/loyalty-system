package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateOrderNumber() string {
	rand.Seed(time.Now().UnixNano())

	firstDigit := rand.Intn(9) + 1
	id := strconv.Itoa(firstDigit)

	for i := 0; i < 14; i++ {
		id += strconv.Itoa(rand.Intn(10))
	}

	return id
}

func IsLuhnValid(orderNumber string) bool {
	var digits []int
	for _, char := range orderNumber {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return false
		}
		digits = append(digits, digit)
	}

	for i := len(digits) - 2; i >= 0; i -= 2 {
		digits[i] *= 2
		if digits[i] > 9 {
			digits[i] -= 9
		}
	}

	sum := 0
	for _, digit := range digits {
		sum += digit
	}
	return sum%10 == 0
}
