package inspection

import (
	"strconv"
	"testing"
	"time"
)

func TestTokenSignsAndParses(t *testing.T) {
	secret := []byte("unit-test-secret")
	issued := time.Unix(1700000000, 0)
	token := ChallengeToken{
		IssuedAt:          issued.Unix(),
		ChallengeSolvedAt: issued.Unix(),
		Domains:           []string{"example.com"},
	}
	value := SignToken(secret, token)
	if value == "" {
		t.Fatal("SignToken returned an empty value")
	}
	parsed, ok := ParseToken(secret, value)
	if !ok {
		t.Fatal("ParseToken rejected a freshly signed token")
	}
	if parsed.ChallengeSolvedAt != issued.Unix() || parsed.CaptchaSolvedAt != 0 {
		t.Fatalf("payload round-trip mismatch: %+v", parsed)
	}
	if len(parsed.Domains) != 1 || parsed.Domains[0] != "example.com" {
		t.Fatalf("domains round-trip mismatch: %+v", parsed.Domains)
	}
}

func TestTokenRejectsBadSignatureAndWrongKey(t *testing.T) {
	secret := []byte("unit-test-secret")
	value := SignToken(secret, ChallengeToken{IssuedAt: 1})
	if _, ok := ParseToken([]byte("other-secret"), value); ok {
		t.Fatal("ParseToken accepted a token signed with a different key")
	}
	tampered := value[:len(value)-2] + "AA"
	if _, ok := ParseToken(secret, tampered); ok {
		t.Fatal("ParseToken accepted a tampered token")
	}
	if _, ok := ParseToken(secret, "not-a-token"); ok {
		t.Fatal("ParseToken accepted a malformed value")
	}
}

func TestTokenCoversHost(t *testing.T) {
	token := ChallengeToken{Domains: []string{"Example.COM"}}
	cases := []struct {
		host string
		want bool
	}{
		{"example.com", true},
		{"www.example.com", true},
		{"deep.www.example.com", true},
		{"notexample.com", false},
		{"example.com.evil.net", false},
		{"", false},
	}
	for _, tc := range cases {
		if got := token.CoversHost(tc.host); got != tc.want {
			t.Errorf("CoversHost(%q) = %v, want %v", tc.host, got, tc.want)
		}
	}
}

func TestTokenSolvedWithin(t *testing.T) {
	now := time.Unix(1700000300, 0)
	token := ChallengeToken{
		ChallengeSolvedAt: now.Unix() - 100,
		CaptchaSolvedAt:   now.Unix() - 400,
	}
	immunity := 300 * time.Second
	if !token.SolvedWithin(ActionChallenge, now, immunity) {
		t.Error("challenge solve within immunity should be valid")
	}
	if token.SolvedWithin(ActionCaptcha, now, immunity) {
		t.Error("captcha solve past immunity should be invalid")
	}
	empty := ChallengeToken{}
	if empty.SolvedWithin(ActionCaptcha, now, immunity) {
		t.Error("an unsolved kind should never be valid")
	}
	if token.SolvedWithin("Monetize", now, immunity) {
		t.Error("an unknown kind should never be valid")
	}
}

func TestVerifyChallengeSolution(t *testing.T) {
	challengeID := "abc123"
	var counter string
	for i := 0; ; i++ {
		candidate := strconv.Itoa(i)
		if VerifyChallengeSolution(challengeID, candidate) {
			counter = candidate
			break
		}
	}
	if !VerifyChallengeSolution(challengeID, counter) {
		t.Fatalf("re-verification failed for the found counter %q", counter)
	}
	if VerifyChallengeSolution(challengeID, "1") {
		t.Fatal("counter 1 must not solve a 4-hex-zero challenge except by astronomical luck")
	}
	if VerifyChallengeSolution("", "1") || VerifyChallengeSolution("x", "") {
		t.Fatal("empty inputs must not verify")
	}
	if VerifyChallengeSolution("x", "-3") {
		t.Fatal("non-numeric counters must not verify")
	}
}
