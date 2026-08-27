package iot

import (
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Core functions for the thing family's describe-support and V2 principal
// listing operations. Handlers on both protocol planes are thin adapters;
// validation and persistence live here only.

// ThingConnectivityResult reports whether any certificate principal of the
// thing is currently connected to a broker.
type ThingConnectivityResult struct {
	ThingName        string
	Connected        bool
	ConnectTime      int64
	TimestampSeconds int64
}

// getThingConnectivityCore resolves the thing and reports connectivity by
// scanning every broker's connected-certificate set (a thing is reachable
// through any region's broker, not just the request region's).
func (s *IoTService) getThingConnectivityCore(store iotstore.IotStoreInterface, thingName string) (*ThingConnectivityResult, error) {
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if _, err := store.GetThing(thingName); err != nil {
		return nil, err
	}

	// Check if any certificate principal attached to this thing is currently
	// connected to any MQTT broker (not just the request-region broker).
	connected := false
	connectedAt := int64(0)
	principals, _ := store.ListPrincipalsForThing(thingName)
	for _, principal := range principals {
		certID := extractCertIDFromPrincipal(principal)
		if certID == "" {
			continue
		}
		for _, brk := range s.brokers {
			if c, ts := brk.IsCertConnected(certID); c {
				connected = true
				connectedAt = ts
				break
			}
		}
		if connected {
			break
		}
	}

	return &ThingConnectivityResult{
		ThingName:        thingName,
		Connected:        connected,
		ConnectTime:      connectedAt,
		TimestampSeconds: time.Now().UTC().Unix(),
	}, nil
}

// listThingPrincipalsV2Core returns the principals attached to a thing.
func (s *IoTService) listThingPrincipalsV2Core(store iotstore.IotStoreInterface, thingName string) ([]string, error) {
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	return store.ListPrincipalsForThing(thingName)
}

// listPrincipalThingsV2Core returns the things attached to a principal.
func (s *IoTService) listPrincipalThingsV2Core(store iotstore.IotStoreInterface, principal string) ([]string, error) {
	if principal == "" {
		return nil, iotstore.ErrMissingParam
	}
	return store.ListThingsForPrincipal(principal)
}
