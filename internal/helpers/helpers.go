package helpers

import (
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pass []byte) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func StringToInt(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return -1, err
	}

	return n, nil
}
