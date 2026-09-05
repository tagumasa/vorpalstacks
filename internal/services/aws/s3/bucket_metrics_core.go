package s3

import (
	types "vorpalstacks/internal/common/tags"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// metricsIdAllowed reports whether a character may appear in a metrics
// configuration id: AWS restricts ids to 64 characters of letters, numbers,
// periods, dashes, and underscores.
func metricsIdAllowed(id string) bool {
	if len(id) == 0 || len(id) > 64 {
		return false
	}
	for _, ch := range id {
		switch {
		case ch >= 'a' && ch <= 'z', ch >= 'A' && ch <= 'Z', ch >= '0' && ch <= '9':
		case ch == '.', ch == '-', ch == '_':
		default:
			return false
		}
	}
	return true
}

// putBucketMetricsConfigurationCore validates and stores one metrics
// configuration under the request's id, replacing any existing configuration
// with that id (a full replacement, as AWS documents).
func (s *S3Service) putBucketMetricsConfigurationCore(bucketStore s3store.BucketStoreInterface, in *PutBucketMetricsConfigurationInput) error {
	if in.Id == "" {
		return NewInvalidArgumentError("the id query parameter is required")
	}
	if !metricsIdAllowed(in.Id) {
		return NewInvalidArgumentError("the metrics configuration id must be at most 64 characters of letters, numbers, periods, dashes, and underscores")
	}
	config := in.MetricsConfiguration
	if config == nil || config.Id == "" {
		return ErrMalformedXML
	}

	stored := &s3store.MetricsConfiguration{ID: in.Id}
	if config.Filter != nil {
		filter := config.Filter
		predicates := 0
		if filter.Prefix != "" {
			predicates++
		}
		if filter.Tag != nil {
			predicates++
		}
		if filter.AccessPointArn != "" {
			predicates++
		}
		if filter.And != nil {
			predicates++
		}
		// The model's MetricsFilter is a union, so a present filter carries
		// exactly one predicate; an empty element is as malformed as a
		// multi-predicate one. A missing filter selects the whole bucket.
		if predicates != 1 {
			return ErrMalformedXML
		}
		storedFilter := &s3store.MetricsFilter{
			Prefix:         filter.Prefix,
			AccessPointArn: filter.AccessPointArn,
		}
		if filter.Tag != nil {
			tag := filter.Tag.ToCommon()
			storedFilter.Tag = &tag
		}
		if filter.And != nil {
			and := &s3store.MetricsAndOperator{
				Prefix:         filter.And.Prefix,
				AccessPointArn: filter.And.AccessPointArn,
			}
			for _, wireTag := range filter.And.Tags {
				and.Tags = append(and.Tags, wireTag.ToCommon())
			}
			andPredicates := 0
			if and.Prefix != "" {
				andPredicates++
			}
			if len(and.Tags) > 0 {
				andPredicates++
			}
			if and.AccessPointArn != "" {
				andPredicates++
			}
			if andPredicates < 2 {
				return ErrMalformedXML
			}
			storedFilter.And = and
		}
		stored.Filter = storedFilter
	}

	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return err
	}
	if _, exists := bucket.MetricsConfigurations[in.Id]; !exists &&
		len(bucket.MetricsConfigurations) >= s3store.MaxBucketConfigurations {
		return ErrTooManyConfigurations
	}
	if err := bucketStore.SetMetricsConfiguration(in.Bucket, in.Id, stored); err != nil {
		return err
	}
	if s.requestMetrics != nil {
		s.requestMetrics.invalidateConfigs(in.Bucket)
	}
	return nil
}

// getBucketMetricsConfigurationCore reads one metrics configuration.
func (s *S3Service) getBucketMetricsConfigurationCore(bucketStore s3store.BucketStoreInterface, in *GetBucketMetricsConfigurationInput) (*GetBucketMetricsConfigurationOutput, error) {
	stored, err := bucketStore.GetMetricsConfiguration(in.Bucket, in.Id)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, ErrNoSuchConfiguration
	}
	return &GetBucketMetricsConfigurationOutput{MetricsConfiguration: metricsConfigurationToOutput(stored)}, nil
}

// deleteBucketMetricsConfigurationCore removes one metrics configuration; an
// absent id deletes nothing, matching the idempotent 204 the AWS API
// documents.
func (s *S3Service) deleteBucketMetricsConfigurationCore(bucketStore s3store.BucketStoreInterface, in *DeleteBucketMetricsConfigurationInput) error {
	if in.Id == "" {
		return NewInvalidArgumentError("the id query parameter is required")
	}
	if err := bucketStore.DeleteMetricsConfiguration(in.Bucket, in.Id); err != nil {
		return err
	}
	if s.requestMetrics != nil {
		s.requestMetrics.invalidateConfigs(in.Bucket)
	}
	return nil
}

// listBucketMetricsConfigurationsCore returns one page of at most 100
// configurations ordered by id; the continuation token is the id of the last
// configuration returned.
func (s *S3Service) listBucketMetricsConfigurationsCore(bucketStore s3store.BucketStoreInterface, in *ListBucketMetricsConfigurationsInput) (*ListBucketMetricsConfigurationsOutput, error) {
	configs, err := bucketStore.ListMetricsConfigurations(in.Bucket)
	if err != nil {
		return nil, err
	}
	const pageSize = 100
	start := 0
	if in.ContinuationToken != "" {
		for start < len(configs) && configs[start].ID <= in.ContinuationToken {
			start++
		}
	}
	end := start + pageSize
	if end > len(configs) {
		end = len(configs)
	}
	out := &ListBucketMetricsConfigurationsOutput{
		ContinuationToken: in.ContinuationToken,
	}
	for _, stored := range configs[start:end] {
		out.MetricsConfigurations = append(out.MetricsConfigurations, *metricsConfigurationToOutput(stored))
	}
	if end < len(configs) {
		out.IsTruncated = true
		out.NextContinuationToken = configs[end-1].ID
	}
	return out, nil
}

func metricsConfigurationToOutput(stored *s3store.MetricsConfiguration) *MetricsConfigurationOutput {
	out := &MetricsConfigurationOutput{Id: stored.ID}
	if stored.Filter != nil {
		filter := &MetricsFilterOutput{
			Prefix:         stored.Filter.Prefix,
			AccessPointArn: stored.Filter.AccessPointArn,
		}
		if stored.Filter.Tag != nil {
			commonTag := CommonToTags([]types.Tag{*stored.Filter.Tag})
			filter.Tag = &commonTag[0]
		}
		if stored.Filter.And != nil {
			and := &MetricsAndOperatorOutput{
				Prefix:         stored.Filter.And.Prefix,
				AccessPointArn: stored.Filter.And.AccessPointArn,
			}
			and.Tags = CommonToTags(stored.Filter.And.Tags)
			filter.And = and
		}
		out.Filter = filter
	}
	return out
}
