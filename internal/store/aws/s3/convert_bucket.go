package s3

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	pb "vorpalstacks/internal/pb/storage/storage_s3"
)

// BucketToProto converts an internal Bucket to a protobuf Bucket.
// Returns nil if the input bucket is nil.
func BucketToProto(b *Bucket) *pb.Bucket {
	if b == nil {
		return nil
	}
	return &pb.Bucket{
		Name:                      b.Name,
		Region:                    b.Region,
		CreationDate:              timestamppb.New(b.CreationDate),
		Acl:                       accessControlPolicyToProto(b.ACL),
		ObjectLockEnabled:         b.ObjectLockEnabled,
		ObjectLockConfig:          objectLockConfigurationToProto(b.ObjectLockConfig),
		VersioningStatus:          bucketVersioningStatusToProto(b.VersioningStatus),
		MfaDelete:                 b.MFADelete,
		EncryptionConfig:          encryptionConfigToProto(b.EncryptionConfig),
		LifecycleConfiguration:    lifecycleConfigurationToProto(b.LifecycleConfiguration),
		WebsiteConfiguration:      websiteConfigurationToProto(b.WebsiteConfiguration),
		CorsConfiguration:         corsConfigurationToProto(b.CORSConfiguration),
		Policy:                    b.Policy,
		PublicAccessBlock:         publicAccessBlockConfigToProto(b.PublicAccessBlock),
		Tags:                      tagsToProto(b.Tags),
		NotificationConfiguration: notificationConfigurationToProto(b.NotificationConfiguration),
		LoggingConfiguration:      loggingConfigurationToProto(b.LoggingConfiguration),
		OwnershipControls:         ownershipControlsToProto(b.OwnershipControls),
		RequestPayment:            requestPaymentConfigurationToProto(b.RequestPayment),
		AccelerateConfiguration:   accelerateConfigurationToProto(b.AccelerateConfiguration),
		ReplicationConfiguration:  replicationConfigurationToProto(b.ReplicationConfiguration),
		InventoryConfigurations:   inventoryConfigurationsToProto(b.InventoryConfigurations),
		MetricsConfigurations:     metricsConfigurationsToProto(b.MetricsConfigurations),
	}
}

// ProtoToBucket converts a protobuf Bucket to an internal Bucket.
// Returns nil if the input protobuf bucket is nil.
func ProtoToBucket(p *pb.Bucket) *Bucket {
	if p == nil {
		return nil
	}
	return &Bucket{
		Name:                      p.Name,
		Region:                    p.Region,
		CreationDate:              p.CreationDate.AsTime(),
		ACL:                       protoToAccessControlPolicy(p.Acl),
		ObjectLockEnabled:         p.ObjectLockEnabled,
		ObjectLockConfig:          protoToObjectLockConfiguration(p.ObjectLockConfig),
		VersioningStatus:          protoToBucketVersioningStatus(p.VersioningStatus),
		MFADelete:                 p.MfaDelete,
		EncryptionConfig:          protoToEncryptionConfig(p.EncryptionConfig),
		LifecycleConfiguration:    protoToLifecycleConfiguration(p.LifecycleConfiguration),
		WebsiteConfiguration:      protoToWebsiteConfiguration(p.WebsiteConfiguration),
		CORSConfiguration:         protoToCORSConfiguration(p.CorsConfiguration),
		Policy:                    p.Policy,
		PublicAccessBlock:         protoToPublicAccessBlockConfig(p.PublicAccessBlock),
		Tags:                      protoToTags(p.Tags),
		NotificationConfiguration: protoToNotificationConfiguration(p.NotificationConfiguration),
		LoggingConfiguration:      protoToLoggingConfiguration(p.LoggingConfiguration),
		OwnershipControls:         protoToOwnershipControls(p.OwnershipControls),
		RequestPayment:            protoToRequestPaymentConfiguration(p.RequestPayment),
		AccelerateConfiguration:   protoToAccelerateConfiguration(p.AccelerateConfiguration),
		ReplicationConfiguration:  protoToReplicationConfiguration(p.ReplicationConfiguration),
		InventoryConfigurations:   protoToInventoryConfigurations(p.InventoryConfigurations),
		MetricsConfigurations:     protoToMetricsConfigurations(p.MetricsConfigurations),
	}
}

func inventoryConfigurationsToProto(configs map[string]*InventoryConfiguration) []*pb.InventoryConfiguration {
	if len(configs) == 0 {
		return nil
	}
	result := make([]*pb.InventoryConfiguration, 0, len(configs))
	for _, c := range configs {
		if c == nil {
			continue
		}
		result = append(result, inventoryConfigurationToProto(c))
	}
	return result
}

func protoToInventoryConfigurations(configs []*pb.InventoryConfiguration) map[string]*InventoryConfiguration {
	if len(configs) == 0 {
		return nil
	}
	result := make(map[string]*InventoryConfiguration, len(configs))
	for _, c := range configs {
		if c == nil {
			continue
		}
		stored := protoToInventoryConfiguration(c)
		result[stored.ID] = stored
	}
	return result
}

func inventoryConfigurationToProto(c *InventoryConfiguration) *pb.InventoryConfiguration {
	if c == nil {
		return nil
	}
	var lastDelivery *timestamppb.Timestamp
	if !c.LastDelivery.IsZero() {
		lastDelivery = timestamppb.New(c.LastDelivery)
	}
	return &pb.InventoryConfiguration{
		Id:                     c.ID,
		IsEnabled:              c.IsEnabled,
		Filter:                 inventoryFilterToProto(c.Filter),
		IncludedObjectVersions: c.IncludedObjectVersions,
		OptionalFields:         c.OptionalFields,
		Schedule:               inventoryScheduleToProto(c.Schedule),
		Destination:            inventoryDestinationToProto(c.Destination),
		LastDelivery:           lastDelivery,
	}
}

func protoToInventoryConfiguration(p *pb.InventoryConfiguration) *InventoryConfiguration {
	if p == nil {
		return nil
	}
	var lastDelivery time.Time
	if p.LastDelivery != nil {
		lastDelivery = p.LastDelivery.AsTime()
	}
	return &InventoryConfiguration{
		ID:                     p.Id,
		IsEnabled:              p.IsEnabled,
		Filter:                 protoToInventoryFilter(p.Filter),
		IncludedObjectVersions: p.IncludedObjectVersions,
		OptionalFields:         p.OptionalFields,
		Schedule:               protoToInventorySchedule(p.Schedule),
		Destination:            protoToInventoryDestination(p.Destination),
		LastDelivery:           lastDelivery,
	}
}

func inventoryFilterToProto(f *InventoryFilter) *pb.InventoryFilter {
	if f == nil {
		return nil
	}
	return &pb.InventoryFilter{Prefix: f.Prefix}
}

func protoToInventoryFilter(p *pb.InventoryFilter) *InventoryFilter {
	if p == nil {
		return nil
	}
	return &InventoryFilter{Prefix: p.Prefix}
}

func inventoryScheduleToProto(s *InventorySchedule) *pb.InventorySchedule {
	if s == nil {
		return nil
	}
	return &pb.InventorySchedule{Frequency: s.Frequency}
}

func protoToInventorySchedule(p *pb.InventorySchedule) *InventorySchedule {
	if p == nil {
		return nil
	}
	return &InventorySchedule{Frequency: p.Frequency}
}

func inventoryDestinationToProto(d *InventoryDestination) *pb.InventoryDestination {
	if d == nil {
		return nil
	}
	return &pb.InventoryDestination{
		S3BucketDestination: inventoryS3BucketDestinationToProto(d.S3BucketDestination),
	}
}

func protoToInventoryDestination(p *pb.InventoryDestination) *InventoryDestination {
	if p == nil {
		return nil
	}
	return &InventoryDestination{
		S3BucketDestination: protoToInventoryS3BucketDestination(p.S3BucketDestination),
	}
}

func inventoryS3BucketDestinationToProto(d *InventoryS3BucketDestination) *pb.InventoryS3BucketDestination {
	if d == nil {
		return nil
	}
	return &pb.InventoryS3BucketDestination{
		AccountId:  d.AccountID,
		Bucket:     d.Bucket,
		Format:     d.Format,
		Prefix:     d.Prefix,
		Encryption: inventoryEncryptionToProto(d.Encryption),
	}
}

func protoToInventoryS3BucketDestination(p *pb.InventoryS3BucketDestination) *InventoryS3BucketDestination {
	if p == nil {
		return nil
	}
	return &InventoryS3BucketDestination{
		AccountID:  p.AccountId,
		Bucket:     p.Bucket,
		Format:     p.Format,
		Prefix:     p.Prefix,
		Encryption: protoToInventoryEncryption(p.Encryption),
	}
}

func inventoryEncryptionToProto(e *InventoryEncryption) *pb.InventoryEncryption {
	if e == nil {
		return nil
	}
	return &pb.InventoryEncryption{
		SseS3:  e.SSES3,
		SseKms: inventorySSEKMSToProto(e.SSEKMS),
	}
}

func protoToInventoryEncryption(p *pb.InventoryEncryption) *InventoryEncryption {
	if p == nil {
		return nil
	}
	return &InventoryEncryption{
		SSES3:  p.SseS3,
		SSEKMS: protoToInventorySSEKMS(p.SseKms),
	}
}

func inventorySSEKMSToProto(k *InventorySSEKMS) *pb.InventorySSEKMS {
	if k == nil {
		return nil
	}
	return &pb.InventorySSEKMS{KeyId: k.KeyID}
}

func protoToInventorySSEKMS(p *pb.InventorySSEKMS) *InventorySSEKMS {
	if p == nil {
		return nil
	}
	return &InventorySSEKMS{KeyID: p.KeyId}
}

func metricsConfigurationsToProto(configs map[string]*MetricsConfiguration) []*pb.MetricsConfiguration {
	if len(configs) == 0 {
		return nil
	}
	result := make([]*pb.MetricsConfiguration, 0, len(configs))
	for _, c := range configs {
		if c == nil {
			continue
		}
		result = append(result, metricsConfigurationToProto(c))
	}
	return result
}

func protoToMetricsConfigurations(configs []*pb.MetricsConfiguration) map[string]*MetricsConfiguration {
	if len(configs) == 0 {
		return nil
	}
	result := make(map[string]*MetricsConfiguration, len(configs))
	for _, c := range configs {
		if c == nil {
			continue
		}
		stored := protoToMetricsConfiguration(c)
		result[stored.ID] = stored
	}
	return result
}

func metricsConfigurationToProto(c *MetricsConfiguration) *pb.MetricsConfiguration {
	if c == nil {
		return nil
	}
	return &pb.MetricsConfiguration{
		Id:     c.ID,
		Filter: metricsFilterToProto(c.Filter),
	}
}

func protoToMetricsConfiguration(p *pb.MetricsConfiguration) *MetricsConfiguration {
	if p == nil {
		return nil
	}
	return &MetricsConfiguration{
		ID:     p.Id,
		Filter: protoToMetricsFilter(p.Filter),
	}
}

func metricsFilterToProto(f *MetricsFilter) *pb.MetricsFilter {
	if f == nil {
		return nil
	}
	return &pb.MetricsFilter{
		Prefix:         f.Prefix,
		Tag:            tagToProtoPtr(f.Tag),
		AccessPointArn: f.AccessPointArn,
		And:            metricsAndOperatorToProto(f.And),
	}
}

func protoToMetricsFilter(p *pb.MetricsFilter) *MetricsFilter {
	if p == nil {
		return nil
	}
	return &MetricsFilter{
		Prefix:         p.Prefix,
		Tag:            protoToTagPtr(p.Tag),
		AccessPointArn: p.AccessPointArn,
		And:            protoToMetricsAndOperator(p.And),
	}
}

func metricsAndOperatorToProto(a *MetricsAndOperator) *pb.MetricsAndOperator {
	if a == nil {
		return nil
	}
	return &pb.MetricsAndOperator{
		Prefix:         a.Prefix,
		Tags:           tagsToProto(a.Tags),
		AccessPointArn: a.AccessPointArn,
	}
}

func protoToMetricsAndOperator(p *pb.MetricsAndOperator) *MetricsAndOperator {
	if p == nil {
		return nil
	}
	return &MetricsAndOperator{
		Prefix:         p.Prefix,
		Tags:           protoToTags(p.Tags),
		AccessPointArn: p.AccessPointArn,
	}
}

func publicAccessBlockConfigToProto(c *PublicAccessBlockConfig) *pb.PublicAccessBlockConfig {
	if c == nil {
		return nil
	}
	return &pb.PublicAccessBlockConfig{
		BlockPublicAcls:       c.BlockPublicAcls,
		BlockPublicPolicy:     c.BlockPublicPolicy,
		IgnorePublicAcls:      c.IgnorePublicAcls,
		RestrictPublicBuckets: c.RestrictPublicBuckets,
	}
}

func protoToPublicAccessBlockConfig(p *pb.PublicAccessBlockConfig) *PublicAccessBlockConfig {
	if p == nil {
		return nil
	}
	return &PublicAccessBlockConfig{
		BlockPublicAcls:       p.BlockPublicAcls,
		BlockPublicPolicy:     p.BlockPublicPolicy,
		IgnorePublicAcls:      p.IgnorePublicAcls,
		RestrictPublicBuckets: p.RestrictPublicBuckets,
	}
}

func loggingConfigurationToProto(c *LoggingConfiguration) *pb.LoggingConfiguration {
	if c == nil {
		return nil
	}
	return &pb.LoggingConfiguration{
		TargetBucket: c.TargetBucket,
		TargetPrefix: c.TargetPrefix,
		TargetGrants: targetGrantsToProto(c.TargetGrants),
	}
}

func protoToLoggingConfiguration(p *pb.LoggingConfiguration) *LoggingConfiguration {
	if p == nil {
		return nil
	}
	return &LoggingConfiguration{
		TargetBucket: p.TargetBucket,
		TargetPrefix: p.TargetPrefix,
		TargetGrants: protoToTargetGrants(p.TargetGrants),
	}
}

func targetGrantsToProto(grants []TargetGrant) []*pb.TargetGrant {
	if grants == nil {
		return nil
	}
	result := make([]*pb.TargetGrant, len(grants))
	for i, g := range grants {
		result[i] = targetGrantToProto(&g)
	}
	return result
}

func protoToTargetGrants(p []*pb.TargetGrant) []TargetGrant {
	if p == nil {
		return nil
	}
	result := make([]TargetGrant, len(p))
	for i, g := range p {
		result[i] = *protoToTargetGrant(g)
	}
	return result
}

func targetGrantToProto(g *TargetGrant) *pb.TargetGrant {
	if g == nil {
		return nil
	}
	return &pb.TargetGrant{
		Grantee:    granteeToProto(g.Grantee),
		Permission: permissionToProto(g.Permission),
	}
}

func protoToTargetGrant(p *pb.TargetGrant) *TargetGrant {
	if p == nil {
		return nil
	}
	return &TargetGrant{
		Grantee:    protoToGrantee(p.Grantee),
		Permission: protoToPermission(p.Permission),
	}
}

func ownershipControlsToProto(o *OwnershipControls) *pb.OwnershipControls {
	if o == nil {
		return nil
	}
	return &pb.OwnershipControls{
		Rules: ownershipControlsRulesToProto(o.Rules),
	}
}

func protoToOwnershipControls(p *pb.OwnershipControls) *OwnershipControls {
	if p == nil {
		return nil
	}
	return &OwnershipControls{
		Rules: protoToOwnershipControlsRules(p.Rules),
	}
}

func ownershipControlsRulesToProto(rules []OwnershipControlsRule) []*pb.OwnershipControlsRule {
	if rules == nil {
		return nil
	}
	result := make([]*pb.OwnershipControlsRule, len(rules))
	for i, r := range rules {
		result[i] = ownershipControlsRuleToProto(&r)
	}
	return result
}

func protoToOwnershipControlsRules(p []*pb.OwnershipControlsRule) []OwnershipControlsRule {
	if p == nil {
		return nil
	}
	result := make([]OwnershipControlsRule, len(p))
	for i, r := range p {
		result[i] = *protoToOwnershipControlsRule(r)
	}
	return result
}

func ownershipControlsRuleToProto(r *OwnershipControlsRule) *pb.OwnershipControlsRule {
	if r == nil {
		return nil
	}
	return &pb.OwnershipControlsRule{
		ObjectOwnership: r.ObjectOwnership,
	}
}

func protoToOwnershipControlsRule(p *pb.OwnershipControlsRule) *OwnershipControlsRule {
	if p == nil {
		return nil
	}
	return &OwnershipControlsRule{
		ObjectOwnership: p.ObjectOwnership,
	}
}

func requestPaymentConfigurationToProto(c *RequestPaymentConfiguration) *pb.RequestPaymentConfiguration {
	if c == nil {
		return nil
	}
	return &pb.RequestPaymentConfiguration{
		Payer: c.Payer,
	}
}

func protoToRequestPaymentConfiguration(p *pb.RequestPaymentConfiguration) *RequestPaymentConfiguration {
	if p == nil {
		return nil
	}
	return &RequestPaymentConfiguration{
		Payer: p.Payer,
	}
}

func accelerateConfigurationToProto(c *AccelerateConfiguration) *pb.AccelerateConfiguration {
	if c == nil {
		return nil
	}
	return &pb.AccelerateConfiguration{
		Status: c.Status,
	}
}

func protoToAccelerateConfiguration(p *pb.AccelerateConfiguration) *AccelerateConfiguration {
	if p == nil {
		return nil
	}
	return &AccelerateConfiguration{
		Status: p.Status,
	}
}

func replicationConfigurationToProto(c *ReplicationConfiguration) *pb.ReplicationConfiguration {
	if c == nil {
		return nil
	}
	return &pb.ReplicationConfiguration{
		Role:  c.Role,
		Rules: replicationRulesToProto(c.Rules),
	}
}

func protoToReplicationConfiguration(p *pb.ReplicationConfiguration) *ReplicationConfiguration {
	if p == nil {
		return nil
	}
	return &ReplicationConfiguration{
		Role:  p.Role,
		Rules: protoToReplicationRules(p.Rules),
	}
}

func replicationRulesToProto(rules []ReplicationRule) []*pb.ReplicationRule {
	if rules == nil {
		return nil
	}
	result := make([]*pb.ReplicationRule, len(rules))
	for i, r := range rules {
		result[i] = replicationRuleToProto(&r)
	}
	return result
}

func protoToReplicationRules(p []*pb.ReplicationRule) []ReplicationRule {
	if p == nil {
		return nil
	}
	result := make([]ReplicationRule, len(p))
	for i, r := range p {
		result[i] = *protoToReplicationRule(r)
	}
	return result
}

func replicationRuleToProto(r *ReplicationRule) *pb.ReplicationRule {
	if r == nil {
		return nil
	}
	return &pb.ReplicationRule{
		Id:                      r.ID,
		Priority:                r.Priority,
		Status:                  r.Status,
		Filter:                  replicationFilterToProto(r.Filter),
		Destination:             replicationDestinationToProto(r.Destination),
		DeleteMarkerReplication: r.DeleteMarkerReplication,
	}
}

func protoToReplicationRule(p *pb.ReplicationRule) *ReplicationRule {
	if p == nil {
		return nil
	}
	return &ReplicationRule{
		ID:                      p.Id,
		Priority:                p.Priority,
		Status:                  p.Status,
		Filter:                  protoToReplicationFilter(p.Filter),
		Destination:             protoToReplicationDestination(p.Destination),
		DeleteMarkerReplication: p.DeleteMarkerReplication,
	}
}

func replicationFilterToProto(f *ReplicationFilter) *pb.ReplicationFilter {
	if f == nil {
		return nil
	}
	return &pb.ReplicationFilter{
		Prefix:      f.Prefix,
		Tag:         replicationTagFilterToProto(f.Tag),
		AndOperator: replicationAndOperatorToProto(f.AndOperator),
	}
}

func protoToReplicationFilter(p *pb.ReplicationFilter) *ReplicationFilter {
	if p == nil {
		return nil
	}
	return &ReplicationFilter{
		Prefix:      p.Prefix,
		Tag:         protoToReplicationTagFilter(p.Tag),
		AndOperator: protoToReplicationAndOperator(p.AndOperator),
	}
}

func replicationTagFilterToProto(t *ReplicationTagFilter) *pb.ReplicationTagFilter {
	if t == nil {
		return nil
	}
	return &pb.ReplicationTagFilter{
		Key:   t.Key,
		Value: t.Value,
	}
}

func protoToReplicationTagFilter(p *pb.ReplicationTagFilter) *ReplicationTagFilter {
	if p == nil {
		return nil
	}
	return &ReplicationTagFilter{
		Key:   p.Key,
		Value: p.Value,
	}
}

func replicationAndOperatorToProto(a *ReplicationAndOperator) *pb.ReplicationAndOperator {
	if a == nil {
		return nil
	}
	return &pb.ReplicationAndOperator{
		Prefix: a.Prefix,
		Tags:   replicationTagFiltersToProto(a.Tags),
	}
}

func protoToReplicationAndOperator(p *pb.ReplicationAndOperator) *ReplicationAndOperator {
	if p == nil {
		return nil
	}
	return &ReplicationAndOperator{
		Prefix: p.Prefix,
		Tags:   protoToReplicationTagFilters(p.Tags),
	}
}

func replicationTagFiltersToProto(tags []ReplicationTagFilter) []*pb.ReplicationTagFilter {
	if tags == nil {
		return nil
	}
	result := make([]*pb.ReplicationTagFilter, len(tags))
	for i, t := range tags {
		result[i] = replicationTagFilterToProto(&t)
	}
	return result
}

func protoToReplicationTagFilters(p []*pb.ReplicationTagFilter) []ReplicationTagFilter {
	if p == nil {
		return nil
	}
	result := make([]ReplicationTagFilter, len(p))
	for i, t := range p {
		result[i] = *protoToReplicationTagFilter(t)
	}
	return result
}

func replicationDestinationToProto(d *ReplicationDestination) *pb.ReplicationDestination {
	if d == nil {
		return nil
	}
	return &pb.ReplicationDestination{
		Bucket:                   d.Bucket,
		StorageClass:             d.StorageClass,
		Account:                  d.Account,
		AccessControlTranslation: d.AccessControlTranslation,
	}
}

func protoToReplicationDestination(p *pb.ReplicationDestination) *ReplicationDestination {
	if p == nil {
		return nil
	}
	return &ReplicationDestination{
		Bucket:                   p.Bucket,
		StorageClass:             p.StorageClass,
		Account:                  p.Account,
		AccessControlTranslation: p.AccessControlTranslation,
	}
}
