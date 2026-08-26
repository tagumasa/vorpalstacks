package secretsmanager

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/http"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
)

const (
	defaultPasswordLength = 32
	minPasswordLength     = 1
	maxPasswordLength     = 4096
)

var (
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits           = "0123456789"
	punctuation      = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	space            = " "
)

// GetRandomPasswordInput carries all fields needed for GetRandomPassword.
type GetRandomPasswordInput struct {
	PasswordLength          int
	ExcludeCharacters       string
	ExcludeNumbers          bool
	ExcludePunctuation      bool
	ExcludeUppercase        bool
	ExcludeLowercase        bool
	IncludeSpace            bool
	RequireEachIncludedType bool
}

// GetRandomPasswordResult holds the transport-agnostic result of
// GetRandomPassword.
type GetRandomPasswordResult struct {
	RandomPassword string
}

// getRandomPasswordCore is the single entry point for random password
// generation: it validates the request and generates a password with the
// requested character-set constraints.
func (s *SecretsManagerService) getRandomPasswordCore(ctx context.Context, in GetRandomPasswordInput) (*GetRandomPasswordResult, error) {
	passwordLength := in.PasswordLength
	if passwordLength == 0 {
		passwordLength = defaultPasswordLength
	}
	if passwordLength < minPasswordLength || passwordLength > maxPasswordLength {
		return nil, awserrors.NewAWSError("InvalidParameterException", "PasswordLength must be between 1 and 4096", http.StatusBadRequest)
	}

	if err := validateExcludeCharacters(in.ExcludeCharacters); err != nil {
		return nil, err
	}

	charset := buildCharset(in.ExcludeCharacters, in.ExcludeNumbers, in.ExcludePunctuation, in.ExcludeUppercase, in.ExcludeLowercase, in.IncludeSpace)
	if len(charset) == 0 {
		return nil, awserrors.NewAWSError("InvalidParameterException", "No characters available for password generation", http.StatusBadRequest)
	}

	var password string
	var err error

	if in.RequireEachIncludedType {
		password, err = generatePasswordWithRequiredTypes(charset, passwordLength, in.ExcludeCharacters, in.ExcludeNumbers, in.ExcludePunctuation, in.ExcludeUppercase, in.ExcludeLowercase)
	} else {
		password, err = generatePassword(charset, passwordLength)
	}

	if err != nil {
		return nil, err
	}

	return &GetRandomPasswordResult{RandomPassword: password}, nil
}

func buildCharset(excludeCharacters string, excludeNumbers, excludePunctuation, excludeUppercase, excludeLowercase, includeSpace bool) string {
	var charset string

	if !excludeLowercase {
		charset += lowercaseLetters
	}
	if !excludeUppercase {
		charset += uppercaseLetters
	}
	if !excludeNumbers {
		charset += digits
	}
	if !excludePunctuation {
		charset += punctuation
	}
	if includeSpace {
		charset += space
	}

	for _, c := range excludeCharacters {
		charset = removeChar(charset, byte(c))
	}

	return charset
}

func removeChar(s string, c byte) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] != c {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

func generatePassword(charset string, length int) (string, error) {
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = charset[n.Int64()]
	}

	return string(result), nil
}

func generatePasswordWithRequiredTypes(charset string, length int, excludeCharacters string, excludeNumbers, excludePunctuation, excludeUppercase, excludeLowercase bool) (string, error) {
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	var requiredChars []byte
	if !excludeLowercase {
		if ch, ok := getRandomCharFromSet(lowercaseLetters, excludeCharacters); ok {
			requiredChars = append(requiredChars, ch)
		}
	}
	if !excludeUppercase {
		if ch, ok := getRandomCharFromSet(uppercaseLetters, excludeCharacters); ok {
			requiredChars = append(requiredChars, ch)
		}
	}
	if !excludeNumbers {
		if ch, ok := getRandomCharFromSet(digits, excludeCharacters); ok {
			requiredChars = append(requiredChars, ch)
		}
	}
	if !excludePunctuation {
		if ch, ok := getRandomCharFromSet(punctuation, excludeCharacters); ok {
			requiredChars = append(requiredChars, ch)
		}
	}

	if err := validatePasswordRequirements(len(requiredChars), length); err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = charset[n.Int64()]
	}

	for i := len(requiredChars) - 1; i > 0; i-- {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		j := int(n.Int64())
		requiredChars[i], requiredChars[j] = requiredChars[j], requiredChars[i]
	}

	positions := make([]int, length)
	for i := range positions {
		positions[i] = i
	}
	for i := length - 1; i > 0; i-- {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		j := int(n.Int64())
		positions[i], positions[j] = positions[j], positions[i]
	}

	for i, ch := range requiredChars {
		result[positions[i]] = ch
	}

	return string(result), nil
}

func getRandomCharFromSet(charset string, excludeCharacters string) (byte, bool) {
	availableChars := charset
	for _, c := range excludeCharacters {
		availableChars = removeChar(availableChars, byte(c))
	}
	if len(availableChars) == 0 {
		return 0, false
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(availableChars))))
	if err != nil {
		return availableChars[0], true
	}
	return availableChars[n.Int64()], true
}
