package s3

import (
	"net/http"
	"testing"

	"vorpalstacks/internal/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

func TestAccessController_IsOwner(t *testing.T) {
	ac := NewAccessController("account-123")

	tests := []struct {
		name      string
		check     *AccessCheck
		bucket    *s3store.Bucket
		wantOwner bool
	}{
		{
			name: "same account as owner",
			check: &AccessCheck{
				PrincipalID:   "account-123",
				PrincipalType: request.PrincipalTypeUser,
			},
			bucket: &s3store.Bucket{
				Name: "test-bucket",
				ACL: &s3store.AccessControlPolicy{
					Owner: &s3store.ACLOwner{ID: "account-123"},
				},
			},
			wantOwner: true,
		},
		{
			name: "different account",
			check: &AccessCheck{
				PrincipalID:   "account-456",
				PrincipalType: request.PrincipalTypeUser,
			},
			bucket: &s3store.Bucket{
				Name: "test-bucket",
				ACL: &s3store.AccessControlPolicy{
					Owner: &s3store.ACLOwner{ID: "account-123"},
				},
			},
			wantOwner: false,
		},
		{
			name: "no ACL defaults to account owner",
			check: &AccessCheck{
				PrincipalID:   "account-123",
				PrincipalType: request.PrincipalTypeUser,
			},
			bucket: &s3store.Bucket{
				Name: "test-bucket",
			},
			wantOwner: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ac.isOwner(tt.check, tt.bucket)
			if got != tt.wantOwner {
				t.Errorf("isOwner() = %v, want %v", got, tt.wantOwner)
			}
		})
	}
}

func TestAccessController_GrantMatchesPrincipal(t *testing.T) {
	ac := NewAccessController("account-123")

	tests := []struct {
		name  string
		grant *s3store.Grant
		check *AccessCheck
		want  bool
	}{
		{
			name: "canonical user match",
			grant: &s3store.Grant{
				Grantee: &s3store.Grantee{
					Type: s3store.GranteeTypeCanonicalUser,
					ID:   "user-123",
				},
			},
			check: &AccessCheck{
				PrincipalID:   "user-123",
				PrincipalType: request.PrincipalTypeUser,
			},
			want: true,
		},
		{
			name: "canonical user no match",
			grant: &s3store.Grant{
				Grantee: &s3store.Grantee{
					Type: s3store.GranteeTypeCanonicalUser,
					ID:   "user-123",
				},
			},
			check: &AccessCheck{
				PrincipalID:   "user-456",
				PrincipalType: request.PrincipalTypeUser,
			},
			want: false,
		},
		{
			name: "all users group matches anonymous",
			grant: &s3store.Grant{
				Grantee: &s3store.Grantee{
					Type: s3store.GranteeTypeGroup,
					URI:  s3store.AllUsersGroup,
				},
			},
			check: &AccessCheck{
				PrincipalID:   "",
				PrincipalType: request.PrincipalTypeAnonymous,
			},
			want: true,
		},
		{
			name: "authenticated users group matches user",
			grant: &s3store.Grant{
				Grantee: &s3store.Grantee{
					Type: s3store.GranteeTypeGroup,
					URI:  s3store.AuthenticatedUsersGroup,
				},
			},
			check: &AccessCheck{
				PrincipalID:   "user-123",
				PrincipalType: request.PrincipalTypeUser,
			},
			want: true,
		},
		{
			name: "authenticated users group does not match anonymous",
			grant: &s3store.Grant{
				Grantee: &s3store.Grantee{
					Type: s3store.GranteeTypeGroup,
					URI:  s3store.AuthenticatedUsersGroup,
				},
			},
			check: &AccessCheck{
				PrincipalID:   "",
				PrincipalType: request.PrincipalTypeAnonymous,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ac.grantMatchesPrincipal(tt.grant, tt.check)
			if got != tt.want {
				t.Errorf("grantMatchesPrincipal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccessController_PermissionMatchesAction(t *testing.T) {
	ac := NewAccessController("account-123")

	tests := []struct {
		name     string
		perm     s3store.Permission
		action   string
		isObject bool
		want     bool
	}{
		{
			name:     "full control matches any action",
			perm:     s3store.PermissionFullControl,
			action:   "s3:GetObject",
			isObject: true,
			want:     true,
		},
		{
			name:     "read permission matches GetObject",
			perm:     s3store.PermissionRead,
			action:   "s3:GetObject",
			isObject: true,
			want:     true,
		},
		{
			name:     "read permission matches ListBucket",
			perm:     s3store.PermissionRead,
			action:   "s3:ListBucket",
			isObject: false,
			want:     true,
		},
		{
			name:     "read permission does not match PutObject",
			perm:     s3store.PermissionRead,
			action:   "s3:PutObject",
			isObject: true,
			want:     false,
		},
		{
			name:     "write permission matches PutObject",
			perm:     s3store.PermissionWrite,
			action:   "s3:PutObject",
			isObject: true,
			want:     true,
		},
		{
			name:     "write permission matches DeleteObject",
			perm:     s3store.PermissionWrite,
			action:   "s3:DeleteObject",
			isObject: true,
			want:     true,
		},
		{
			name:     "read_acp matches GetObjectAcl",
			perm:     s3store.PermissionReadACP,
			action:   "s3:GetObjectAcl",
			isObject: true,
			want:     true,
		},
		{
			name:     "write_acp matches PutObjectAcl",
			perm:     s3store.PermissionWriteACP,
			action:   "s3:PutObjectAcl",
			isObject: true,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ac.permissionMatchesAction(tt.perm, tt.action, tt.isObject)
			if got != tt.want {
				t.Errorf("permissionMatchesAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccessController_ExtractAccountFromPrincipal(t *testing.T) {
	ac := NewAccessController("account-123")

	tests := []struct {
		name      string
		principal string
		want      string
	}{
		{
			name:      "valid IAM user ARN",
			principal: "arn:aws:iam::123456789012:user/test-user",
			want:      "123456789012",
		},
		{
			name:      "valid IAM role ARN",
			principal: "arn:aws:iam::123456789012:role/test-role",
			want:      "123456789012",
		},
		{
			name:      "wildcard principal",
			principal: "*",
			want:      "",
		},
		{
			name:      "empty principal",
			principal: "",
			want:      "",
		},
		{
			name:      "invalid ARN format",
			principal: "not-an-arn",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ac.extractAccountFromPrincipal(tt.principal)
			if got != tt.want {
				t.Errorf("extractAccountFromPrincipal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccessController_AclContainsPublicAccess(t *testing.T) {
	ac := NewAccessController("account-123")

	tests := []struct {
		name string
		acl  *s3store.AccessControlPolicy
		want bool
	}{
		{
			name: "nil ACL",
			acl:  nil,
			want: false,
		},
		{
			name: "ACL with public grant",
			acl: &s3store.AccessControlPolicy{
				Grants: []*s3store.Grant{
					{
						Grantee: &s3store.Grantee{
							Type: s3store.GranteeTypeGroup,
							URI:  s3store.AllUsersGroup,
						},
						Permission: s3store.PermissionRead,
					},
				},
			},
			want: true,
		},
		{
			name: "ACL without public grant",
			acl: &s3store.AccessControlPolicy{
				Grants: []*s3store.Grant{
					{
						Grantee: &s3store.Grantee{
							Type: s3store.GranteeTypeCanonicalUser,
							ID:   "user-123",
						},
						Permission: s3store.PermissionFullControl,
					},
				},
			},
			want: false,
		},
		{
			name: "ACL with authenticated users only",
			acl: &s3store.AccessControlPolicy{
				Grants: []*s3store.Grant{
					{
						Grantee: &s3store.Grantee{
							Type: s3store.GranteeTypeGroup,
							URI:  s3store.AuthenticatedUsersGroup,
						},
						Permission: s3store.PermissionRead,
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ac.aclContainsPublicAccess(tt.acl)
			if got != tt.want {
				t.Errorf("aclContainsPublicAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetermineObjectAction(t *testing.T) {
	tests := []struct {
		method     string
		query      map[string]string
		hasCopySrc bool
		want       string
	}{
		{"GET", nil, false, "s3:GetObject"},
		{"HEAD", nil, false, "s3:GetObject"},
		{"PUT", nil, false, "s3:PutObject"},
		{"DELETE", nil, false, "s3:DeleteObject"},
		{"GET", map[string]string{"acl": ""}, false, "s3:GetObjectAcl"},
		{"PUT", map[string]string{"acl": ""}, false, "s3:PutObjectAcl"},
		{"GET", map[string]string{"tagging": ""}, false, "s3:GetObjectTagging"},
		{"PUT", map[string]string{"tagging": ""}, false, "s3:PutObjectTagging"},
		{"DELETE", map[string]string{"tagging": ""}, false, "s3:DeleteObjectTagging"},
		{"GET", map[string]string{"uploadId": "123"}, false, "s3:ListMultipartUploadParts"},
		{"PUT", map[string]string{"uploadId": "123"}, false, "s3:UploadPart"},
		{"POST", map[string]string{"uploadId": "123"}, false, "s3:CompleteMultipartUpload"},
		{"DELETE", map[string]string{"uploadId": "123"}, false, "s3:AbortMultipartUpload"},
		{"POST", map[string]string{"uploads": ""}, false, "s3:CreateMultipartUpload"},
		{"PUT", nil, true, "s3:PutObject"},
	}

	for _, tt := range tests {
		t.Run(tt.method+"_"+tt.want, func(t *testing.T) {
			req := createMockRequest(tt.method, tt.query, tt.hasCopySrc)
			got := determineObjectAction(req)
			if got != tt.want {
				t.Errorf("determineObjectAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createMockRequest(method string, query map[string]string, hasCopySrc bool) *http.Request {
	req, _ := http.NewRequest(method, "/bucket/key", nil)
	if query != nil {
		q := req.URL.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	if hasCopySrc {
		req.Header.Set("x-amz-copy-source", "/source-bucket/source-key")
	}
	return req
}
