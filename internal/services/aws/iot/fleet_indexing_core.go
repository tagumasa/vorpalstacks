package iot

import (
	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Core functions for the fleet-indexing search path. The query grammar and
// matching engine live in fleet_indexing.go; this file is the
// store-and-broker boundary the search handlers funnel through, so both
// protocol planes share one persistence path.

// searchThingsCore resolves the query and returns every thing matching it,
// paginating through the whole thing list internally.
func (s *IoTService) searchThingsCore(store iotstore.IotStoreInterface, queryString string) ([]*iotstore.Thing, error) {
	qNode, err := parseQuery(queryString)
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	// Pre-compute the set of connected thing names from all brokers in
	// the current region.  This is used by the isConnected query filter.
	conn := s.buildConnectedSetCore(store)

	// Paginate through ALL things, applying the query filter.
	var matched []*iotstore.Thing
	var opts storecommon.ListOptions
	for {
		result, err := store.ListThings(opts, "", "")
		if err != nil {
			return nil, err
		}
		for i := range result.Items {
			if qNode.match(result.Items[i], conn) {
				matched = append(matched, result.Items[i])
			}
		}
		if result.NextMarker == "" {
			break
		}
		opts.Marker = result.NextMarker
	}
	return matched, nil
}

// buildConnectedSetCore returns a set of thing names that have at least
// one principal certificate currently connected to any broker. The
// account and region are required to construct the full certificate ARN
// that ListThingsForPrincipal uses as a PebbleDB key prefix.
func (s *IoTService) buildConnectedSetCore(store iotstore.IotStoreInterface) connectedSet {
	conn := connectedSet{}
	accountID := store.GetAccountID()
	region := store.GetRegion()
	for _, b := range s.brokers {
		for _, certID := range b.ConnectedCertIDs() {
			// ListThingsForPrincipal keys on the full certificate ARN,
			// not the raw SHA-256 hash.
			principal := iotstore.BuildCertificateARN(accountID, region, certID)
			thingNames, err := store.ListThingsForPrincipal(principal)
			if err != nil {
				continue
			}
			for _, tn := range thingNames {
				conn[tn] = true
			}
		}
	}
	return conn
}
