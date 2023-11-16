package utils

import (
	"strconv"
)

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
