package iot

import (
	"strings"
	"context"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
)
func (s *IotStore) CreatePolicy(policy *Policy) (*Policy, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if policy.PolicyName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.Policy{}
	if err := s.policiesBase.GetProto(policy.PolicyName, existing); err == nil {
		return nil, ErrPolicyAlreadyExists
	}
	policy.PolicyARN = BuildPolicyARN(s.accountID, s.region, policy.PolicyName)
	return policy, s.policyPS.Create(policy)
}

func (s *IotStore) GetPolicy(policyName string) (*Policy, error) {
	return s.policyPS.Get(policyName)
}

func (s *IotStore) DeletePolicy(policyName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	principals, err := s.listPrincipalsForPolicyLocked(policyName)
	if err != nil {
		return err
	}
	if len(principals) > 0 {
		return ErrDeleteConflict
	}
	return s.policyPS.DeleteIfExists(policyName)
}

func (s *IotStore) ListPolicies(opts common.ListOptions) (*common.ListResult[Policy], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.policiesBase, opts, func() *pb.Policy { return &pb.Policy{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*Policy, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToPolicy(p))
	}
	return &common.ListResult[Policy]{Items: items, NextMarker: result.NextMarker}, nil
}

func (s *IotStore) AttachPolicyToPrincipal(policyName, principal string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.policiesBase.Exists(policyName) {
		return ErrPolicyNotFound
	}
	ak := policyAttachKey(policyName, principal)
	pk := policyPrincipalKey(policyName, principal)
	paBucket := "iot-policy-attach" + s.rs
	ppBucket := "iot-policy-principal" + s.rs
	return s.ts.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(paBucket).Put([]byte(ak), []byte("1")); err != nil {
			return err
		}
		return txn.Bucket(ppBucket).Put([]byte(pk), []byte("1"))
	})
}

func (s *IotStore) DetachPolicyFromPrincipal(policyName, principal string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.policiesBase.Exists(policyName) {
		return ErrPolicyNotFound
	}
	ak := policyAttachKey(policyName, principal)
	pk := policyPrincipalKey(policyName, principal)
	paBucket := "iot-policy-attach" + s.rs
	ppBucket := "iot-policy-principal" + s.rs
	return s.ts.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(paBucket).Delete([]byte(ak)); err != nil {
			return err
		}
		return txn.Bucket(ppBucket).Delete([]byte(pk))
	})
}

// listPrincipalsForPolicyLocked returns principals attached to a policy.
// Caller must hold s.mu.
func (s *IotStore) listPrincipalsForPolicyLocked(policyName string) ([]string, error) {
	var principals []string
	err := s.policyPrincipalBase.ScanPrefix(policyName+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			principals = append(principals, parts[1])
		}
		return nil
	})
	return principals, err
}

func (s *IotStore) ListPrincipalsForPolicy(policyName string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.listPrincipalsForPolicyLocked(policyName)
}

func (s *IotStore) ListPoliciesForPrincipal(principal string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var policies []string
	err := s.policyAttachBase.ScanPrefix(principal+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			policies = append(policies, parts[1])
		}
		return nil
	})
	return policies, err
}

func policyAttachKey(policyName, principal string) string {
	return principal + "\x00" + policyName
}

func policyPrincipalKey(policyName, principal string) string {
	return policyName + "\x00" + principal
}
