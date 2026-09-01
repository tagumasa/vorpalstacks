package sesv2

import (
	pagination "vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// ---------------------------------------------------------------------------
// Input DTOs — dedicated IP pool family
// ---------------------------------------------------------------------------

// CreateDedicatedIpPoolInput carries the pool-create members.
type CreateDedicatedIpPoolInput struct {
	PoolName    string
	ScalingMode string
	Tags        []tags.Tag
}

// ---------------------------------------------------------------------------
// Core functions — dedicated IP pool family
// ---------------------------------------------------------------------------

// createDedicatedIpPoolCore is the single entry point for pool creation.
// Per Smithy com.amazonaws.sesv2#CreateDedicatedIpPoolRequest, Tags is a
// top-level member that must be persisted alongside the pool.
func (s *SESv2Service) createDedicatedIpPoolCore(store sesv2store.SESv2StoreInterface, in CreateDedicatedIpPoolInput) error {
	if in.PoolName == "" {
		return ErrMissingParameter
	}
	if !validatePoolName(in.PoolName) {
		return ErrBadRequest
	}

	scalingMode := in.ScalingMode
	if scalingMode == "" {
		scalingMode = "STANDARD"
	}
	if !validateScalingMode(scalingMode) {
		return ErrBadRequest
	}

	pool := &sesv2store.DedicatedIpPool{
		PoolName:    in.PoolName,
		ScalingMode: scalingMode,
	}

	if err := store.CreateDedicatedIpPool(pool); err != nil {
		return err
	}

	if len(in.Tags) > 0 {
		arn := store.BuildDedicatedIpPoolArn(in.PoolName)
		if err := store.TagFromSlice(arn, in.Tags); err != nil {
			return err
		}
	}

	return nil
}

// getDedicatedIpPoolCore is the single entry point for reading a pool.
func (s *SESv2Service) getDedicatedIpPoolCore(store sesv2store.SESv2StoreInterface, poolName string) (map[string]interface{}, error) {
	if poolName == "" {
		return nil, ErrMissingParameter
	}

	pool, err := store.GetDedicatedIpPool(poolName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DedicatedIpPool": map[string]interface{}{
			"PoolName":    pool.PoolName,
			"ScalingMode": pool.ScalingMode,
		},
	}, nil
}

// deleteDedicatedIpPoolCore is the single entry point for deleting a pool.
func (s *SESv2Service) deleteDedicatedIpPoolCore(store sesv2store.SESv2StoreInterface, poolName string) error {
	if poolName == "" {
		return ErrMissingParameter
	}
	return store.DeleteDedicatedIpPool(poolName)
}

// listDedicatedIpPoolsCore is the single entry point for listing pools.
func (s *SESv2Service) listDedicatedIpPoolsCore(store sesv2store.SESv2StoreInterface, maxItems int, nextToken string) (map[string]interface{}, error) {
	result, err := store.ListDedicatedIpPools(common.ListOptions{
		MaxItems: maxItems,
		Marker:   nextToken,
	})
	if err != nil {
		return nil, err
	}

	pools := make([]string, 0, len(result.Items))
	for _, pool := range result.Items {
		pools = append(pools, pool.PoolName)
	}

	resp := map[string]interface{}{
		"DedicatedIpPools": pools,
	}

	pagination.SetNextToken(resp, "NextToken", result.NextMarker)

	return resp, nil
}

// ---------------------------------------------------------------------------
// HTTP handlers — parse → DTO → Core → serialise
// ---------------------------------------------------------------------------
