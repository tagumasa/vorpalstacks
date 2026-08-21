package athena

import (
	"net/http"
	"vorpalstacks/internal/common/defaults"

	pb "vorpalstacks/internal/pb/aws/athena"
	"vorpalstacks/internal/utils/timeutils"
)

// getStores extracts the region from request headers and returns the full
// athenaStores for that region.
func (h *AdminHandler) getStores(headers http.Header) (*athenaStores, error) {
	region := defaults.GetRegionFromHeader(headers)
	return h.service.getStoresForRegion(region)
}

// protoToCreateInput converts the protobuf CreateWorkGroupInput message into
// the transport-neutral WorkGroupCreateInput DTO used by createWorkGroupCore.
func protoToCreateInput(msg *pb.CreateWorkGroupInput) WorkGroupCreateInput {
	input := WorkGroupCreateInput{
		Name:        msg.Name,
		Description: msg.Description,
	}
	if msg.Configuration != nil {
		protoCfg := msg.Configuration
		cfg := &WorkGroupConfigInput{
			EnforceConfig:           protoCfg.GetEnforceworkgroupconfiguration(),
			PublishMetrics:          protoCfg.GetPublishcloudwatchmetricsenabled(),
			RequesterPaysEnabled:    protoCfg.GetRequesterpaysenabled(),
			AdditionalConfiguration: protoCfg.GetAdditionalconfiguration(),
			ExecutionRole:           protoCfg.GetExecutionrole(),
		}
		if protoCfg.Bytesscannedcutoffperquery != nil {
			v := int64(protoCfg.GetBytesscannedcutoffperquery())
			cfg.BytesScannedCutoff = &v
		}
		if protoCfg.Resultconfiguration != nil {
			cfg.OutputLocation = protoCfg.Resultconfiguration.Outputlocation
		}
		if protoCfg.Engineversion != nil {
			cfg.EngineVersionSelected = protoCfg.Engineversion.Selectedengineversion
			cfg.EngineVersionEffective = protoCfg.Engineversion.Effectiveengineversion
		}
		if protoCfg.Customercontentencryptionconfiguration != nil {
			cfg.CustomerContentEncryptionKmsKey = protoCfg.Customercontentencryptionconfiguration.GetKmskey()
		}
		if protoCfg.Enableminimumencryptionconfiguration != nil {
			cfg.EnableMinimumEncryptionConfiguration = protoCfg.GetEnableminimumencryptionconfiguration()
		}
		input.Config = cfg
	}

	for _, tag := range msg.Tags {
		if input.Tags == nil {
			input.Tags = make(map[string]string)
		}
		input.Tags[tag.Key] = tag.Value
	}

	return input
}

// toPbWorkGroupSummaries converts service-layer WorkGroupOut items into the
// protobuf WorkGroupSummary messages for the admin handler response.
func toPbWorkGroupSummaries(items []WorkGroupOut) []*pb.WorkGroupSummary {
	summaries := make([]*pb.WorkGroupSummary, 0, len(items))
	for _, wg := range items {
		state := pb.WorkGroupState_WORK_GROUP_STATE_DISABLED
		if wg.State == "ENABLED" {
			state = pb.WorkGroupState_WORK_GROUP_STATE_ENABLED
		}
		summary := &pb.WorkGroupSummary{
			Name:  wg.Name,
			State: state,
		}
		if wg.Description != "" {
			summary.Description = wg.Description
		}
		if !wg.CreationTime.IsZero() {
			summary.Creationtime = wg.CreationTime.Format(timeutils.ISO8601UTCFormat)
		}
		summaries = append(summaries, summary)
	}
	return summaries
}
