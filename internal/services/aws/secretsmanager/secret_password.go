package secretsmanager

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/http"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
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

// GetRandomPassword generates a random password with configurable options.
// It supports custom length, character sets, and requiring at least one character
// from each included character type.
func (s *SecretsManagerService) GetRandomPassword(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	passwordLength := request.GetIntParam(req.Parameters, "PasswordLength")
	if passwordLength == 0 {
		passwordLength = defaultPasswordLength
	}
	if passwordLength < minPasswordLength || passwordLength > maxPasswordLength {
		return nil, awserrors.NewAWSError("InvalidParameterException", "PasswordLength must be between 1 and 4096", http.StatusBadRequest)
	}

	excludeCharacters := request.GetStringParam(req.Parameters, "ExcludeCharacters")
	excludeNumbers := request.GetBoolParam(req.Parameters, "ExcludeNumbers")
	excludePunctuation := request.GetBoolParam(req.Parameters, "ExcludePunctuation")
	excludeUppercase := request.GetBoolParam(req.Parameters, "ExcludeUppercase")
	excludeLowercase := request.GetBoolParam(req.Parameters, "ExcludeLowercase")
	includeSpace := request.GetBoolParam(req.Parameters, "IncludeSpace")
	requireEachIncludedType := request.GetBoolParam(req.Parameters, "RequireEachIncludedType")

	charset := buildCharset(excludeCharacters, excludeNumbers, excludePunctuation, excludeUppercase, excludeLowercase, includeSpace)
	if len(charset) == 0 {
		return nil, awserrors.NewAWSError("InvalidParameterException", "No characters available for password generation", http.StatusBadRequest)
	}

	var password string
	var err error

	if requireEachIncludedType {
		password, err = generatePasswordWithRequiredTypes(charset, passwordLength, excludeCharacters, excludeNumbers, excludePunctuation, excludeUppercase, excludeLowercase)
	} else {
		password, err = generatePassword(charset, passwordLength)
	}

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"RandomPassword": password,
	}, nil
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
		if i < length {
			result[positions[i]] = ch
		}
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
