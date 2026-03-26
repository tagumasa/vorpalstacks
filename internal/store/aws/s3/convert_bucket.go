package s3

import (
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
