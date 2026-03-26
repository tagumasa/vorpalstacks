// Package s3 provides S3 (Simple Storage Service) functionality for vorpalstacks.
package s3

import (
	"fmt"
	"strings"

	s3store "vorpalstacks/internal/store/aws/s3"
)

// BuildACLXML builds an XML representation of an S3 access control list.
func BuildACLXML(owner *s3store.ACLOwner, grants []*s3store.Grant) string {
	var result strings.Builder
	result.WriteString(`<AccessControlPolicy xmlns="http://s3.amazonaws.com/doc/2006-03-01/">`)
	if owner != nil {
		result.WriteString(`<Owner><ID>`)
		result.WriteString(xmlEscape(owner.ID))
		result.WriteString(`</ID><DisplayName>`)
		result.WriteString(xmlEscape(owner.DisplayName))
		result.WriteString(`</DisplayName></Owner>`)
	}
	if len(grants) > 0 {
		result.WriteString(`<AccessControlList>`)
		for _, g := range grants {
			result.WriteString(`<Grant>`)
			if g.Grantee != nil {
				result.WriteString(`<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="`)
				result.WriteString(string(g.Grantee.Type))
				result.WriteString(`">`)
				if g.Grantee.ID != "" {
					result.WriteString(`<ID>`)
					result.WriteString(xmlEscape(g.Grantee.ID))
					result.WriteString(`</ID>`)
				}
				if g.Grantee.DisplayName != "" {
					result.WriteString(`<DisplayName>`)
					result.WriteString(xmlEscape(g.Grantee.DisplayName))
					result.WriteString(`</DisplayName>`)
				}
				if g.Grantee.Email != "" {
					result.WriteString(`<EmailAddress>`)
					result.WriteString(xmlEscape(g.Grantee.Email))
					result.WriteString(`</EmailAddress>`)
				}
				if g.Grantee.URI != "" {
					result.WriteString(`<URI>`)
					result.WriteString(xmlEscape(g.Grantee.URI))
					result.WriteString(`</URI>`)
				}
				result.WriteString(`</Grantee>`)
			}
			result.WriteString(`<Permission>`)
			result.WriteString(string(g.Permission))
			result.WriteString(`</Permission>`)
			result.WriteString(`</Grant>`)
		}
		result.WriteString(`</AccessControlList>`)
	}
	result.WriteString(`</AccessControlPolicy>`)
	return result.String()
}

// CannedACLToPolicy converts a canned ACL string to an access control policy.
func CannedACLToPolicy(cannedACL string, owner *s3store.ACLOwner) (*s3store.AccessControlPolicy, error) {
	if owner == nil {
		owner = &s3store.ACLOwner{ID: "owner", DisplayName: "owner"}
	}

	grants := []*s3store.Grant{
		{
			Grantee:    &s3store.Grantee{Type: s3store.GranteeTypeCanonicalUser, ID: owner.ID, DisplayName: owner.DisplayName},
			Permission: s3store.PermissionFullControl,
		},
	}

	switch cannedACL {
	case "private":
	case "public-read":
		grants = append(grants, &s3store.Grant{
			Grantee:    &s3store.Grantee{Type: s3store.GranteeTypeGroup, URI: s3store.AllUsersGroup},
			Permission: s3store.PermissionRead,
		})
	case "public-read-write":
		grants = append(grants,
			&s3store.Grant{Grantee: &s3store.Grantee{Type: s3store.GranteeTypeGroup, URI: s3store.AllUsersGroup}, Permission: s3store.PermissionRead},
			&s3store.Grant{Grantee: &s3store.Grantee{Type: s3store.GranteeTypeGroup, URI: s3store.AllUsersGroup}, Permission: s3store.PermissionWrite},
		)
	case "authenticated-read":
		grants = append(grants, &s3store.Grant{
			Grantee:    &s3store.Grantee{Type: s3store.GranteeTypeGroup, URI: s3store.AuthenticatedUsersGroup},
			Permission: s3store.PermissionRead,
		})
	case "log-delivery-write":
		grants = append(grants,
			&s3store.Grant{Grantee: &s3store.Grantee{Type: s3store.GranteeTypeGroup, URI: s3store.LogDeliveryGroup}, Permission: s3store.PermissionReadACP},
			&s3store.Grant{Grantee: &s3store.Grantee{Type: s3store.GranteeTypeGroup, URI: s3store.LogDeliveryGroup}, Permission: s3store.PermissionWrite},
		)
	case "bucket-owner-read", "bucket-owner-full-control", "aws-exec-read":
	default:
		return nil, fmt.Errorf("invalid canned ACL: %s", cannedACL)
	}

	return &s3store.AccessControlPolicy{Owner: owner, Grants: grants}, nil
}

// ParseGrantHeaders parses grant headers into a list of grants.
func ParseGrantHeaders(grantFullControl, grantRead, grantReadACP, grantWrite, grantWriteACP string) []*s3store.Grant {
	var grants []*s3store.Grant

	parseGrantees := func(header string, permission s3store.Permission) {
		if header == "" {
			return
		}
		grantees := strings.Split(header, ",")
		for _, g := range grantees {
			g = strings.TrimSpace(g)
			if g == "" {
				continue
			}
			grantee := ParseGrantee(g)
			if grantee != nil {
				grants = append(grants, &s3store.Grant{Grantee: grantee, Permission: permission})
			}
		}
	}

	parseGrantees(grantFullControl, s3store.PermissionFullControl)
	parseGrantees(grantRead, s3store.PermissionRead)
	parseGrantees(grantReadACP, s3store.PermissionReadACP)
	parseGrantees(grantWrite, s3store.PermissionWrite)
	parseGrantees(grantWriteACP, s3store.PermissionWriteACP)

	return grants
}

// ParseGrantee parses a grantee string into a Grantee struct.
func ParseGrantee(s string) *s3store.Grantee {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return nil
	}

	grantType := strings.TrimSpace(parts[0])
	id := strings.Trim(strings.TrimSpace(parts[1]), `"`)

	switch grantType {
	case "id":
		return &s3store.Grantee{Type: s3store.GranteeTypeCanonicalUser, ID: id, DisplayName: "owner"}
	case "uri":
		return &s3store.Grantee{Type: s3store.GranteeTypeGroup, URI: id}
	case "emailAddress":
		return &s3store.Grantee{Type: s3store.GranteeTypeAmazonCustomerByEmail, Email: id}
	default:
		return nil
	}
}
