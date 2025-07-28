package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordConfig holds configuration for password hashing
type PasswordConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// DefaultPasswordConfig returns a secure default configuration for password hashing
func DefaultPasswordConfig() *PasswordConfig {
	return &PasswordConfig{
		Time:    1,
		Memory:  64 * 1024, // 64 MB
		Threads: 4,
		KeyLen:  32,
	}
}

var (
	ErrInvalidHash         = errors.New("invalid hash format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

// HashPassword creates a hash of the password using Argon2id
func HashPassword(password string, config *PasswordConfig) (string, error) {
	if config == nil {
		config = DefaultPasswordConfig()
	}

	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Generate the hash
	hash := argon2.IDKey([]byte(password), salt, config.Time, config.Memory, config.Threads, config.KeyLen)

	// Encode the salt and hash
	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)
	hashEncoded := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, config.Memory, config.Time, config.Threads, saltEncoded, hashEncoded), nil
}

// VerifyPassword verifies a password against its hash
func VerifyPassword(password, hash string) (bool, error) {
	// Parse the hash
	config, salt, hashBytes, err := parseHash(hash)
	if err != nil {
		return false, err
	}

	// Generate hash from provided password
	providedHash := argon2.IDKey([]byte(password), salt, config.Time, config.Memory, config.Threads, config.KeyLen)

	// Compare hashes using constant time comparison
	return subtle.ConstantTimeCompare(hashBytes, providedHash) == 1, nil
}

// parseHash extracts the configuration and salt from an encoded hash
func parseHash(hash string) (*PasswordConfig, []byte, []byte, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	config := &PasswordConfig{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &config.Memory, &config.Time, &config.Threads); err != nil {
		return nil, nil, nil, err
	}
	config.KeyLen = 32

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}

	hashBytes, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}

	return config, salt, hashBytes, nil
}

// ValidatePasswordStrength validates password strength
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	var missing []string
	if !hasUpper {
		missing = append(missing, "uppercase letter")
	}
	if !hasLower {
		missing = append(missing, "lowercase letter")
	}
	if !hasDigit {
		missing = append(missing, "digit")
	}
	if !hasSpecial {
		missing = append(missing, "special character")
	}

	if len(missing) > 0 {
		return fmt.Errorf("password must contain at least one: %s", strings.Join(missing, ", "))
	}

	return nil
}
