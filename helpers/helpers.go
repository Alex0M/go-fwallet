package helpers

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		log.Printf("Error while hashing password, Reason: %v\n", err)
	}

	return string(hashed)
}
