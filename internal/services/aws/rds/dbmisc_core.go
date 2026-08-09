package rds

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/core/logs"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/rds"
	storerds "vorpalstacks/internal/store/aws/rds"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"
)

// ---------------------------------------------------------------------------
// Input DTOs
// ---------------------------------------------------------------------------

// paginateRDSItems applies Marker/MaxRecords pagination to a sorted slice.
// Delegates to the shared pagination.PaginateSlice so that marker-not-found
// returns an empty page (not the first page) — the same behaviour Neptune's
// paginateItems relies on. maxRecords is clamped to 20-100 with a default
// of 100 when zero.
func paginateRDSItems[T any](items []T, marker string, maxRecords int32, idFunc func(T) string) ([]T, string) {
	if maxRecords == 0 {
		maxRecords = 100
	}
	if maxRecords < 20 {
		maxRecords = 20
	}
	if maxRecords > 100 {
		maxRecords = 100
	}
	result := pagination.PaginateSlice(items, marker, int(maxRecords), pagination.KeyExtractor[T](idFunc))
	return result.Items, result.NextMarker
}

type DescribeDBSubnetGroupsInput struct {
	DBSubnetGroupName string
	Filters           []*pb.Filter
	Marker            string
	MaxRecords        int32
}

type DescribeGlobalClustersInput struct {
	GlobalClusterIdentifier string
	Filters                 []*pb.Filter
	Marker                  string
	MaxRecords              int32
}

type DescribeEventSubscriptionsInput struct {
	SubscriptionName string
	Filters          []*pb.Filter
	Marker           string
	MaxRecords       int32
}

type DescribeEventsInput struct {
	SourceType       pb.SourceType
	SourceIdentifier string
	StartTime        string
	EndTime          string
	Duration         int32
	EventCategories  []string
	Marker           string
	MaxRecords       int32
	Filters          []*pb.Filter
}

type ListTagsForResourceInput struct {
	ResourceName string
}

type AddTagsToResourceInput struct {
	ResourceName string
	Tags         []*pb.Tag
}

type RemoveTagsFromResourceInput struct {
	ResourceName string
	TagKeys      []string
}

// ---------------------------------------------------------------------------
// Core methods
// ---------------------------------------------------------------------------

func (s *RDSService) describeDBSubnetGroupsCore(stores *rdsStores, in DescribeDBSubnetGroupsInput) (*pb.DBSubnetGroupMessage, error) {
	groups, err := stores.store.ListSubnetGroups()
	if err != nil {
		return nil, translateStoreError(err)
	}

	pbGroups := make([]*pb.DBSubnetGroup, 0, len(groups))
	for _, g := range groups {
		if in.DBSubnetGroupName != "" && g.DBSubnetGroupName != in.DBSubnetGroupName {
			continue
		}
		if !applyRDSFilters(in.Filters, subnetGroupFilterGetter(g)) {
			continue
		}
		pbGroups = append(pbGroups, subnetGroupToPb(g))
	}

	pbGroups, nextMarker := paginateRDSItems(pbGroups, in.Marker, in.MaxRecords, func(g *pb.DBSubnetGroup) string {
		return g.Dbsubnetgroupname
	})
	return &pb.DBSubnetGroupMessage{Dbsubnetgroups: pbGroups, Marker: nextMarker}, nil
}

func (s *RDSService) describeGlobalClustersCore(stores *rdsStores, in DescribeGlobalClustersInput) (*pb.GlobalClustersMessage, error) {
	clusters, err := stores.store.ListGlobalClusters()
	if err != nil {
		return nil, translateStoreError(err)
	}

	pbClusters := make([]*pb.GlobalCluster, 0, len(clusters))
	for _, c := range clusters {
		if in.GlobalClusterIdentifier != "" && c.GlobalClusterIdentifier != in.GlobalClusterIdentifier {
			continue
		}
		if !applyRDSFilters(in.Filters, globalClusterFilterGetter(c)) {
			continue
		}
		pbClusters = append(pbClusters, globalClusterToPb(c))
	}

	pbClusters, nextMarker := paginateRDSItems(pbClusters, in.Marker, in.MaxRecords, func(c *pb.GlobalCluster) string {
		return c.Globalclusteridentifier
	})
	return &pb.GlobalClustersMessage{Globalclusters: pbClusters, Marker: nextMarker}, nil
}

func (s *RDSService) describeEventSubscriptionsCore(stores *rdsStores, in DescribeEventSubscriptionsInput) (*pb.EventSubscriptionsMessage, error) {
	subs, err := stores.store.ListEventSubscriptions()
	if err != nil {
		return nil, translateStoreError(err)
	}

	pbSubs := make([]*pb.EventSubscription, 0, len(subs))
	for _, sub := range subs {
		if in.SubscriptionName != "" && sub.CustSubscriptionId != in.SubscriptionName {
			continue
		}
		if !applyRDSFilters(in.Filters, eventSubscriptionFilterGetter(sub)) {
			continue
		}
		pbSubs = append(pbSubs, eventSubscriptionToPb(sub))
	}

	pbSubs, nextMarker := paginateRDSItems(pbSubs, in.Marker, in.MaxRecords, func(sub *pb.EventSubscription) string {
		return sub.Custsubscriptionid
	})
	return &pb.EventSubscriptionsMessage{Eventsubscriptionslist: pbSubs, Marker: nextMarker}, nil
}

func (s *RDSService) describeEventsCore(stores *rdsStores, in DescribeEventsInput) (*pb.EventsMessage, error) {
	var sourceTypeStr string
	for _, f := range in.Filters {
		if strings.EqualFold(f.Name, "source-type") && len(f.Values) > 0 {
			sourceTypeStr = f.Values[0]
		} else {
			logs.Warn("rds-admin: DescribeEvents ignores unsupported Filter",
				logs.String("filter", f.Name))
		}
	}
	if sourceTypeStr == "" {
		var err error
		sourceTypeStr, err = sourceTypeToString(in.SourceType)
		if err != nil {
			return nil, newValidationError("%v", err)
		}
	}

	var startTime time.Time
	if st := in.StartTime; st != "" {
		var err error
		startTime, err = time.Parse(time.RFC3339, st)
		if err != nil {
			return nil, newValidationError("invalid StartTime: %v", err)
		}
	}
	var endTime time.Time
	if et := in.EndTime; et != "" {
		var err error
		endTime, err = time.Parse(time.RFC3339, et)
		if err != nil {
			return nil, newValidationError("invalid EndTime: %v", err)
		}
	}
	if in.Duration > 0 && !startTime.IsZero() && endTime.IsZero() {
		endTime = startTime.Add(time.Duration(in.Duration) * time.Minute)
	}

	maxRecords := int(in.MaxRecords)
	if maxRecords == 0 {
		maxRecords = 100
	} else if maxRecords < 20 {
		maxRecords = 20
	} else if maxRecords > 100 {
		maxRecords = 100
	}

	opts := storerds.EventListOptions{
		SourceType:       sourceTypeStr,
		SourceIdentifier: in.SourceIdentifier,
		StartTime:        startTime,
		EndTime:          endTime,
		EventCategories:  in.EventCategories,
		Marker:           in.Marker,
		MaxRecords:       maxRecords,
	}

	result, err := stores.store.ListEvents(opts)
	if err != nil {
		if errors.Is(err, storerds.ErrInvalidEventMarker) {
			return nil, newValidationError("invalid Marker: %s", in.Marker)
		}
		return nil, translateStoreError(err)
	}

	events := make([]*pb.Event, 0, len(result.Events))
	for _, evt := range result.Events {
		st, stErr := stringToSourceType(evt.SourceType)
		if stErr != nil {
			logs.Warn("rds-admin: skipping event with unknown SourceType",
				logs.String("event_id", evt.EventID),
				logs.String("source_type", evt.SourceType),
				logs.Err(stErr))
			continue
		}
		events = append(events, &pb.Event{
			Date:             evt.Date.UTC().Format(timeutils.ISO8601UTCFormat),
			Message:          evt.Message,
			Sourcearn:        evt.SourceArn,
			Sourceidentifier: evt.SourceIdentifier,
			Sourcetype:       st,
			Eventcategories:  evt.EventCategories,
		})
	}

	resp := &pb.EventsMessage{Events: events}
	if result.IsTruncated && result.Marker != "" {
		resp.Marker = result.Marker
	}
	return resp, nil
}

func (s *RDSService) listTagsForResourceCore(stores *rdsStores, in ListTagsForResourceInput) (*pb.TagListMessage, error) {
	tags, err := stores.store.GetTags(in.ResourceName)
	if err != nil {
		return nil, translateStoreError(err)
	}

	pbTags := make([]*pb.Tag, len(tags))
	for i, t := range tags {
		pbTags[i] = &pb.Tag{Key: t.Key, Value: t.Value}
	}

	return &pb.TagListMessage{Taglist: pbTags}, nil
}

func (s *RDSService) addTagsToResourceCore(stores *rdsStores, in AddTagsToResourceInput) (*pbcommon.Empty, error) {
	tags := make([]types.Tag, len(in.Tags))
	for i, t := range in.Tags {
		tags[i] = types.Tag{Key: t.Key, Value: t.Value}
	}

	if err := stores.store.AddTags(in.ResourceName, tags); err != nil {
		return nil, translateStoreError(err)
	}

	return &pbcommon.Empty{}, nil
}

func (s *RDSService) removeTagsFromResourceCore(stores *rdsStores, in RemoveTagsFromResourceInput) (*pbcommon.Empty, error) {
	if err := stores.store.RemoveTags(in.ResourceName, in.TagKeys); err != nil {
		return nil, translateStoreError(err)
	}

	return &pbcommon.Empty{}, nil
}

// ---------------------------------------------------------------------------
// Conversion helpers (store -> protobuf)
// ---------------------------------------------------------------------------

func subnetGroupToPb(g *storerds.DBSubnetGroup) *pb.DBSubnetGroup {
	p := &pb.DBSubnetGroup{
		Dbsubnetgroupname:        g.DBSubnetGroupName,
		Dbsubnetgroupdescription: g.DBSubnetGroupDescription,
		Vpcid:                    g.VpcId,
		Subnetgroupstatus:        g.SubnetGroupStatus,
		Dbsubnetgrouparn:         g.ARN,
	}
	for _, s := range g.Subnets {
		p.Subnets = append(p.Subnets, &pb.Subnet{
			Subnetidentifier:       s.SubnetIdentifier,
			Subnetavailabilityzone: &pb.AvailabilityZone{Name: s.SubnetAvailabilityZone},
			Subnetstatus:           s.SubnetStatus,
		})
	}
	return p
}

func globalClusterToPb(c *storerds.GlobalCluster) *pb.GlobalCluster {
	p := &pb.GlobalCluster{
		Globalclusteridentifier: c.GlobalClusterIdentifier,
		Globalclusterresourceid: c.GlobalClusterResourceId,
		Globalclusterarn:        c.GlobalClusterArn,
		Engine:                  c.Engine,
		Engineversion:           c.EngineVersion,
		Status:                  c.Status,
		Storageencrypted:        proto.Bool(c.StorageEncrypted),
		Deletionprotection:      proto.Bool(c.DeletionProtection),
	}
	for _, m := range c.GlobalClusterMembers {
		p.Globalclustermembers = append(p.Globalclustermembers, &pb.GlobalClusterMember{
			Dbclusterarn: m.DBClusterArn,
			Iswriter:     proto.Bool(m.IsWriter),
			Readers:      m.Readers,
		})
	}
	return p
}

func eventSubscriptionToPb(sub *storerds.EventSubscription) *pb.EventSubscription {
	return &pb.EventSubscription{
		Custsubscriptionid:   sub.CustSubscriptionId,
		Snstopicarn:          sub.SnsTopicArn,
		Status:               sub.Status,
		Sourcetype:           sub.SourceType,
		Sourceidslist:        sub.SourceIdsList,
		Eventcategorieslist:  sub.EventCategoriesList,
		Enabled:              proto.Bool(sub.Enabled),
		Eventsubscriptionarn: sub.CustSubscriptionArn,
	}
}

// ---------------------------------------------------------------------------
// Filter getters
// ---------------------------------------------------------------------------

func subnetGroupFilterGetter(g *storerds.DBSubnetGroup) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "db-subnet-group-name":
			return g.DBSubnetGroupName, true
		default:
			return "", false
		}
	}
}

func globalClusterFilterGetter(c *storerds.GlobalCluster) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "global-cluster-id":
			return c.GlobalClusterIdentifier, true
		case "engine":
			return c.Engine, true
		default:
			return "", false
		}
	}
}

func eventSubscriptionFilterGetter(sub *storerds.EventSubscription) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "event-subscription-id":
			return sub.CustSubscriptionId, true
		default:
			return "", false
		}
	}
}

// ---------------------------------------------------------------------------
// SourceType helpers (proto enum <-> string)
// ---------------------------------------------------------------------------

func sourceTypeToString(st pb.SourceType) (string, error) {
	switch st {
	case pb.SourceType_SOURCE_TYPE_DB_PARAMETER_GROUP:
		return "db-parameter-group", nil
	case pb.SourceType_SOURCE_TYPE_DB_SHARD_GROUP:
		return "db-shard-group", nil
	case pb.SourceType_SOURCE_TYPE_CUSTOM_ENGINE_VERSION:
		return "custom-engine-version", nil
	case pb.SourceType_SOURCE_TYPE_DB_PROXY:
		return "db-proxy", nil
	case pb.SourceType_SOURCE_TYPE_DB_INSTANCE:
		return "db-instance", nil
	case pb.SourceType_SOURCE_TYPE_ZERO_ETL:
		return "zero-etl", nil
	case pb.SourceType_SOURCE_TYPE_DB_CLUSTER:
		return "db-cluster", nil
	case pb.SourceType_SOURCE_TYPE_DB_SECURITY_GROUP:
		return "db-security-group", nil
	case pb.SourceType_SOURCE_TYPE_DB_CLUSTER_SNAPSHOT:
		return "db-cluster-snapshot", nil
	case pb.SourceType_SOURCE_TYPE_DB_SNAPSHOT:
		return "db-snapshot", nil
	default:
		if st == 0 {
			return "", nil
		}
		return "", fmt.Errorf("unsupported SourceType: %d", st)
	}
}

func stringToSourceType(s string) (pb.SourceType, error) {
	switch s {
	case "blue-green-deployment":
		return pb.SourceType_SOURCE_TYPE_BLUE_GREEN_DEPLOYMENT, nil
	case "db-parameter-group":
		return pb.SourceType_SOURCE_TYPE_DB_PARAMETER_GROUP, nil
	case "db-shard-group":
		return pb.SourceType_SOURCE_TYPE_DB_SHARD_GROUP, nil
	case "custom-engine-version":
		return pb.SourceType_SOURCE_TYPE_CUSTOM_ENGINE_VERSION, nil
	case "db-proxy":
		return pb.SourceType_SOURCE_TYPE_DB_PROXY, nil
	case "db-instance":
		return pb.SourceType_SOURCE_TYPE_DB_INSTANCE, nil
	case "zero-etl":
		return pb.SourceType_SOURCE_TYPE_ZERO_ETL, nil
	case "db-cluster":
		return pb.SourceType_SOURCE_TYPE_DB_CLUSTER, nil
	case "db-security-group":
		return pb.SourceType_SOURCE_TYPE_DB_SECURITY_GROUP, nil
	case "db-cluster-snapshot":
		return pb.SourceType_SOURCE_TYPE_DB_CLUSTER_SNAPSHOT, nil
	case "db-snapshot":
		return pb.SourceType_SOURCE_TYPE_DB_SNAPSHOT, nil
	default:
		return 0, fmt.Errorf("unsupported SourceType string: %q", s)
	}
}
