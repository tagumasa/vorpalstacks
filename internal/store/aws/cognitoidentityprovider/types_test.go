package cognitoidentityprovider

import "testing"

// The ID token's email_verified claim derives from the stored attribute:
// the string "true" is verified, anything else (including absent) is not.
func TestUserGetEmailVerified(t *testing.T) {
	for name, tc := range map[string]struct {
		attrs map[string]string
		want  bool
	}{
		"verified":       {map[string]string{"email_verified": "true"}, true},
		"explicit false": {map[string]string{"email_verified": "false"}, false},
		"no attribute":   {map[string]string{"email": "a@b.c"}, false},
		"nil attributes": {nil, false},
	} {
		u := &User{Attributes: tc.attrs}
		if got := u.GetEmailVerified(); got != tc.want {
			t.Errorf("%s: got %v, want %v", name, got, tc.want)
		}
	}
}
