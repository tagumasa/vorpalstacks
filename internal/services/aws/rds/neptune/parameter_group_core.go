package neptune

import (
	"context"
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
)

// CreateDBClusterParameterGroupInput carries the wire-parsed
// CreateDBClusterParameterGroup request.
type CreateDBClusterParameterGroupInput struct {
	DBClusterParameterGroupName string
	DBParameterGroupFamily      string
	Description                 string
	AccountID                   string
	Region                      string
}

// DeleteDBClusterParameterGroupInput carries the wire-parsed
// DeleteDBClusterParameterGroup request.
type DeleteDBClusterParameterGroupInput struct {
	DBClusterParameterGroupName string
}

// DescribeDBClusterParametersInput carries the wire-parsed
// DescribeDBClusterParameters request.
type DescribeDBClusterParametersInput struct {
	DBClusterParameterGroupName string
}

// ModifyDBClusterParameterGroupInput carries the wire-parsed
// ModifyDBClusterParameterGroup request.
type ModifyDBClusterParameterGroupInput struct {
	DBClusterParameterGroupName string
	Parameters                  []neptunestore.Parameter
}

// ResetDBClusterParameterGroupInput carries the wire-parsed
// ResetDBClusterParameterGroup request. The documented default is a full
// reset, so the subset path is selected purely by the presence of the
// Parameters list.
type ResetDBClusterParameterGroupInput struct {
	DBClusterParameterGroupName string
	Parameters                  []neptunestore.Parameter
}

// CopyDBClusterParameterGroupInput carries the wire-parsed
// CopyDBClusterParameterGroup request.
type CopyDBClusterParameterGroupInput struct {
	SourceDBClusterParameterGroupIdentifier  string
	TargetDBClusterParameterGroupIdentifier  string
	TargetDBClusterParameterGroupDescription string
	AccountID                                string
	Region                                   string
}

// CreateDBParameterGroupInput carries the wire-parsed CreateDBParameterGroup
// request.
type CreateDBParameterGroupInput struct {
	DBParameterGroupName   string
	DBParameterGroupFamily string
	Description            string
	AccountID              string
	Region                 string
}

// DeleteDBParameterGroupInput carries the wire-parsed DeleteDBParameterGroup
// request.
type DeleteDBParameterGroupInput struct {
	DBParameterGroupName string
}

// DescribeDBParametersInput carries the wire-parsed DescribeDBParameters
// request.
type DescribeDBParametersInput struct {
	DBParameterGroupName string
}

// ModifyDBParameterGroupInput carries the wire-parsed ModifyDBParameterGroup
// request.
type ModifyDBParameterGroupInput struct {
	DBParameterGroupName string
	Parameters           []neptunestore.Parameter
}

// ResetDBParameterGroupInput carries the wire-parsed ResetDBParameterGroup
// request. The documented default is a full reset, so the subset path is
// selected purely by the presence of the Parameters list.
type ResetDBParameterGroupInput struct {
	DBParameterGroupName string
	Parameters           []neptunestore.Parameter
}

// CopyDBParameterGroupInput carries the wire-parsed CopyDBParameterGroup
// request.
type CopyDBParameterGroupInput struct {
	SourceDBParameterGroupIdentifier  string
	TargetDBParameterGroupIdentifier  string
	TargetDBParameterGroupDescription string
	AccountID                         string
	Region                            string
}

// getClusterParameterGroupCore resolves a cluster parameter group by name.
func (s *NeptuneService) getClusterParameterGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DescribeDBClusterParametersInput) (*neptunestore.DBClusterParameterGroup, error) {
	name := in.DBClusterParameterGroupName
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBClusterParameterGroupName is required")
	}
	return store.GetClusterParameterGroup(name)
}

// getParameterGroupCore resolves an instance parameter group by name.
func (s *NeptuneService) getParameterGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DescribeDBParametersInput) (*neptunestore.DBParameterGroup, error) {
	name := in.DBParameterGroupName
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBParameterGroupName is required")
	}
	return store.GetParameterGroup(name)
}

// createDBClusterParameterGroupCore validates and persists a new cluster
// parameter group.
func (s *NeptuneService) createDBClusterParameterGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *CreateDBClusterParameterGroupInput) (interface{}, error) {
	name := in.DBClusterParameterGroupName
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBClusterParameterGroupName is required")
	}
	family := in.DBParameterGroupFamily
	if family == "" {
		family = "neptune1"
	}

	pg := &neptunestore.DBClusterParameterGroup{
		DBClusterParameterGroupName: name,
		DBParameterGroupFamily:      family,
		Description:                 in.Description,
		ARN:                         neptunestore.ClusterParameterGroupARN(in.AccountID, in.Region, name),
	}
	if err := store.CreateClusterParameterGroup(pg); err != nil {
		return nil, translateStoreError(err)
	}
	return map[string]interface{}{"DBClusterParameterGroup": pg}, nil
}

// deleteDBClusterParameterGroupCore deletes a cluster parameter group and
// clears its tags.
func (s *NeptuneService) deleteDBClusterParameterGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DeleteDBClusterParameterGroupInput) (interface{}, error) {
	name := in.DBClusterParameterGroupName
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBClusterParameterGroupName is required")
	}
	pg, err := store.GetClusterParameterGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}
	if err := store.DeleteClusterParameterGroup(name); err != nil {
		return nil, translateStoreError(err)
	}
	removeTagsForResource(store, pg.ARN)
	return map[string]interface{}{}, nil
}

// validateParameterEntries enforces the ApplyMethod enum shared by the
// Modify and Reset parameter-list members: a supplied value must name a
// documented apply method.
func validateParameterEntries(parameters []neptunestore.Parameter) error {
	for _, p := range parameters {
		switch p.ApplyMethod {
		case "", "immediate", "pending-reboot":
		default:
			return awserrors.NewAWSError("InvalidParameterValue",
				fmt.Sprintf("ApplyMethod must be immediate or pending-reboot, got %s", p.ApplyMethod), http.StatusBadRequest)
		}
	}
	return nil
}

// modifyDBClusterParameterGroupCore merges parameter modifications into a
// cluster parameter group.
func (s *NeptuneService) modifyDBClusterParameterGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *ModifyDBClusterParameterGroupInput) (interface{}, error) {
	name := in.DBClusterParameterGroupName
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBClusterParameterGroupName is required")
	}
	if err := validateParameterEntries(in.Parameters); err != nil {
		return nil, err
	}
	pg, err := store.GetClusterParameterGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}
	if len(in.Parameters) > 0 {
		pg.Parameters = applyParameterModifications(pg.Parameters, in.Parameters)
		if err := store.UpdateClusterParameterGroup(pg); err != nil {
			return nil, translateStoreError(err)
		}
	}
	return map[string]interface{}{"DBClusterParameterGroupName": name}, nil
}

// resetDBClusterParameterGroupCore resets a cluster parameter group's
// parameters to their defaults. A request carrying the documented
// Parameters list resets only the named entries; every other request
// resets the whole group.
func (s *NeptuneService) resetDBClusterParameterGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *ResetDBClusterParameterGroupInput) (interface{}, error) {
	name := in.DBClusterParameterGroupName
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBClusterParameterGroupName is required")
	}
	if err := validateParameterEntries(in.Parameters); err != nil {
		return nil, err
	}
	pg, err := store.GetClusterParameterGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}
	if len(in.Parameters) > 0 {
		pg.Parameters = resetNamedParameters(pg.Parameters, in.Parameters)
	} else {
		pg.Parameters = nil
	}
	if err := store.UpdateClusterParameterGroup(pg); err != nil {
		return nil, translateStoreError(err)
	}
	return map[string]interface{}{"DBClusterParameterGroupName": name}, nil
}

// copyDBClusterParameterGroupCore copies a cluster parameter group under a
// new identifier.
func (s *NeptuneService) copyDBClusterParameterGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *CopyDBClusterParameterGroupInput) (interface{}, error) {
	sourceName := in.SourceDBClusterParameterGroupIdentifier
	if sourceName == "" {
		return nil, awserrors.NewMissingParameter("SourceDBClusterParameterGroupIdentifier is required")
	}
	targetName := in.TargetDBClusterParameterGroupIdentifier
	if targetName == "" {
		return nil, awserrors.NewMissingParameter("TargetDBClusterParameterGroupIdentifier is required")
	}
	source, err := store.GetClusterParameterGroup(sourceName)
	if err != nil {
		return nil, translateStoreError(err)
	}
	desc := in.TargetDBClusterParameterGroupDescription
	if desc == "" {
		desc = source.Description
	}
	pg := &neptunestore.DBClusterParameterGroup{
		DBClusterParameterGroupName: targetName,
		DBParameterGroupFamily:      source.DBParameterGroupFamily,
		Description:                 desc,
		ARN:                         neptunestore.ClusterParameterGroupARN(in.AccountID, in.Region, targetName),
		Parameters:                  append([]neptunestore.Parameter(nil), source.Parameters...),
	}
	if err := store.CreateClusterParameterGroup(pg); err != nil {
		return nil, translateStoreError(err)
	}
	return map[string]interface{}{"DBClusterParameterGroup": pg}, nil
}

// createDBParameterGroupCore validates and persists a new instance parameter
// group.
func (s *NeptuneService) createDBParameterGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *CreateDBParameterGroupInput) (interface{}, error) {
	name := in.DBParameterGroupName
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBParameterGroupName is required")
	}
	family := in.DBParameterGroupFamily
	if family == "" {
		family = "neptune1"
	}
	pg := &neptunestore.DBParameterGroup{
		DBParameterGroupName:   name,
		DBParameterGroupFamily: family,
		Description:            in.Description,
		ARN:                    neptunestore.ParameterGroupARN(in.AccountID, in.Region, name),
	}
	if err := store.CreateParameterGroup(pg); err != nil {
		return nil, translateStoreError(err)
	}
	return map[string]interface{}{"DBParameterGroup": pg}, nil
}

// deleteDBParameterGroupCore deletes an instance parameter group and clears
// its tags.
func (s *NeptuneService) deleteDBParameterGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DeleteDBParameterGroupInput) (interface{}, error) {
	name := in.DBParameterGroupName
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBParameterGroupName is required")
	}
	pg, err := store.GetParameterGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}
	if err := store.DeleteParameterGroup(name); err != nil {
		return nil, translateStoreError(err)
	}
	removeTagsForResource(store, pg.ARN)
	return map[string]interface{}{}, nil
}

// modifyDBParameterGroupCore merges parameter modifications into an instance
// parameter group.
func (s *NeptuneService) modifyDBParameterGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *ModifyDBParameterGroupInput) (interface{}, error) {
	name := in.DBParameterGroupName
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBParameterGroupName is required")
	}
	if err := validateParameterEntries(in.Parameters); err != nil {
		return nil, err
	}
	pg, err := store.GetParameterGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}
	if len(in.Parameters) > 0 {
		pg.Parameters = applyParameterModifications(pg.Parameters, in.Parameters)
		if err := store.UpdateParameterGroup(pg); err != nil {
			return nil, translateStoreError(err)
		}
	}
	return map[string]interface{}{"DBParameterGroupName": name}, nil
}

// resetDBParameterGroupCore resets an instance parameter group's parameters
// to their defaults. A request carrying the documented Parameters list
// resets only the named entries; every other request resets the whole
// group.
func (s *NeptuneService) resetDBParameterGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *ResetDBParameterGroupInput) (interface{}, error) {
	name := in.DBParameterGroupName
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBParameterGroupName is required")
	}
	if err := validateParameterEntries(in.Parameters); err != nil {
		return nil, err
	}
	pg, err := store.GetParameterGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}
	if len(in.Parameters) > 0 {
		pg.Parameters = resetNamedParameters(pg.Parameters, in.Parameters)
	} else {
		pg.Parameters = nil
	}
	if err := store.UpdateParameterGroup(pg); err != nil {
		return nil, translateStoreError(err)
	}
	return map[string]interface{}{"DBParameterGroupName": name}, nil
}

// copyDBParameterGroupCore copies an instance parameter group under a new
// identifier.
func (s *NeptuneService) copyDBParameterGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *CopyDBParameterGroupInput) (interface{}, error) {
	sourceName := in.SourceDBParameterGroupIdentifier
	if sourceName == "" {
		return nil, awserrors.NewMissingParameter("SourceDBParameterGroupIdentifier is required")
	}
	targetName := in.TargetDBParameterGroupIdentifier
	if targetName == "" {
		return nil, awserrors.NewMissingParameter("TargetDBParameterGroupIdentifier is required")
	}
	source, err := store.GetParameterGroup(sourceName)
	if err != nil {
		return nil, translateStoreError(err)
	}
	desc := in.TargetDBParameterGroupDescription
	if desc == "" {
		desc = source.Description
	}
	pg := &neptunestore.DBParameterGroup{
		DBParameterGroupName:   targetName,
		DBParameterGroupFamily: source.DBParameterGroupFamily,
		Description:            desc,
		ARN:                    neptunestore.ParameterGroupARN(in.AccountID, in.Region, targetName),
		Parameters:             append([]neptunestore.Parameter(nil), source.Parameters...),
	}
	if err := store.CreateParameterGroup(pg); err != nil {
		return nil, translateStoreError(err)
	}
	return map[string]interface{}{"DBParameterGroup": pg}, nil
}
