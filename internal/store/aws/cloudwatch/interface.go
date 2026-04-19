package cloudwatch

// AlarmStoreInterface defines operations for managing CloudWatch alarms.
type AlarmStoreInterface interface {
	CreateAlarm(alarm *Alarm) (*Alarm, error)
	GetAlarm(name string) (*Alarm, error)
	UpdateAlarm(alarm *Alarm) error
	DeleteAlarm(name string) error
	ListAlarms(alarmNamePrefix string) ([]*Alarm, error)
	SetAlarmState(name, state, reason string) error
	SetAlarmActionsEnabled(name string, enabled bool) error
	AddAlarmHistory(entry *AlarmHistoryEntry) error
	ListAlarmHistory(alarmName string, historyItemType string) ([]*AlarmHistoryEntry, error)
	Tag(arn string, tags map[string]string) error
	List(resourceKey string) (map[string]string, error)
	Untag(resourceKey string, tagKeys []string) error
}

// MetricChunkStoreInterface defines operations for managing CloudWatch metric data.
type MetricChunkStoreInterface interface {
	PutMetricData(namespace string, metricData []MetricDatum) error
	GetMetricStatistics(query MetricQuery) ([]*MetricStatistics, error)
	ListMetrics(namespace, metricName string, dimensions []Dimension) ([]MetricDatum, error)
}

// DashboardStoreInterface defines operations for managing CloudWatch dashboards.
type DashboardStoreInterface interface {
	PutDashboard(name, body string) (*Dashboard, error)
	GetDashboard(name string) (*Dashboard, error)
	DeleteDashboards(names []string) ([]string, []string, error)
	ListDashboards(prefix string) ([]*Dashboard, error)
}
