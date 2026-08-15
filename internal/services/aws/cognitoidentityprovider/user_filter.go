package cognitoidentityprovider

import (
	"regexp"
	"strings"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

var userFilterRe = regexp.MustCompile(`^"?(\w[\w:.\-+]*)\s*(=|\^=)\s*"?([^"]+?)"?\s*$`)

func matchUserFilter(user *cognitostore.User, filter string) bool {
	f := strings.TrimSpace(filter)
	m := userFilterRe.FindStringSubmatch(f)
	if m == nil {
		return false
	}
	attrName := m[1]
	op := m[2]
	attrValue := m[3]

	var actual string
	switch strings.ToLower(attrName) {
	case "username":
		actual = user.Username
	case "cognito:user_status":
		actual = user.UserStatus
	case "status":
		if user.Enabled {
			actual = "Enabled"
		} else {
			actual = "Disabled"
		}
	default:
		if user.Attributes != nil {
			actual = user.Attributes[attrName]
		}
	}

	switch op {
	case "=":
		return strings.EqualFold(actual, attrValue)
	case "^=":
		return strings.HasPrefix(strings.ToLower(actual), strings.ToLower(attrValue))
	}
	return false
}
