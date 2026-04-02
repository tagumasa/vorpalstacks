package cloudwatchlogs

import (
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// CloudWatchLogsStoreInterface defines operations for managing CloudWatch Logs.
type CloudWatchLogsStoreInterface interface {
	ARNBuilder() *svcarn.ARNBuilder
	CreateLogGroup(lg *LogGroup) error
	GetLogGroup(name string) (*LogGroup, error)
	PutLogGroup(lg *LogGroup) error
	DeleteLogGroup(name string) error
	ListLogGroups(prefix, marker string, maxItems int) ([]*LogGroup, string, error)
	CreateLogStream(ls *LogStream) error
	GetLogStream(logGroupName, logStreamName string) (*LogStream, error)
	DeleteLogStream(logGroupName, logStreamName string) error
	ListLogStreams(logGroupName, prefix, marker string, maxItems int) ([]*LogStream, string, error)
	PutLogEvents(logGroupName, logStreamName string, events []LogEntry) (string, error)
	GetLogEvents(logGroupName, logStreamName string, startTime, endTime int64, limit int, startFromHead bool, nextToken string) ([]*OutputLogEvent, string, int, error)
	FilterLogEvents(logGroupName string, logStreamNames []string, startTime, endTime int64, filterPattern string, limit int, nextToken string) ([]*OutputLogEvent, map[string]bool, string, error)
	PutSubscriptionFilter(filter *SubscriptionFilter) error
	DeleteSubscriptionFilter(logGroupName, filterName string) error
	GetSubscriptionFilter(logGroupName, filterName string) (*SubscriptionFilter, error)
	ListSubscriptionFilters(logGroupName, filterNamePrefix string) ([]*SubscriptionFilter, error)
	PurgeExpiredChunks(logGroupName string, cutoffTime int64) (int64, error)
	PurgeAllExpiredChunks() error
	PutDestination(dest *Destination) error
	GetDestination(name string) (*Destination, error)
	DeleteDestination(name string) error
	PutDestinationPolicy(name, accessPolicy string) error
	ListDestinations(prefix string) ([]*Destination, error)
	Raw() *Store
}

var _ CloudWatchLogsStoreInterface = (*Store)(nil)

// Raw returns the underlying CloudWatch Logs store.
func (s *Store) Raw() *Store {
	return s
}
