package cognitoidentityprovider

import (
	"sync"
	"testing"

	"golang.org/x/crypto/bcrypt"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// Two concurrent password writes through the UpdateUserFunc seam must both
// land in the password history: the serialised read-modify-write sees the
// other writer's hash instead of overwriting it with a stale snapshot. With
// the detached read-modify-write this test pins against, the later write
// recorded only the seed hash and one racer's password vanished.
func TestConcurrentPasswordWritesKeepHistory(t *testing.T) {
	_, _, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("password-race", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	var policy *cognitostore.PasswordPolicy
	if err := store.UpdateUserPoolFunc(pool.ID, func(p *cognitostore.UserPool) error {
		policy = &cognitostore.PasswordPolicy{PasswordHistorySize: 2}
		p.PasswordPolicy = policy
		return nil
	}); err != nil {
		t.Fatalf("set password policy: %v", err)
	}

	user := cognitostore.NewUser(pool.ID, "racer")
	user.Attributes = map[string]string{}
	if err := store.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := store.UpdateUserFunc(pool.ID, "racer", func(u *cognitostore.User) error {
		return setNativePasswordCredentials(u, policy, "Initial123!")
	}); err != nil {
		t.Fatalf("seed native credentials: %v", err)
	}

	const pwOne = "RacerOne123!"
	const pwTwo = "RacerTwo123!"
	var wg sync.WaitGroup
	start := make(chan struct{})
	for _, pw := range []string{pwOne, pwTwo} {
		wg.Add(1)
		go func(pw string) {
			defer wg.Done()
			<-start
			_ = store.UpdateUserFunc(pool.ID, "racer", func(u *cognitostore.User) error {
				return setNativePasswordCredentials(u, policy, pw)
			})
		}(pw)
	}
	close(start)
	wg.Wait()

	got, err := store.GetUser(pool.ID, "racer")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	oneCurrent := bcrypt.CompareHashAndPassword([]byte(got.PasswordHash), []byte(pwOne)) == nil
	twoCurrent := bcrypt.CompareHashAndPassword([]byte(got.PasswordHash), []byte(pwTwo)) == nil
	if oneCurrent == twoCurrent {
		t.Fatalf("current hash matches neither or both racers (one=%v two=%v)", oneCurrent, twoCurrent)
	}
	superseded := pwOne
	if oneCurrent {
		superseded = pwTwo
	}
	if len(got.PasswordHistory) == 0 {
		t.Fatal("password history empty — one racer's write was lost")
	}
	// The most recent superseded hash is the other racer's password; the
	// seed hash follows it.
	if bcrypt.CompareHashAndPassword([]byte(got.PasswordHistory[0]), []byte(superseded)) != nil {
		t.Fatalf("superseded racer password missing from history (%d entries)", len(got.PasswordHistory))
	}
}
