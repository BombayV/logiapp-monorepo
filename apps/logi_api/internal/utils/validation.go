package utils

import (
	"strconv"
)

// ValidateEcuadorianCedula checks if the provided string is a valid Ecuadorian cedula.
func ValidateEcuadorianCedula(cedula string) bool {
	if len(cedula) != 10 {
		return false
	}

	// Check if it contains only digits
	for _, char := range cedula {
		if char < '0' || char > '9' {
			return false
		}
	}

	// First two digits (province code) must be between 01 and 24
	provinceCode, _ := strconv.Atoi(cedula[0:2])
	if provinceCode < 1 || provinceCode > 24 {
		return false
	}

	// Third digit must be less than 6
	thirdDigit, _ := strconv.Atoi(string(cedula[2]))
	if thirdDigit >= 6 {
		return false
	}

	// Calculate check digit
	coefficients := []int{2, 1, 2, 1, 2, 1, 2, 1, 2}
	total := 0

	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(cedula[i]))
		product := digit * coefficients[i]
		if product >= 10 {
			product -= 9
		}
		total += product
	}

	checkDigit := 10 - (total % 10)
	if checkDigit == 10 {
		checkDigit = 0
	}

	lastDigit, _ := strconv.Atoi(string(cedula[9]))

	return checkDigit == lastDigit
}
