package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// HashService provides password hashing and comparison functionalities.
type argonParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// HashPassword creates a salted Argon2id hash of a password.
func HashPassword(password string) (string, error) {
	// Use the recommended Argon2id variant.
	// These parameters are a good starting point.
	p := &argonParams{
		memory:      64 * 1024, // 64 MB
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	// Generate a cryptographically secure random salt.
	salt := make([]byte, p.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Generate the hash.
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Encode the salt and hash into a single string for storage.
	// Format: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// CheckPasswordHash compares a password with a stored Argon2id hash.
func CheckPasswordHash(password, encodedHash string) (bool, error) {
	// Decode the hash string.
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Generate a hash of the provided password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Compare the two hashes in constant time to prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

// decodeHash parses the encoded hash string and returns its components.
func decodeHash(encodedHash string) (p *argonParams, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("invalid encoded hash format")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible argon2 version")
	}

	p = &argonParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
