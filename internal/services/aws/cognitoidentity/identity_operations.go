package cognitoidentity

import (
	"context"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"

	"github.com/google/uuid"
)

// GetId obtains a unique identity ID for a Cognito identity pool.
func (s *CognitoIdentityService) GetId(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := getParam(req, "IdentityPoolId")
	if poolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetIdentityPool(poolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	identity := cognitoidentitystore.NewIdentity(poolID)
	if logins := parseLogins(req); len(logins) > 0 {
		identity.Logins = logins
	}

	if err := store.CreateIdentity(identity); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"IdentityId": identity.ID,
	}, nil
}

// GetCredentialsForIdentity returns temporary credentials for an identity.
func (s *CognitoIdentityService) GetCredentialsForIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identityID := getParam(req, "IdentityId")
	if identityID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetIdentityByID(identityID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	now := time.Now().UTC()
	expiration := now.Add(1 * time.Hour)

	credentials := map[string]interface{}{
		"AccessKeyId":  "ASIA" + uuid.New().String()[:16],
		"SecretKey":    uuid.New().String(),
		"SessionToken": uuid.New().String() + uuid.New().String(),
		"Expiration":   expiration.Unix(),
	}

	return map[string]interface{}{
		"IdentityId":  identityID,
		"Credentials": credentials,
	}, nil
}

// DescribeIdentity returns information about a Cognito identity.
func (s *CognitoIdentityService) DescribeIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identityID := getParam(req, "IdentityId")
	if identityID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetIdentityByID(identityID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{
		"IdentityId":       identity.ID,
		"CreationDate":     identity.CreationDate.Unix(),
		"LastModifiedDate": identity.LastModifiedDate.Unix(),
		"Logins":           identity.Logins,
	}, nil
}

func parseLogins(req *request.ParsedRequest) map[string]string {
	if val, ok := req.Parameters["Logins"]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			result := make(map[string]string)
			for k, v := range m {
				if s, ok := v.(string); ok {
					result[k] = s
				}
			}
			return result
		}
	}
	return nil
}
