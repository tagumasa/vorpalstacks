package cloudwatch

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func dashboardBucketName(region string) string {
	return "cw_dashboards-" + region
}

// DashboardStore provides storage operations for CloudWatch dashboards.
type DashboardStore struct {
	*common.BaseStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	mu         sync.Mutex
}

// NewDashboardStore creates a new DashboardStore instance.
func NewDashboardStore(store storage.BasicStorage, accountID, region string) *DashboardStore {
	return &DashboardStore{
		BaseStore:  common.NewBaseStore(store.Bucket(dashboardBucketName(region)), "cloudwatch-dashboards"),
		arnBuilder: svcarn.NewARNBuilder(accountID, region),
		accountID:  accountID,
		region:     region,
	}
}

func (s *DashboardStore) buildDashboardArn(name string) string {
	return s.arnBuilder.Build("cloudwatch", "dashboard:"+name)
}

// PutDashboard creates or updates a CloudWatch dashboard with the given name and body.
func (s *DashboardStore) PutDashboard(name, body string) (*Dashboard, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	key := dashboardKey(name)

	existing := &Dashboard{}
	if err := s.BaseStore.Get(key, existing); err == nil {
		existing.Body = body
		existing.UpdatedAt = now
		if err := s.BaseStore.Put(key, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	dashboard := &Dashboard{
		Name:      name,
		ARN:       s.buildDashboardArn(name),
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.BaseStore.Put(key, dashboard); err != nil {
		return nil, err
	}
	return dashboard, nil
}

// GetDashboard retrieves a CloudWatch dashboard by name.
func (s *DashboardStore) GetDashboard(name string) (*Dashboard, error) {
	var dashboard Dashboard
	key := dashboardKey(name)
	if err := s.BaseStore.Get(key, &dashboard); err != nil {
		return nil, fmt.Errorf("dashboard not found: %s", name)
	}
	return &dashboard, nil
}

// DeleteDashboards deletes the specified dashboards, returning lists of successfully deleted and not-found names.
func (s *DashboardStore) DeleteDashboards(names []string) ([]string, []string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var deleted, notFound []string
	for _, name := range names {
		key := dashboardKey(name)
		if !s.BaseStore.Exists(key) {
			notFound = append(notFound, name)
			continue
		}
		if err := s.BaseStore.Delete(key); err != nil {
			return deleted, notFound, err
		}
		deleted = append(deleted, name)
	}
	return deleted, notFound, nil
}

// ListDashboards returns all dashboards whose names match the given prefix.
func (s *DashboardStore) ListDashboards(prefix string) ([]*Dashboard, error) {
	var dashboards []*Dashboard
	prefixKey := dashboardKey(prefix)
	err := s.BaseStore.ScanPrefix(prefixKey, func(key string, value []byte) error {
		var dashboard Dashboard
		if err := json.Unmarshal(value, &dashboard); err != nil {
			return err
		}
		dashboards = append(dashboards, &dashboard)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dashboards, nil
}

func dashboardKey(name string) string {
	return "dashboard:" + name
}
