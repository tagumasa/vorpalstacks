package rds

import (
	storerds "vorpalstacks/internal/store/aws/rds"
	"vorpalstacks/internal/utils/aws/types"
)

// QueryClusters queries, filters, and paginates DBClusters from the given
// store. Both the admin handler Core method and the Neptune HTTP handler
// call this function to avoid duplicated store interaction logic.
func QueryClusters(store storerds.StoreInterface, in DescribeDBClustersInput) ([]*storerds.DBCluster, string, error) {
	clusters, err := store.ListClusters()
	if err != nil {
		return nil, "", translateStoreError(err)
	}

	if in.DBClusterIdentifier != "" {
		c, getErr := store.GetCluster(in.DBClusterIdentifier)
		if getErr != nil {
			return nil, "", translateStoreError(getErr)
		}
		clusters = []*storerds.DBCluster{c}
	}

	filtered := make([]*storerds.DBCluster, 0, len(clusters))
	for _, c := range clusters {
		if !applyRDSFilters(in.Filters, clusterFilterGetter(c)) {
			continue
		}
		filtered = append(filtered, c)
	}

	paginated, nextMarker := paginateRDSItems(filtered, in.Marker, in.MaxRecords, func(c *storerds.DBCluster) string {
		return c.DBClusterIdentifier
	})
	return paginated, nextMarker, nil
}

// QueryInstances queries, filters, and paginates DBInstances from the given
// store.
func QueryInstances(store storerds.StoreInterface, in DescribeDBInstancesInput) ([]*storerds.DBInstance, string, error) {
	instances, err := store.ListInstances()
	if err != nil {
		return nil, "", translateStoreError(err)
	}

	if in.DBInstanceIdentifier != "" {
		i, getErr := store.GetInstance(in.DBInstanceIdentifier)
		if getErr != nil {
			return nil, "", translateStoreError(getErr)
		}
		instances = []*storerds.DBInstance{i}
	}

	filtered := make([]*storerds.DBInstance, 0, len(instances))
	for _, i := range instances {
		if !applyRDSFilters(in.Filters, instanceFilterGetter(i)) {
			continue
		}
		filtered = append(filtered, i)
	}

	paginated, nextMarker := paginateRDSItems(filtered, in.Marker, in.MaxRecords, func(i *storerds.DBInstance) string {
		return i.DBInstanceIdentifier
	})
	return paginated, nextMarker, nil
}

// QueryClusterSnapshots queries, filters, and paginates DBClusterSnapshots.
func QueryClusterSnapshots(store storerds.StoreInterface, in DescribeDBClusterSnapshotsInput) ([]*storerds.DBClusterSnapshot, string, error) {
	if in.DBClusterSnapshotIdentifier != "" {
		snap, err := store.GetSnapshot(in.DBClusterSnapshotIdentifier)
		if err != nil {
			return nil, "", translateStoreError(err)
		}
		paginated, nextMarker := paginateRDSItems([]*storerds.DBClusterSnapshot{snap}, in.Marker, in.MaxRecords, func(s *storerds.DBClusterSnapshot) string {
			return s.DBClusterSnapshotIdentifier
		})
		return paginated, nextMarker, nil
	}

	snapshots, err := store.ListSnapshots()
	if err != nil {
		return nil, "", translateStoreError(err)
	}

	filtered := make([]*storerds.DBClusterSnapshot, 0, len(snapshots))
	for _, snap := range snapshots {
		if in.DBClusterSnapshotIdentifier != "" && snap.DBClusterSnapshotIdentifier != in.DBClusterSnapshotIdentifier {
			continue
		}
		if in.DBClusterIdentifier != "" && snap.DBClusterIdentifier != in.DBClusterIdentifier {
			continue
		}
		if in.SnapshotType != "" && snap.SnapshotType != in.SnapshotType {
			continue
		}
		if !applyRDSFilters(in.Filters, clusterSnapshotFilterGetter(snap)) {
			continue
		}
		filtered = append(filtered, snap)
	}

	paginated, nextMarker := paginateRDSItems(filtered, in.Marker, in.MaxRecords, func(s *storerds.DBClusterSnapshot) string {
		return s.DBClusterSnapshotIdentifier
	})
	return paginated, nextMarker, nil
}

// QuerySnapshots queries, filters, and paginates DBSnapshots.
func QuerySnapshots(store storerds.StoreInterface, in DescribeDBSnapshotsInput) ([]*storerds.DBInstanceSnapshot, string, error) {
	if in.DBSnapshotIdentifier != "" {
		snap, err := store.GetInstanceSnapshot(in.DBSnapshotIdentifier)
		if err != nil {
			return nil, "", translateStoreError(err)
		}
		paginated, nextMarker := paginateRDSItems([]*storerds.DBInstanceSnapshot{snap}, in.Marker, in.MaxRecords, func(s *storerds.DBInstanceSnapshot) string {
			return s.DBSnapshotIdentifier
		})
		return paginated, nextMarker, nil
	}

	snapshots, err := store.ListInstanceSnapshots()
	if err != nil {
		return nil, "", translateStoreError(err)
	}

	filtered := make([]*storerds.DBInstanceSnapshot, 0, len(snapshots))
	for _, snap := range snapshots {
		if in.DBSnapshotIdentifier != "" && snap.DBSnapshotIdentifier != in.DBSnapshotIdentifier {
			continue
		}
		if in.DBInstanceIdentifier != "" && snap.DBInstanceIdentifier != in.DBInstanceIdentifier {
			continue
		}
		if in.SnapshotType != "" && snap.SnapshotType != in.SnapshotType {
			continue
		}
		if !applyRDSFilters(in.Filters, instanceSnapshotFilterGetter(snap)) {
			continue
		}
		filtered = append(filtered, snap)
	}

	paginated, nextMarker := paginateRDSItems(filtered, in.Marker, in.MaxRecords, func(s *storerds.DBInstanceSnapshot) string {
		return s.DBSnapshotIdentifier
	})
	return paginated, nextMarker, nil
}

// QuerySubnetGroups queries, filters, and paginates DBSubnetGroups.
func QuerySubnetGroups(store storerds.StoreInterface, in DescribeDBSubnetGroupsInput) ([]*storerds.DBSubnetGroup, string, error) {
	if in.DBSubnetGroupName != "" {
		group, err := store.GetSubnetGroup(in.DBSubnetGroupName)
		if err != nil {
			return nil, "", translateStoreError(err)
		}
		paginated, nextMarker := paginateRDSItems([]*storerds.DBSubnetGroup{group}, in.Marker, in.MaxRecords, func(g *storerds.DBSubnetGroup) string {
			return g.DBSubnetGroupName
		})
		return paginated, nextMarker, nil
	}

	groups, err := store.ListSubnetGroups()
	if err != nil {
		return nil, "", translateStoreError(err)
	}

	filtered := make([]*storerds.DBSubnetGroup, 0, len(groups))
	for _, g := range groups {
		if in.DBSubnetGroupName != "" && g.DBSubnetGroupName != in.DBSubnetGroupName {
			continue
		}
		if !applyRDSFilters(in.Filters, subnetGroupFilterGetter(g)) {
			continue
		}
		filtered = append(filtered, g)
	}

	paginated, nextMarker := paginateRDSItems(filtered, in.Marker, in.MaxRecords, func(g *storerds.DBSubnetGroup) string {
		return g.DBSubnetGroupName
	})
	return paginated, nextMarker, nil
}

// QueryGlobalClusters queries, filters, and paginates GlobalClusters.
func QueryGlobalClusters(store storerds.StoreInterface, in DescribeGlobalClustersInput) ([]*storerds.GlobalCluster, string, error) {
	if in.GlobalClusterIdentifier != "" {
		gc, err := store.GetGlobalCluster(in.GlobalClusterIdentifier)
		if err != nil {
			return nil, "", translateStoreError(err)
		}
		paginated, nextMarker := paginateRDSItems([]*storerds.GlobalCluster{gc}, in.Marker, in.MaxRecords, func(g *storerds.GlobalCluster) string {
			return g.GlobalClusterIdentifier
		})
		return paginated, nextMarker, nil
	}

	clusters, err := store.ListGlobalClusters()
	if err != nil {
		return nil, "", translateStoreError(err)
	}

	filtered := make([]*storerds.GlobalCluster, 0, len(clusters))
	for _, gc := range clusters {
		if in.GlobalClusterIdentifier != "" && gc.GlobalClusterIdentifier != in.GlobalClusterIdentifier {
			continue
		}
		if !applyRDSFilters(in.Filters, globalClusterFilterGetter(gc)) {
			continue
		}
		filtered = append(filtered, gc)
	}

	paginated, nextMarker := paginateRDSItems(filtered, in.Marker, in.MaxRecords, func(gc *storerds.GlobalCluster) string {
		return gc.GlobalClusterIdentifier
	})
	return paginated, nextMarker, nil
}

// QueryEventSubscriptions queries, filters, and paginates EventSubscriptions.
func QueryEventSubscriptions(store storerds.StoreInterface, in DescribeEventSubscriptionsInput) ([]*storerds.EventSubscription, string, error) {
	if in.SubscriptionName != "" {
		sub, err := store.GetEventSubscription(in.SubscriptionName)
		if err != nil {
			return nil, "", translateStoreError(err)
		}
		paginated, nextMarker := paginateRDSItems([]*storerds.EventSubscription{sub}, in.Marker, in.MaxRecords, func(es *storerds.EventSubscription) string {
			return es.CustSubscriptionId
		})
		return paginated, nextMarker, nil
	}

	subs, err := store.ListEventSubscriptions()
	if err != nil {
		return nil, "", translateStoreError(err)
	}

	filtered := make([]*storerds.EventSubscription, 0, len(subs))
	for _, es := range subs {
		if in.SubscriptionName != "" && es.CustSubscriptionId != in.SubscriptionName {
			continue
		}
		if !applyRDSFilters(in.Filters, eventSubscriptionFilterGetter(es)) {
			continue
		}
		filtered = append(filtered, es)
	}

	paginated, nextMarker := paginateRDSItems(filtered, in.Marker, in.MaxRecords, func(es *storerds.EventSubscription) string {
		return es.CustSubscriptionId
	})
	return paginated, nextMarker, nil
}

// QueryClusterParameterGroups queries, filters, and paginates
// DBClusterParameterGroups.
func QueryClusterParameterGroups(store storerds.StoreInterface, in DescribeDBClusterParameterGroupsInput) ([]*storerds.DBClusterParameterGroup, string, error) {
	if in.DBClusterParameterGroupName != "" {
		group, err := store.GetClusterParameterGroup(in.DBClusterParameterGroupName)
		if err != nil {
			return nil, "", translateStoreError(err)
		}
		paginated, nextMarker := paginateRDSItems([]*storerds.DBClusterParameterGroup{group}, in.Marker, in.MaxRecords, func(g *storerds.DBClusterParameterGroup) string {
			return g.DBClusterParameterGroupName
		})
		return paginated, nextMarker, nil
	}

	groups, err := store.ListClusterParameterGroups()
	if err != nil {
		return nil, "", translateStoreError(err)
	}

	filtered := make([]*storerds.DBClusterParameterGroup, 0, len(groups))
	for _, g := range groups {
		if in.DBClusterParameterGroupName != "" && g.DBClusterParameterGroupName != in.DBClusterParameterGroupName {
			continue
		}
		if !applyRDSFilters(in.Filters, clusterParamGroupFilterGetter(g)) {
			continue
		}
		filtered = append(filtered, g)
	}

	paginated, nextMarker := paginateRDSItems(filtered, in.Marker, in.MaxRecords, func(g *storerds.DBClusterParameterGroup) string {
		return g.DBClusterParameterGroupName
	})
	return paginated, nextMarker, nil
}

// QueryParameterGroups queries, filters, and paginates DBParameterGroups.
func QueryParameterGroups(store storerds.StoreInterface, in DescribeDBParameterGroupsInput) ([]*storerds.DBParameterGroup, string, error) {
	if in.DBParameterGroupName != "" {
		group, err := store.GetParameterGroup(in.DBParameterGroupName)
		if err != nil {
			return nil, "", translateStoreError(err)
		}
		paginated, nextMarker := paginateRDSItems([]*storerds.DBParameterGroup{group}, in.Marker, in.MaxRecords, func(g *storerds.DBParameterGroup) string {
			return g.DBParameterGroupName
		})
		return paginated, nextMarker, nil
	}

	groups, err := store.ListParameterGroups()
	if err != nil {
		return nil, "", translateStoreError(err)
	}

	filtered := make([]*storerds.DBParameterGroup, 0, len(groups))
	for _, g := range groups {
		if in.DBParameterGroupName != "" && g.DBParameterGroupName != in.DBParameterGroupName {
			continue
		}
		if !applyRDSFilters(in.Filters, paramGroupFilterGetter(g)) {
			continue
		}
		filtered = append(filtered, g)
	}

	paginated, nextMarker := paginateRDSItems(filtered, in.Marker, in.MaxRecords, func(g *storerds.DBParameterGroup) string {
		return g.DBParameterGroupName
	})
	return paginated, nextMarker, nil
}

// QueryTags retrieves tags for a resource ARN.
func QueryTags(store storerds.StoreInterface, arn string) ([]types.Tag, error) {
	tags, err := store.GetTags(arn)
	if err != nil {
		return nil, translateStoreError(err)
	}
	return tags, nil
}
