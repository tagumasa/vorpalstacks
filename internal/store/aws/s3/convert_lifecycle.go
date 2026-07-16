package s3

import (
	pb "vorpalstacks/internal/pb/storage/storage_s3"
)

func lifecycleConfigurationToProto(c *LifecycleConfiguration) *pb.LifecycleConfiguration {
	if c == nil {
		return nil
	}
	return &pb.LifecycleConfiguration{
		Rules: lifecycleRulesToProto(c.Rules),
	}
}

func protoToLifecycleConfiguration(p *pb.LifecycleConfiguration) *LifecycleConfiguration {
	if p == nil {
		return nil
	}
	return &LifecycleConfiguration{
		Rules: protoToLifecycleRules(p.Rules),
	}
}

func lifecycleRulesToProto(rules []LifecycleRule) []*pb.LifecycleRule {
	if rules == nil {
		return nil
	}
	result := make([]*pb.LifecycleRule, len(rules))
	for i, r := range rules {
		result[i] = lifecycleRuleToProto(&r)
	}
	return result
}

func protoToLifecycleRules(p []*pb.LifecycleRule) []LifecycleRule {
	if p == nil {
		return nil
	}
	result := make([]LifecycleRule, len(p))
	for i, r := range p {
		result[i] = *protoToLifecycleRule(r)
	}
	return result
}

func lifecycleRuleToProto(r *LifecycleRule) *pb.LifecycleRule {
	if r == nil {
		return nil
	}
	return &pb.LifecycleRule{
		Id:                             r.ID,
		Status:                         r.Status,
		Filter:                         lifecycleRuleFilterToProto(r.Filter),
		Expiration:                     lifecycleExpirationToProto(r.Expiration),
		Transitions:                    lifecycleTransitionsToProto(r.Transitions),
		NoncurrentVersionExpiration:    noncurrentVersionExpirationToProto(r.NoncurrentVersionExpiration),
		NoncurrentVersionTransitions:   noncurrentVersionTransitionsToProto(r.NoncurrentVersionTransitions),
		AbortIncompleteMultipartUpload: abortIncompleteUploadToProto(r.AbortIncompleteMultipartUpload),
	}
}

func protoToLifecycleRule(p *pb.LifecycleRule) *LifecycleRule {
	if p == nil {
		return nil
	}
	return &LifecycleRule{
		ID:                             p.Id,
		Status:                         p.Status,
		Filter:                         protoToLifecycleRuleFilter(p.Filter),
		Expiration:                     protoToLifecycleExpiration(p.Expiration),
		Transitions:                    protoToLifecycleTransitions(p.Transitions),
		NoncurrentVersionExpiration:    protoToNoncurrentVersionExpiration(p.NoncurrentVersionExpiration),
		NoncurrentVersionTransitions:   protoToNoncurrentVersionTransitions(p.NoncurrentVersionTransitions),
		AbortIncompleteMultipartUpload: protoToAbortIncompleteUpload(p.AbortIncompleteMultipartUpload),
	}
}

func lifecycleRuleFilterToProto(f *LifecycleRuleFilter) *pb.LifecycleRuleFilter {
	if f == nil {
		return nil
	}
	return &pb.LifecycleRuleFilter{
		Prefix:                f.Prefix,
		ObjectSizeGreaterThan: int64PtrToInt64(f.ObjectSizeGreaterThan),
		ObjectSizeLessThan:    int64PtrToInt64(f.ObjectSizeLessThan),
		And:                   lifecycleRuleAndOperatorToProto(f.And),
		Tag:                   tagToProtoPtr(f.Tag),
	}
}

func protoToLifecycleRuleFilter(p *pb.LifecycleRuleFilter) *LifecycleRuleFilter {
	if p == nil {
		return nil
	}
	return &LifecycleRuleFilter{
		Prefix:                p.Prefix,
		ObjectSizeGreaterThan: int64ToInt64Ptr(p.ObjectSizeGreaterThan),
		ObjectSizeLessThan:    int64ToInt64Ptr(p.ObjectSizeLessThan),
		And:                   protoToLifecycleRuleAndOperator(p.And),
		Tag:                   protoToTagPtr(p.Tag),
	}
}

func lifecycleRuleAndOperatorToProto(o *LifecycleRuleAndOperator) *pb.LifecycleRuleAndOperator {
	if o == nil {
		return nil
	}
	return &pb.LifecycleRuleAndOperator{
		Prefix:                o.Prefix,
		Tags:                  tagsToProto(o.Tags),
		ObjectSizeGreaterThan: int64PtrToInt64(o.ObjectSizeGreaterThan),
		ObjectSizeLessThan:    int64PtrToInt64(o.ObjectSizeLessThan),
	}
}

func protoToLifecycleRuleAndOperator(p *pb.LifecycleRuleAndOperator) *LifecycleRuleAndOperator {
	if p == nil {
		return nil
	}
	return &LifecycleRuleAndOperator{
		Prefix:                p.Prefix,
		Tags:                  protoToTags(p.Tags),
		ObjectSizeGreaterThan: int64ToInt64Ptr(p.ObjectSizeGreaterThan),
		ObjectSizeLessThan:    int64ToInt64Ptr(p.ObjectSizeLessThan),
	}
}

func lifecycleExpirationToProto(e *LifecycleExpiration) *pb.LifecycleExpiration {
	if e == nil {
		return nil
	}
	return &pb.LifecycleExpiration{
		Date:                      timeToProto(e.Date),
		Days:                      intPtrToInt32(e.Days),
		ExpiredObjectDeleteMarker: boolToBool(e.ExpiredObjectDeleteMarker),
	}
}

func protoToLifecycleExpiration(p *pb.LifecycleExpiration) *LifecycleExpiration {
	if p == nil {
		return nil
	}
	return &LifecycleExpiration{
		Date:                      protoToTime(p.Date),
		Days:                      int32ToIntPtr(p.Days),
		ExpiredObjectDeleteMarker: boolToBoolPtr(p.ExpiredObjectDeleteMarker),
	}
}

func lifecycleTransitionsToProto(t []LifecycleTransition) []*pb.LifecycleTransition {
	if t == nil {
		return nil
	}
	result := make([]*pb.LifecycleTransition, len(t))
	for i, tr := range t {
		result[i] = lifecycleTransitionToProto(&tr)
	}
	return result
}

func protoToLifecycleTransitions(p []*pb.LifecycleTransition) []LifecycleTransition {
	if p == nil {
		return nil
	}
	result := make([]LifecycleTransition, len(p))
	for i, tr := range p {
		result[i] = *protoToLifecycleTransition(tr)
	}
	return result
}

func lifecycleTransitionToProto(t *LifecycleTransition) *pb.LifecycleTransition {
	if t == nil {
		return nil
	}
	return &pb.LifecycleTransition{
		Date:         timeToProto(t.Date),
		Days:         intPtrToInt32(t.Days),
		StorageClass: objectStorageClassToProto(t.StorageClass),
	}
}

func protoToLifecycleTransition(p *pb.LifecycleTransition) *LifecycleTransition {
	if p == nil {
		return nil
	}
	return &LifecycleTransition{
		Date:         protoToTime(p.Date),
		Days:         int32ToIntPtr(p.Days),
		StorageClass: protoToObjectStorageClass(p.StorageClass),
	}
}

func noncurrentVersionExpirationToProto(e *NoncurrentVersionExpiration) *pb.NoncurrentVersionExpiration {
	if e == nil {
		return nil
	}
	return &pb.NoncurrentVersionExpiration{
		NoncurrentDays:          intPtrToInt32(e.NoncurrentDays),
		NewerNoncurrentVersions: intPtrToInt32(e.NewerNoncurrentVersions),
	}
}

func protoToNoncurrentVersionExpiration(p *pb.NoncurrentVersionExpiration) *NoncurrentVersionExpiration {
	if p == nil {
		return nil
	}
	return &NoncurrentVersionExpiration{
		NoncurrentDays:          int32ToIntPtr(p.NoncurrentDays),
		NewerNoncurrentVersions: int32ToIntPtr(p.NewerNoncurrentVersions),
	}
}

func noncurrentVersionTransitionsToProto(t []NoncurrentVersionTransition) []*pb.NoncurrentVersionTransition {
	if t == nil {
		return nil
	}
	result := make([]*pb.NoncurrentVersionTransition, len(t))
	for i, tr := range t {
		result[i] = noncurrentVersionTransitionToProto(&tr)
	}
	return result
}

func protoToNoncurrentVersionTransitions(p []*pb.NoncurrentVersionTransition) []NoncurrentVersionTransition {
	if p == nil {
		return nil
	}
	result := make([]NoncurrentVersionTransition, len(p))
	for i, tr := range p {
		result[i] = *protoToNoncurrentVersionTransition(tr)
	}
	return result
}

func noncurrentVersionTransitionToProto(t *NoncurrentVersionTransition) *pb.NoncurrentVersionTransition {
	if t == nil {
		return nil
	}
	return &pb.NoncurrentVersionTransition{
		NoncurrentDays:          intPtrToInt32(t.NoncurrentDays),
		NewerNoncurrentVersions: intPtrToInt32(t.NewerNoncurrentVersions),
		StorageClass:            objectStorageClassToProto(t.StorageClass),
	}
}

func protoToNoncurrentVersionTransition(p *pb.NoncurrentVersionTransition) *NoncurrentVersionTransition {
	if p == nil {
		return nil
	}
	return &NoncurrentVersionTransition{
		NoncurrentDays:          int32ToIntPtr(p.NoncurrentDays),
		NewerNoncurrentVersions: int32ToIntPtr(p.NewerNoncurrentVersions),
		StorageClass:            protoToObjectStorageClass(p.StorageClass),
	}
}

func abortIncompleteUploadToProto(a *AbortIncompleteUpload) *pb.AbortIncompleteUpload {
	if a == nil {
		return nil
	}
	return &pb.AbortIncompleteUpload{
		DaysAfterInitiation: intPtrToInt32(a.DaysAfterInitiation),
	}
}

func protoToAbortIncompleteUpload(p *pb.AbortIncompleteUpload) *AbortIncompleteUpload {
	if p == nil {
		return nil
	}
	return &AbortIncompleteUpload{
		DaysAfterInitiation: int32ToIntPtr(p.DaysAfterInitiation),
	}
}
