package iot

import (
	"time"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
)

func (s *IotStore) CreateAuthorizer(a *Authorizer) (*Authorizer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if a.AuthorizerName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.Authorizer{}
	if err := s.authorizersBase.GetProto(a.AuthorizerName, existing); err == nil {
		return nil, ErrAuthorizerAlreadyExists
	}
	a.AuthorizerARN = BuildAuthorizerARN(s.accountID, s.region, a.AuthorizerName)
	return a, s.authorizerPS.Create(a)
}

func (s *IotStore) GetAuthorizer(name string) (*Authorizer, error) {
	return s.authorizerPS.Get(name)
}

func (s *IotStore) UpdateAuthorizer(name string, opts AuthorizerUpdateOpts) (*Authorizer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, err := s.authorizerPS.Get(name)
	if err != nil {
		return nil, ErrAuthorizerNotFound
	}
	if opts.FunctionARN != "" {
		existing.AuthorizerFunctionARN = opts.FunctionARN
	}
	if opts.TokenName != "" {
		existing.TokenName = opts.TokenName
	}
	if opts.TokenSignature != "" {
		existing.TokenSignature = opts.TokenSignature
	}
	if opts.EnableCaching != nil {
		existing.EnableCachingForHTTP = *opts.EnableCaching
	}
	if opts.Status != "" {
		existing.Status = opts.Status == "ACTIVE"
	}
	existing.LastModifiedDate = time.Now().UTC()
	return existing, s.authorizerPS.Update(existing)
}

func (s *IotStore) DeleteAuthorizer(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.authorizerPS.DeleteIfExists(name)
}

// ListAuthorizers returns all custom authorizers.
func (s *IotStore) ListAuthorizers(opts common.ListOptions) (*common.ListResult[Authorizer], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.authorizersBase, opts, func() *pb.Authorizer { return &pb.Authorizer{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*Authorizer, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToAuthorizer(p))
	}
	return &common.ListResult[Authorizer]{Items: items, NextMarker: result.NextMarker}, nil
}

// CreateProvisioningTemplate persists a new provisioning template.
