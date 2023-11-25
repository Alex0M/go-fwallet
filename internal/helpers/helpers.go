package helpers

import (
	"fmt"
	"log"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		log.Printf("Error while hashing password, Reason: %v\n", err)
	}

	return string(hashed)
}

func StringToInt(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return -1, fmt.Errorf("cannot convert string %s to int", s)
	}

	return n, nil
}
