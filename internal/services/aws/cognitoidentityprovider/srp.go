package cognitoidentityprovider

// This file implements the AWS Cognito variant of the Secure Remote Password
// (SRP) protocol. The protocol differs from standard SRP-6a in two ways:
//
//  1. The shared secret S is fed through HKDF (info="Caldera Derived Key",
//     salt=u) rather than a plain hash, producing a 16-byte symmetric key K.
//  2. The client proof is HMAC-SHA256(K, poolName||userId||secretBlock||
//     timestamp) rather than the standard SRP-6a M1. The inclusion of the
//     per-challenge SECRET_BLOCK and TIMESTAMP binds the proof to a single
//     challenge, preventing replay.
//
// Reference implementation:
//   - amazon-cognito-identity-js/src/AuthenticationHelper.js (canonical)
//   - github.com/alexrudd/cognito-srp (de facto Go port)
//
// The 2048-bit safe prime and generator below are copied verbatim from the
// amazon-cognito-identity-js source; they match RFC 5054 Appendix A.

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
)

// cognitoSrpNHex is the hex-encoded 2048-bit safe prime N from RFC 5054.
const cognitoSrpNHex = "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
	"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
	"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
	"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
	"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
	"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
	"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
	"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
	"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
	"DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
	"15728E5A8AAAC42DAD33170D04507A33A85521ABDF1CBA64" +
	"ECFB850458DBEF0A8AEA71575D060C7DB3970F85A6E1E4C7" +
	"ABF5AE8CDB0933D71E8C94E04A25619DCEE3D2261AD2EE6B" +
	"F12FFA06D98A0864D87602733EC86A64521F2B18177B200C" +
	"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB31" +
	"43DB5BFCE0FD108E4B82D120A93AD2CAFFFFFFFFFFFFFFFF"

const (
	cognitoSrpGHex     = "2"
	cognitoSrpInfoBits = "Caldera Derived Key"
)

var (
	cognitoSrpN = mustHexToBig(cognitoSrpNHex)
	cognitoSrpG = mustHexToBig(cognitoSrpGHex)
	// cognitoSrpK is the SRP multiplier parameter k = H(N || pad(g)). It is
	// computed once at package initialisation.
	cognitoSrpK = mustHexToBig(hexHash("00" + cognitoSrpNHex + "0" + cognitoSrpGHex))
)

// ErrInvalidSrpA is returned when the client-supplied SRP_A value is invalid
// (specifically, when A mod N == 0, which would allow trivial attacks).
var ErrInvalidSrpA = errors.New("cognito srp: invalid SRP_A (A mod N == 0)")

// ComputeVerifier derives the SRP verifier v = g^x mod N, where
// x = H(padHex(saltHex) || H(poolName || userId || ":" || password)).
//
// The inner hash deliberately includes poolName: this is Cognito's SRP variant
// and is required for interoperability with AWS SDK clients. saltHex must be
// the hex-encoded random salt (typically 16 random bytes => 32 hex chars).
// poolName is the portion of the user pool ID after the underscore (e.g.
// "us-east-1_abc123" => "abc123").
//
// The salt is normalised by routing it through big.Int before hashing, mirroring
// the reference client (amazon-cognito-identity-js / alexrudd/cognito-srp).
// This strips any leading zero bytes that hex.EncodeToString would otherwise
// preserve, ensuring the server and client compute identical x values even
// when the random salt begins with one or more 0x00 bytes.
func ComputeVerifier(saltHex, poolName, userID, password string) *big.Int {
	saltInt, ok := new(big.Int).SetString(saltHex, 16)
	if !ok {
		saltInt = big.NewInt(0)
	}
	normalisedSalt := padHex(saltInt.Text(16))
	userPassHash := hashSHA256Hex([]byte(poolName + userID + ":" + password))
	x := mustHexToBig(hexHash(normalisedSalt + userPassHash))
	return new(big.Int).Exp(cognitoSrpG, x, cognitoSrpN)
}

// GenerateB produces the server's ephemeral public value B and the matching
// private scalar b. b is a fresh 256-bit random integer reduced mod N, and
// B is computed as (k*v + g^b) mod N.
func GenerateB(verifier *big.Int) (B, b *big.Int, err error) {
	for attempt := 0; attempt < 32; attempt++ {
		b, err = randBigIntModN()
		if err != nil {
			return nil, nil, err
		}
		gb := new(big.Int).Exp(cognitoSrpG, b, cognitoSrpN)
		kv := new(big.Int).Mul(cognitoSrpK, verifier)
		kv.Mod(kv, cognitoSrpN)
		B = new(big.Int).Add(kv, gb)
		B.Mod(B, cognitoSrpN)
		if B.Sign() != 0 {
			return B, b, nil
		}
	}
	return nil, nil, errors.New("SRP: failed to generate B != 0 mod N after 32 attempts")
}

// ComputeServerS derives the shared secret S = (A * v^u)^b mod N. It rejects
// A values for which A mod N == 0 to prevent forgery.
func ComputeServerS(A, B, b, verifier *big.Int) (*big.Int, error) {
	if new(big.Int).Mod(A, cognitoSrpN).Sign() == 0 {
		return nil, ErrInvalidSrpA
	}
	u := calculateU(A, B)
	vu := new(big.Int).Exp(verifier, u, cognitoSrpN)
	aTimesVu := new(big.Int).Mul(A, vu)
	aTimesVu.Mod(aTimesVu, cognitoSrpN)
	return new(big.Int).Exp(aTimesVu, b, cognitoSrpN), nil
}

// DeriveServerKey is the high-level helper used by the PASSWORD_VERIFIER
// handler. It computes the shared secret S and derives the 16-byte HMAC key K
// in a single call, returning K ready for use with VerifyClaim.
func DeriveServerKey(A, B, b, verifier *big.Int) ([]byte, error) {
	S, err := ComputeServerS(A, B, b, verifier)
	if err != nil {
		return nil, err
	}
	u := calculateU(A, B)
	return DeriveKey(S, u), nil
}

// DeriveKey extracts the 16-byte symmetric key K from the shared secret S
// and the scrambling parameter u via HKDF. Cognito's HKDF uses u as the salt
// and the literal string "Caldera Derived Key" as the info field; the output
// is the first 16 bytes of the HMAC-SHA256 PRK expansion.
func DeriveKey(S, u *big.Int) []byte {
	return computeHKDF(padHex(S.Text(16)), padHex(u.Text(16)))
}

// VerifyClaim computes the expected HMAC-SHA256 claim over the Cognito message
//
//	msg = poolName || userId || secretBlock || timestamp
//
// using the derived key K. The returned slice is the digest that the client's
// PASSWORD_CLAIM_SIGNATURE must match (constant-time comparison is the
// caller's responsibility).
func VerifyClaim(K []byte, poolName, userID string, secretBlock []byte, timestamp string) []byte {
	mac := hmac.New(sha256.New, K)
	mac.Write([]byte(poolName))
	mac.Write([]byte(userID))
	mac.Write(secretBlock)
	mac.Write([]byte(timestamp))
	return mac.Sum(nil)
}

// calculateU computes the SRP scrambling parameter u = H(padHex(A) || padHex(B)).
// Both A and B are first serialised to lowercase hex and padded, then hashed.
func calculateU(A, B *big.Int) *big.Int {
	return mustHexToBig(hexHash(padHex(A.Text(16)) + padHex(B.Text(16))))
}

// computeHKDF performs Cognito's HKDF extraction-and-expansion. ikm and salt
// are hex strings; the function decodes them and returns the first 16 bytes
// of the expansion.
func computeHKDF(ikmHex, saltHex string) []byte {
	ikm, _ := hex.DecodeString(ikmHex)
	salt, _ := hex.DecodeString(saltHex)

	extractor := hmac.New(sha256.New, salt)
	extractor.Write(ikm)
	prk := extractor.Sum(nil)

	// infoBits is followed by a single 0x01 byte per the HKDF counter layout.
	info := append([]byte(cognitoSrpInfoBits), 1)
	extractor = hmac.New(sha256.New, prk)
	extractor.Write(info)
	return extractor.Sum(nil)[:16]
}

// hashSHA256Hex returns the lowercase hex-encoded SHA-256 digest of data.
func hashSHA256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// hexHash decodes a hex string, hashes the resulting bytes with SHA-256, and
// returns the hex-encoded digest. It is the canonical Cognito "hash of hex"
// operation used throughout the SRP computations.
func hexHash(hexStr string) string {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		// All callers pass constants or already-validated hex; reaching here
		// indicates a programming error.
		panic("cognito srp: invalid hex input to hexHash: " + err.Error())
	}
	return hashSHA256Hex(b)
}

// padHex ensures the hex string represents an unsigned big-endian integer.
// Odd-length strings are left-padded with "0"; even-length strings whose
// first nibble is 8..f are left-padded with "00" to keep the value unsigned.
func padHex(hexStr string) string {
	if len(hexStr)%2 == 1 {
		return "0" + hexStr
	}
	if len(hexStr) > 0 && strings.IndexByte("89abcdef", hexStr[0]) >= 0 {
		return "00" + hexStr
	}
	return hexStr
}

// randBigIntModN returns a fresh uniformly-distributed scalar in [0, N).
// 32 bytes (256 bits) of entropy exceed the 128-bit security level of the
// 2048-bit SRP group, in line with RFC 5054 recommendations.
func randBigIntModN() (*big.Int, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	r := new(big.Int).SetBytes(b)
	return r.Mod(r, cognitoSrpN), nil
}

// mustHexToBig parses a hex string into *big.Int. It panics on malformed input
// because all callers pass compile-time constants.
func mustHexToBig(hexStr string) *big.Int {
	n, ok := new(big.Int).SetString(hexStr, 16)
	if !ok {
		panic("cognito srp: invalid hex constant")
	}
	return n
}
