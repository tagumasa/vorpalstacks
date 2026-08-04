package rds

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	pb "vorpalstacks/internal/pb/aws/rds"
	storerds "vorpalstacks/internal/store/aws/rds"
)

// translateStoreError converts a common RDS store-level sentinel error
// into the corresponding AWS-compatible *AWSError so that the gRPC-Web
// admin handler can map it via svcerrors.AWSErrorToGRPC. Unknown errors
// are returned unchanged and will fall through to CodeInternal.
//
// The sentinel variables live in the shared store package
// (internal/store/aws/rds/errors.go) and are re-exported by each
// engine-specific sub-store (neptune, etc.), so a single translation
// table covers every engine.
func translateStoreError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	// --- NotFound family (HTTP 404) ---
	case errors.Is(err, storerds.ErrDBClusterNotFound):
		return awserrors.NewAWSError("DBClusterNotFoundFault", err.Error(), http.StatusNotFound)
	case errors.Is(err, storerds.ErrDBInstanceNotFound):
		return awserrors.NewAWSError("DBInstanceNotFoundFault", err.Error(), http.StatusNotFound)
	case errors.Is(err, storerds.ErrDBClusterSnapshotNotFound):
		return awserrors.NewAWSError("DBClusterSnapshotNotFoundFault", err.Error(), http.StatusNotFound)
	case errors.Is(err, storerds.ErrDBSnapshotNotFound):
		return awserrors.NewAWSError("DBSnapshotNotFoundFault", err.Error(), http.StatusNotFound)
	case errors.Is(err, storerds.ErrDBClusterParameterGroupNotFound):
		return awserrors.NewAWSError("DBClusterParameterGroupNotFoundFault", err.Error(), http.StatusNotFound)
	case errors.Is(err, storerds.ErrDBParameterGroupNotFound):
		return awserrors.NewAWSError("DBParameterGroupNotFoundFault", err.Error(), http.StatusNotFound)
	case errors.Is(err, storerds.ErrDBSubnetGroupNotFound):
		return awserrors.NewAWSError("DBSubnetGroupNotFoundFault", err.Error(), http.StatusNotFound)
	case errors.Is(err, storerds.ErrGlobalClusterNotFound):
		return awserrors.NewAWSError("GlobalClusterNotFoundFault", err.Error(), http.StatusNotFound)
	case errors.Is(err, storerds.ErrEventSubscriptionNotFound):
		return awserrors.NewAWSError("SubscriptionNotFoundFault", err.Error(), http.StatusNotFound)
	case errors.Is(err, storerds.ErrEventNotFound):
		return awserrors.NewAWSError("EventNotFoundFault", err.Error(), http.StatusNotFound)

	// --- AlreadyExists family (HTTP 409) ---
	case errors.Is(err, storerds.ErrDBClusterAlreadyExists):
		return awserrors.NewAWSError("DBClusterAlreadyExistsFault", err.Error(), http.StatusConflict)
	case errors.Is(err, storerds.ErrDBInstanceAlreadyExists):
		return awserrors.NewAWSError("DBInstanceAlreadyExistsFault", err.Error(), http.StatusConflict)
	case errors.Is(err, storerds.ErrDBClusterSnapshotAlreadyExists):
		return awserrors.NewAWSError("DBClusterSnapshotAlreadyExistsFault", err.Error(), http.StatusConflict)
	case errors.Is(err, storerds.ErrDBSnapshotAlreadyExists):
		return awserrors.NewAWSError("DBSnapshotAlreadyExistsFault", err.Error(), http.StatusConflict)
	case errors.Is(err, storerds.ErrDBClusterParameterGroupAlreadyExists):
		return awserrors.NewAWSError("DBClusterParameterGroupAlreadyExistsFault", err.Error(), http.StatusConflict)
	case errors.Is(err, storerds.ErrDBParameterGroupAlreadyExists):
		return awserrors.NewAWSError("DBParameterGroupAlreadyExistsFault", err.Error(), http.StatusConflict)
	case errors.Is(err, storerds.ErrDBSubnetGroupAlreadyExists):
		return awserrors.NewAWSError("DBSubnetGroupAlreadyExistsFault", err.Error(), http.StatusConflict)
	case errors.Is(err, storerds.ErrGlobalClusterAlreadyExists):
		return awserrors.NewAWSError("GlobalClusterAlreadyExistsFault", err.Error(), http.StatusConflict)
	case errors.Is(err, storerds.ErrEventSubscriptionAlreadyExists):
		return awserrors.NewAWSError("SubscriptionAlreadyExistFault", err.Error(), http.StatusConflict)
	case errors.Is(err, storerds.ErrEventAlreadyExists):
		return awserrors.NewAWSError("EventAlreadyExistsFault", err.Error(), http.StatusConflict)

	// --- InvalidParameter family (HTTP 400) ---
	case errors.Is(err, storerds.ErrInvalidParameterGroupState):
		return awserrors.NewAWSError("InvalidDBParameterGroupStateFault", err.Error(), http.StatusBadRequest)
	case errors.Is(err, storerds.ErrInvalidEventMarker):
		return awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	return err
}

// ---------------------------------------------------------------------------
// Error helpers used by Core methods. These return *AWSError so that
// svcerrors.AWSErrorToGRPC maps them to the correct connect.Code.
// ---------------------------------------------------------------------------

func newValidationError(format string, args ...interface{}) error {
	return awserrors.NewValidationException(fmt.Sprintf(format, args...))
}

func newFailedPreconditionError(format string, args ...interface{}) error {
	return awserrors.NewDeleteConflictException(fmt.Sprintf(format, args...))
}

func newInternalError(format string, args ...interface{}) error {
	return awserrors.NewInternalErrorException(fmt.Sprintf(format, args...))
}

// sortParameters sorts a slice of protobuf Parameters by name. Used by
// DescribeDBClusterParameters and DescribeDBParameters Core methods.
func sortParameters(params []*pb.Parameter) {
	sort.Slice(params, func(i, j int) bool {
		return params[i].Parametername < params[j].Parametername
	})
}

// applyRDSFilters reports whether a candidate resource matches every filter
// in the supplied list. A resource matches a single filter when the value
// returned by getter(name) is equal (case-insensitive) to at least one of
// the filter's Values; an empty filter list matches everything.
//
// Semantics mirror AWS RDS: OR within a single filter's Values, AND across
// multiple filters. Unknown filter names cause the resource to be excluded
// rather than silently matching.
func applyRDSFilters(filters []*pb.Filter, getter func(name string) (string, bool)) bool {
	if len(filters) == 0 {
		return true
	}
	for _, f := range filters {
		if f == nil {
			continue
		}
		v, ok := getter(f.Name)
		if !ok {
			return false
		}
		matched := false
		for _, want := range f.Values {
			if strings.EqualFold(v, want) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}
