package core

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var defaultParams = params{
	memory:      32768, // 32 MB
	iterations:  4,
	parallelism: 1,
	saltLength:  16, // 16 bytes
	keyLength:   32, // 32 bytes
}

func generateRandomBytes(length uint32) ([]byte, error) {
	// Create a byte slice of the specified length.
	bytes := make([]byte, length)
	// Fill the byte slice with cryptographically secure random bytes.
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func generateFromPassword(password string, p *params) (hash []byte, err error) {
	// Generate a cryptographically secure random salt.
	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return nil, err
	}
	
	hash = argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	return hash, nil
}
