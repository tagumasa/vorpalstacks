// Package kinesis provides AWS Kinesis stream service operations for vorpalstacks.
package kinesis

import (
	"fmt"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// KinesisService provides AWS Kinesis stream operations.
type KinesisService struct {
	accountID string
}

// NewKinesisService creates a new Kinesis service instance.
func NewKinesisService(store storage.BasicStorage, accountID, region string) *KinesisService {
	return &KinesisService{
		accountID: accountID,
	}
}

func (s *KinesisService) store(reqCtx *request.RequestContext) (*kinesisstore.KinesisStore, error) {
	basicStore, err := reqCtx.GetStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage: %w", err)
	}
	tstore, ok := basicStore.(storage.TransactionalStorageWith2PC)
	if !ok {
		return nil, fmt.Errorf("storage does not support TransactionalStorageWith2PC")
	}
	return kinesisstore.NewKinesisStore(tstore, reqCtx.GetAccountID(), reqCtx.GetRegion()), nil
}

// RegisterHandlers registers the Kinesis service handlers with the dispatcher.
func (s *KinesisService) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("kinesis", "CreateStream", s.CreateStream)
	d.RegisterHandlerForService("kinesis", "DeleteStream", s.DeleteStream)
	d.RegisterHandlerForService("kinesis", "DescribeStream", s.DescribeStream)
	d.RegisterHandlerForService("kinesis", "DescribeStreamSummary", s.DescribeStreamSummary)
	d.RegisterHandlerForService("kinesis", "ListStreams", s.ListStreams)
	d.RegisterHandlerForService("kinesis", "UpdateStreamMode", s.UpdateStreamMode)

	d.RegisterHandlerForService("kinesis", "PutRecord", s.PutRecord)
	d.RegisterHandlerForService("kinesis", "PutRecords", s.PutRecords)
	d.RegisterHandlerForService("kinesis", "GetRecords", s.GetRecords)
	d.RegisterHandlerForService("kinesis", "GetShardIterator", s.GetShardIterator)

	d.RegisterHandlerForService("kinesis", "ListShards", s.ListShards)
	d.RegisterHandlerForService("kinesis", "SplitShard", s.SplitShard)
	d.RegisterHandlerForService("kinesis", "MergeShards", s.MergeShards)
	d.RegisterHandlerForService("kinesis", "UpdateShardCount", s.UpdateShardCount)

	d.RegisterHandlerForService("kinesis", "RegisterStreamConsumer", s.RegisterStreamConsumer)
	d.RegisterHandlerForService("kinesis", "DeregisterStreamConsumer", s.DeregisterStreamConsumer)
	d.RegisterHandlerForService("kinesis", "DescribeStreamConsumer", s.DescribeStreamConsumer)
	d.RegisterHandlerForService("kinesis", "ListStreamConsumers", s.ListStreamConsumers)
	d.RegisterHandlerForService("kinesis", "SubscribeToShard", s.SubscribeToShard)

	d.RegisterHandlerForService("kinesis", "AddTagsToStream", s.AddTagsToStream)
	d.RegisterHandlerForService("kinesis", "RemoveTagsFromStream", s.RemoveTagsFromStream)
	d.RegisterHandlerForService("kinesis", "ListTagsForStream", s.ListTagsForStream)
	d.RegisterHandlerForService("kinesis", "TagResource", s.TagResource)
	d.RegisterHandlerForService("kinesis", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("kinesis", "ListTagsForResource", s.ListTagsForResource)

	d.RegisterHandlerForService("kinesis", "IncreaseStreamRetentionPeriod", s.IncreaseStreamRetentionPeriod)
	d.RegisterHandlerForService("kinesis", "DecreaseStreamRetentionPeriod", s.DecreaseStreamRetentionPeriod)
	d.RegisterHandlerForService("kinesis", "DescribeLimits", s.DescribeLimits)
	d.RegisterHandlerForService("kinesis", "DescribeAccountSettings", s.DescribeAccountSettings)
	d.RegisterHandlerForService("kinesis", "UpdateAccountSettings", s.UpdateAccountSettings)

	d.RegisterHandlerForService("kinesis", "EnableEnhancedMonitoring", s.EnableEnhancedMonitoring)
	d.RegisterHandlerForService("kinesis", "DisableEnhancedMonitoring", s.DisableEnhancedMonitoring)

	d.RegisterHandlerForService("kinesis", "StartStreamEncryption", s.StartStreamEncryption)
	d.RegisterHandlerForService("kinesis", "StopStreamEncryption", s.StopStreamEncryption)

	d.RegisterHandlerForService("kinesis", "GetResourcePolicy", s.GetResourcePolicy)
	d.RegisterHandlerForService("kinesis", "PutResourcePolicy", s.PutResourcePolicy)
	d.RegisterHandlerForService("kinesis", "DeleteResourcePolicy", s.DeleteResourcePolicy)

	d.RegisterHandlerForService("kinesis", "UpdateMaxRecordSize", s.UpdateMaxRecordSize)
	d.RegisterHandlerForService("kinesis", "UpdateStreamWarmThroughput", s.UpdateStreamWarmThroughput)
}
