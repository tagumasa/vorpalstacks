package cognitoidentityprovider

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// This file implements verification of CSV-imported password hashes.
// Formats and parameter ceilings follow the Amazon Cognito developer guide
// ("Importing users with password hashes"): all four algorithms are
// self-describing, so every parameter is extracted from the hash string
// itself. Constants live in the store package as the single definitions.

var (
	bcryptHashPattern = regexp.MustCompile(`^\$2[abxy]\$(\d+)\$[./A-Za-z0-9]{53}$`)
	scryptHashPattern = regexp.MustCompile(`^(\d+)\$(\d+)\$(\d+)\$([0-9a-fA-F]+)\$([0-9a-fA-F]+)$`)
	argon2HashPattern = regexp.MustCompile(`^\$argon2id\$v=(\d+)\$m=(\d+),t=(\d+),p=(\d+)\$([A-Za-z0-9+/=]+)\$([A-Za-z0-9+/=]+)$`)
	pbkdf2HashPattern = regexp.MustCompile(`^\$pbkdf2-sha256\$(\d+)\$([A-Za-z0-9+/=]+)\$([A-Za-z0-9+/=]+)$`)
)

// validateImportedHashParams reports whether the encoded hash matches its
// algorithm's expected format and stays within the documented parameter
// ceilings. The import worker rejects rows whose hashes violate these
// bounds, mirroring AWS: "If a password hash has parameter values that
// exceed the maximum bounds listed above, the import fails for that user."
func validateImportedHashParams(algo, encoded string) error {
	switch algo {
	case "BCRYPT":
		m := bcryptHashPattern.FindStringSubmatch(encoded)
		if m == nil {
			return fmt.Errorf("malformed BCRYPT hash")
		}
		cost, _ := strconv.Atoi(m[1])
		if cost < cognitostore.MinImportHashParamValue {
			return fmt.Errorf("BCRYPT cost %d is below the minimum of %d", cost, cognitostore.MinImportHashParamValue)
		}
		if cost > cognitostore.MaxImportHashBcryptCost {
			return fmt.Errorf("BCRYPT cost %d exceeds the maximum of %d", cost, cognitostore.MaxImportHashBcryptCost)
		}
		return nil
	case "SCRYPT":
		m := scryptHashPattern.FindStringSubmatch(encoded)
		if m == nil {
			return fmt.Errorf("malformed SCRYPT hash")
		}
		n, _ := strconv.Atoi(m[1])
		r, _ := strconv.Atoi(m[2])
		p, _ := strconv.Atoi(m[3])
		if n < cognitostore.MinImportHashParamValue || r < cognitostore.MinImportHashParamValue || p < cognitostore.MinImportHashParamValue {
			return fmt.Errorf("SCRYPT parameters must all be at least %d", cognitostore.MinImportHashParamValue)
		}
		if n > cognitostore.MaxImportHashScryptN {
			return fmt.Errorf("SCRYPT N %d exceeds the maximum of %d", n, cognitostore.MaxImportHashScryptN)
		}
		if r > cognitostore.MaxImportHashScryptR {
			return fmt.Errorf("SCRYPT r %d exceeds the maximum of %d", r, cognitostore.MaxImportHashScryptR)
		}
		if p > cognitostore.MaxImportHashScryptP {
			return fmt.Errorf("SCRYPT p %d exceeds the maximum of %d", p, cognitostore.MaxImportHashScryptP)
		}
		return nil
	case "ARGON2ID":
		m := argon2HashPattern.FindStringSubmatch(encoded)
		if m == nil {
			return fmt.Errorf("malformed ARGON2ID hash")
		}
		mem, _ := strconv.Atoi(m[2])
		tCost, _ := strconv.Atoi(m[3])
		par, _ := strconv.Atoi(m[4])
		// RFC 9106 requires t >= 1, p >= 1, and m >= 8*p; a zero-valued
		// parameter is not a legitimate hash, and the argon2 library
		// panics on rounds or parallelism below one.
		if tCost < cognitostore.MinImportHashParamValue || par < cognitostore.MinImportHashParamValue {
			return fmt.Errorf("ARGON2ID t and p must be at least %d", cognitostore.MinImportHashParamValue)
		}
		if mem < 8*par {
			return fmt.Errorf("ARGON2ID m %d KiB is below the RFC 9106 minimum of 8*p (%d KiB)", mem, 8*par)
		}
		if mem > cognitostore.MaxImportHashArgon2MemKiB {
			return fmt.Errorf("ARGON2ID m %d KiB exceeds the maximum of %d KiB", mem, cognitostore.MaxImportHashArgon2MemKiB)
		}
		if tCost > cognitostore.MaxImportHashArgon2Time {
			return fmt.Errorf("ARGON2ID t %d exceeds the maximum of %d", tCost, cognitostore.MaxImportHashArgon2Time)
		}
		if par > cognitostore.MaxImportHashArgon2Parallelism {
			return fmt.Errorf("ARGON2ID p %d exceeds the maximum of %d", par, cognitostore.MaxImportHashArgon2Parallelism)
		}
		return nil
	case "PBKDF2_SHA256":
		m := pbkdf2HashPattern.FindStringSubmatch(encoded)
		if m == nil {
			return fmt.Errorf("malformed PBKDF2_SHA256 hash")
		}
		iter, _ := strconv.Atoi(m[1])
		if iter < cognitostore.MinImportHashParamValue {
			return fmt.Errorf("PBKDF2_SHA256 iterations %d are below the minimum of %d", iter, cognitostore.MinImportHashParamValue)
		}
		if iter > cognitostore.MaxImportHashPbkdf2Iterations {
			return fmt.Errorf("PBKDF2_SHA256 iterations %d exceed the maximum of %d", iter, cognitostore.MaxImportHashPbkdf2Iterations)
		}
		return nil
	default:
		return fmt.Errorf("unknown password hashing algorithm %q", algo)
	}
}

// verifyImportedPasswordHash reports whether the password verifies against
// the encoded hash produced by the named algorithm. Callers invoke it for
// users whose stored PasswordHashAlgo is set (CSV import with password
// hashes); on success the credentials migrate to the native bcrypt+SRP pair.
func verifyImportedPasswordHash(algo, encoded, password string) bool {
	if err := validateImportedHashParams(algo, encoded); err != nil {
		return false
	}
	switch algo {
	case "BCRYPT":
		return bcrypt.CompareHashAndPassword([]byte(encoded), []byte(password)) == nil
	case "SCRYPT":
		m := scryptHashPattern.FindStringSubmatch(encoded)
		n, _ := strconv.Atoi(m[1])
		r, _ := strconv.Atoi(m[2])
		p, _ := strconv.Atoi(m[3])
		salt, err := hex.DecodeString(m[4])
		if err != nil {
			return false
		}
		want, err := hex.DecodeString(m[5])
		if err != nil {
			return false
		}
		got, err := scrypt.Key([]byte(password), salt, n, r, p, len(want))
		if err != nil {
			return false
		}
		return subtle.ConstantTimeCompare(got, want) == 1
	case "ARGON2ID":
		m := argon2HashPattern.FindStringSubmatch(encoded)
		mem, _ := strconv.Atoi(m[2])
		tCost, _ := strconv.Atoi(m[3])
		par, _ := strconv.Atoi(m[4])
		salt, err := base64.RawStdEncoding.DecodeString(m[5])
		if err != nil {
			return false
		}
		want, err := base64.RawStdEncoding.DecodeString(m[6])
		if err != nil {
			return false
		}
		got := argon2.IDKey([]byte(password), salt, uint32(tCost), uint32(mem), uint8(par), uint32(len(want)))
		return subtle.ConstantTimeCompare(got, want) == 1
	case "PBKDF2_SHA256":
		m := pbkdf2HashPattern.FindStringSubmatch(encoded)
		iter, _ := strconv.Atoi(m[1])
		salt, err := base64.RawStdEncoding.DecodeString(m[2])
		if err != nil {
			return false
		}
		want, err := base64.RawStdEncoding.DecodeString(m[3])
		if err != nil {
			return false
		}
		got := pbkdf2.Key([]byte(password), salt, iter, len(want), sha256.New)
		return subtle.ConstantTimeCompare(got, want) == 1
	default:
		return false
	}
}
