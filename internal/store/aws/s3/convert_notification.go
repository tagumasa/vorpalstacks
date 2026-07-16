package s3

import (
	pb "vorpalstacks/internal/pb/storage/storage_s3"
)

func notificationConfigurationToProto(c *NotificationConfiguration) *pb.NotificationConfiguration {
	if c == nil {
		return nil
	}
	return &pb.NotificationConfiguration{
		TopicConfigurations:  topicNotificationConfigsToProto(c.TopicConfigurations),
		QueueConfigurations:  queueNotificationConfigsToProto(c.QueueConfigurations),
		LambdaConfigurations: lambdaNotificationConfigsToProto(c.LambdaConfigurations),
	}
}

func protoToNotificationConfiguration(p *pb.NotificationConfiguration) *NotificationConfiguration {
	if p == nil {
		return nil
	}
	return &NotificationConfiguration{
		TopicConfigurations:  protoToTopicNotificationConfigs(p.TopicConfigurations),
		QueueConfigurations:  protoToQueueNotificationConfigs(p.QueueConfigurations),
		LambdaConfigurations: protoToLambdaNotificationConfigs(p.LambdaConfigurations),
	}
}

func topicNotificationConfigsToProto(configs []TopicNotificationConfiguration) []*pb.TopicNotificationConfiguration {
	if configs == nil {
		return nil
	}
	result := make([]*pb.TopicNotificationConfiguration, len(configs))
	for i, c := range configs {
		result[i] = topicNotificationConfigToProto(&c)
	}
	return result
}

func protoToTopicNotificationConfigs(p []*pb.TopicNotificationConfiguration) []TopicNotificationConfiguration {
	if p == nil {
		return nil
	}
	result := make([]TopicNotificationConfiguration, len(p))
	for i, c := range p {
		result[i] = *protoToTopicNotificationConfig(c)
	}
	return result
}

func topicNotificationConfigToProto(c *TopicNotificationConfiguration) *pb.TopicNotificationConfiguration {
	if c == nil {
		return nil
	}
	return &pb.TopicNotificationConfiguration{
		Id:       c.Id,
		TopicArn: c.TopicArn,
		Events:   c.Events,
		Filter:   notificationConfigFilterToProto(c.Filter),
	}
}

func protoToTopicNotificationConfig(p *pb.TopicNotificationConfiguration) *TopicNotificationConfiguration {
	if p == nil {
		return nil
	}
	return &TopicNotificationConfiguration{
		Id:       p.Id,
		TopicArn: p.TopicArn,
		Events:   p.Events,
		Filter:   protoToNotificationConfigFilter(p.Filter),
	}
}

func queueNotificationConfigsToProto(configs []QueueNotificationConfiguration) []*pb.QueueNotificationConfiguration {
	if configs == nil {
		return nil
	}
	result := make([]*pb.QueueNotificationConfiguration, len(configs))
	for i, c := range configs {
		result[i] = queueNotificationConfigToProto(&c)
	}
	return result
}

func protoToQueueNotificationConfigs(p []*pb.QueueNotificationConfiguration) []QueueNotificationConfiguration {
	if p == nil {
		return nil
	}
	result := make([]QueueNotificationConfiguration, len(p))
	for i, c := range p {
		result[i] = *protoToQueueNotificationConfig(c)
	}
	return result
}

func queueNotificationConfigToProto(c *QueueNotificationConfiguration) *pb.QueueNotificationConfiguration {
	if c == nil {
		return nil
	}
	return &pb.QueueNotificationConfiguration{
		Id:       c.Id,
		QueueArn: c.QueueArn,
		Events:   c.Events,
		Filter:   notificationConfigFilterToProto(c.Filter),
	}
}

func protoToQueueNotificationConfig(p *pb.QueueNotificationConfiguration) *QueueNotificationConfiguration {
	if p == nil {
		return nil
	}
	return &QueueNotificationConfiguration{
		Id:       p.Id,
		QueueArn: p.QueueArn,
		Events:   p.Events,
		Filter:   protoToNotificationConfigFilter(p.Filter),
	}
}

func lambdaNotificationConfigsToProto(configs []LambdaNotificationConfiguration) []*pb.LambdaNotificationConfiguration {
	if configs == nil {
		return nil
	}
	result := make([]*pb.LambdaNotificationConfiguration, len(configs))
	for i, c := range configs {
		result[i] = lambdaNotificationConfigToProto(&c)
	}
	return result
}

func protoToLambdaNotificationConfigs(p []*pb.LambdaNotificationConfiguration) []LambdaNotificationConfiguration {
	if p == nil {
		return nil
	}
	result := make([]LambdaNotificationConfiguration, len(p))
	for i, c := range p {
		result[i] = *protoToLambdaNotificationConfig(c)
	}
	return result
}

func lambdaNotificationConfigToProto(c *LambdaNotificationConfiguration) *pb.LambdaNotificationConfiguration {
	if c == nil {
		return nil
	}
	return &pb.LambdaNotificationConfiguration{
		Id:                c.Id,
		LambdaFunctionArn: c.LambdaFunctionArn,
		Events:            c.Events,
		Filter:            notificationConfigFilterToProto(c.Filter),
	}
}

func protoToLambdaNotificationConfig(p *pb.LambdaNotificationConfiguration) *LambdaNotificationConfiguration {
	if p == nil {
		return nil
	}
	return &LambdaNotificationConfiguration{
		Id:                p.Id,
		LambdaFunctionArn: p.LambdaFunctionArn,
		Events:            p.Events,
		Filter:            protoToNotificationConfigFilter(p.Filter),
	}
}

func notificationConfigFilterToProto(f *NotificationConfigurationFilter) *pb.NotificationConfigurationFilter {
	if f == nil {
		return nil
	}
	return &pb.NotificationConfigurationFilter{
		Key: s3KeyFilterToProto(f.Key),
	}
}

func protoToNotificationConfigFilter(p *pb.NotificationConfigurationFilter) *NotificationConfigurationFilter {
	if p == nil {
		return nil
	}
	return &NotificationConfigurationFilter{
		Key: protoToS3KeyFilter(p.Key),
	}
}

func s3KeyFilterToProto(f *S3KeyFilter) *pb.S3KeyFilter {
	if f == nil {
		return nil
	}
	return &pb.S3KeyFilter{
		FilterRules: filterRulesToProto(f.FilterRules),
	}
}

func protoToS3KeyFilter(p *pb.S3KeyFilter) *S3KeyFilter {
	if p == nil {
		return nil
	}
	return &S3KeyFilter{
		FilterRules: protoToFilterRules(p.FilterRules),
	}
}

func filterRulesToProto(rules []FilterRule) []*pb.FilterRule {
	if rules == nil {
		return nil
	}
	result := make([]*pb.FilterRule, len(rules))
	for i, r := range rules {
		result[i] = filterRuleToProto(&r)
	}
	return result
}

func protoToFilterRules(p []*pb.FilterRule) []FilterRule {
	if p == nil {
		return nil
	}
	result := make([]FilterRule, len(p))
	for i, r := range p {
		result[i] = *protoToFilterRule(r)
	}
	return result
}

func filterRuleToProto(r *FilterRule) *pb.FilterRule {
	if r == nil {
		return nil
	}
	return &pb.FilterRule{
		Name:  r.Name,
		Value: r.Value,
	}
}

func protoToFilterRule(p *pb.FilterRule) *FilterRule {
	if p == nil {
		return nil
	}
	return &FilterRule{
		Name:  p.Name,
		Value: p.Value,
	}
}
