package ssm

import (
	"fmt"
	"net/http"

	"google.golang.org/protobuf/proto"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/utils/timeutils"

	pb "vorpalstacks/internal/pb/aws/ssm"
	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

// getStore returns the regional SSM store from the gRPC request headers.
func (h *AdminHandler) getStore(headers http.Header) (*ssmstore.Store, error) {
	region := svccommon.GetRegionFromHeader(headers)
	store, err := h.service.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	return store.(*ssmstore.Store), nil
}

// toPbParameterMetadata converts a store Parameter to the proto
// ParameterMetadata message used by the admin console.
func toPbParameterMetadata(p *ssmstore.Parameter) *pb.ParameterMetadata {
	meta := &pb.ParameterMetadata{
		Name:             p.Name,
		Version:          proto.Int64(p.Version),
		Lastmodifieddate: p.LastModifiedDate.Format(timeutils.ISO8601UTCFormat),
		Datatype:         p.DataType,
		Arn:              p.ARN,
	}
	if p.Description != "" {
		meta.Description = p.Description
	}
	if p.KeyID != "" {
		meta.Keyid = p.KeyID
	}
	if p.AllowedPattern != "" {
		meta.Allowedpattern = p.AllowedPattern
	}
	switch p.Type {
	case ssmstore.ParameterTypeString:
		meta.Type = pb.ParameterType_PARAMETER_TYPE_STRING
	case ssmstore.ParameterTypeStringList:
		meta.Type = pb.ParameterType_PARAMETER_TYPE_STRING_LIST
	case ssmstore.ParameterTypeSecureString:
		meta.Type = pb.ParameterType_PARAMETER_TYPE_SECURE_STRING
	}
	switch p.Tier {
	case ssmstore.ParameterTierStandard:
		meta.Tier = pb.ParameterTier_PARAMETER_TIER_STANDARD
	case ssmstore.ParameterTierAdvanced:
		meta.Tier = pb.ParameterTier_PARAMETER_TIER_ADVANCED
	case ssmstore.ParameterTierIntelligentTiering:
		meta.Tier = pb.ParameterTier_PARAMETER_TIER_INTELLIGENT_TIERING
	}
	return meta
}

// toStoreFilters converts proto filter messages to store ParameterFilter
// structs, validating each key.
func toStoreFilters(pbFilters []*pb.ParametersFilter) ([]ssmstore.ParameterFilter, error) {
	var filters []ssmstore.ParameterFilter
	for _, f := range pbFilters {
		key := f.Key.String()
		if key == "" || key == "PARAMETERS_FILTER_KEY_INVALID" {
			continue
		}
		if !ssmstore.ValidateParameterFilterKey(key) {
			return nil, fmt.Errorf("invalid filter key: %s", key)
		}
		filters = append(filters, ssmstore.ParameterFilter{
			Key:    key,
			Option: "",
			Values: f.Values,
		})
	}
	return filters, nil
}

// toPbParameterMetadataFromMeta converts a store ParameterMetadata to the
// proto ParameterMetadata message.
func toPbParameterMetadataFromMeta(p *ssmstore.ParameterMetadata) *pb.ParameterMetadata {
	meta := &pb.ParameterMetadata{
		Name:             p.Name,
		Version:          proto.Int64(p.Version),
		Lastmodifieddate: p.LastModifiedDate.Format(timeutils.ISO8601UTCFormat),
		Datatype:         p.DataType,
		Arn:              p.ARN,
	}
	if p.Description != "" {
		meta.Description = p.Description
	}
	if p.KeyID != "" {
		meta.Keyid = p.KeyID
	}
	if p.AllowedPattern != "" {
		meta.Allowedpattern = p.AllowedPattern
	}
	switch p.Type {
	case ssmstore.ParameterTypeString:
		meta.Type = pb.ParameterType_PARAMETER_TYPE_STRING
	case ssmstore.ParameterTypeStringList:
		meta.Type = pb.ParameterType_PARAMETER_TYPE_STRING_LIST
	case ssmstore.ParameterTypeSecureString:
		meta.Type = pb.ParameterType_PARAMETER_TYPE_SECURE_STRING
	}
	switch p.Tier {
	case ssmstore.ParameterTierStandard:
		meta.Tier = pb.ParameterTier_PARAMETER_TIER_STANDARD
	case ssmstore.ParameterTierAdvanced:
		meta.Tier = pb.ParameterTier_PARAMETER_TIER_ADVANCED
	case ssmstore.ParameterTierIntelligentTiering:
		meta.Tier = pb.ParameterTier_PARAMETER_TIER_INTELLIGENT_TIERING
	}
	return meta
}
