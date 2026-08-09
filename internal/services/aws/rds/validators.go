package rds

import (
	"fmt"
	"strings"

	pb "vorpalstacks/internal/pb/aws/rds"
)

// AWS RDS identifier rules (CreateDBInstance / CreateDBCluster API reference):
//
//   - 1 to 63 characters
//   - lowercase letters, digits, and hyphens
//   - first character must be a letter
//   - cannot end with a hyphen
//   - cannot contain consecutive hyphens
//
// AWS RDS database name rules:
//
//   - 1 to 64 alphanumeric characters
//   - first character must be a letter
//
// Apply at every Create entry point so that the store never receives a
// malformed name that would later confuse filter logic or pollute
// downstream API responses.

// ValidateDBInstanceIdentifier rejects identifiers that violate the AWS
// RDS DBInstanceIdentifier constraints.
func ValidateDBInstanceIdentifier(id string) error {
	if len(id) < 1 || len(id) > 63 {
		return fmt.Errorf("DBInstanceIdentifier must be 1-63 characters, got %d", len(id))
	}
	if !isASCIILetter(id[0]) {
		return fmt.Errorf("DBInstanceIdentifier %q must start with a letter", id)
	}
	for i := 0; i < len(id); i++ {
		c := id[i]
		if !(isASCIILetter(c) || isASCIIDigit(c) || c == '-') {
			return fmt.Errorf("DBInstanceIdentifier %q contains invalid character %q (allowed: a-z, 0-9, -)", id, string(rune(c)))
		}
	}
	if id[len(id)-1] == '-' {
		return fmt.Errorf("DBInstanceIdentifier %q cannot end with a hyphen", id)
	}
	if strings.Contains(id, "--") {
		return fmt.Errorf("DBInstanceIdentifier %q cannot contain consecutive hyphens", id)
	}
	return nil
}

// ValidateDBClusterIdentifier mirrors ValidateDBInstanceIdentifier for
// the DBClusterIdentifier parameter. AWS uses identical rules.
func ValidateDBClusterIdentifier(id string) error {
	if len(id) < 1 || len(id) > 63 {
		return fmt.Errorf("DBClusterIdentifier must be 1-63 characters, got %d", len(id))
	}
	if !isASCIILetter(id[0]) {
		return fmt.Errorf("DBClusterIdentifier %q must start with a letter", id)
	}
	for i := 0; i < len(id); i++ {
		c := id[i]
		if !(isASCIILetter(c) || isASCIIDigit(c) || c == '-') {
			return fmt.Errorf("DBClusterIdentifier %q contains invalid character %q", id, string(rune(c)))
		}
	}
	if id[len(id)-1] == '-' {
		return fmt.Errorf("DBClusterIdentifier %q cannot end with a hyphen", id)
	}
	if strings.Contains(id, "--") {
		return fmt.Errorf("DBClusterIdentifier %q cannot contain consecutive hyphens", id)
	}
	return nil
}

// ValidateDatabaseName enforces AWS RDS DatabaseName constraints.
func ValidateDatabaseName(name string) error {
	if len(name) < 1 || len(name) > 64 {
		return fmt.Errorf("DatabaseName must be 1-64 characters, got %d", len(name))
	}
	if !isASCIILetter(name[0]) {
		return fmt.Errorf("DatabaseName %q must start with a letter", name)
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		if !(isASCIILetter(c) || isASCIIDigit(c) || c == '_') {
			return fmt.Errorf("DatabaseName %q contains invalid character %q (allowed: a-z, 0-9, _)", name, string(rune(c)))
		}
	}
	return nil
}

// ValidateDBInstanceClass enforces that DBInstanceClass is supplied.
// AWS RDS rejects empty DBInstanceClass with InvalidParameterValue.
func ValidateDBInstanceClass(class string) error {
	if class == "" {
		return fmt.Errorf("DBInstanceClass is required")
	}
	return nil
}

// MysqlEngineVersion describes a single MySQL engine version and its
// associated metadata. This is the single source of truth consumed by
// both ValidateEngineVersion and DescribeDBEngineVersions so that the
// list a client sees via DescribeDBEngineVersions is exactly the list
// the validator accepts.
type MysqlEngineVersion struct {
	Version   string
	Family    string
	Major     string
	DescShort string
}

// NeptuneEngineVersion describes a single Neptune engine version.
type NeptuneEngineVersion struct {
	Version       string
	Family        string
	ParallelQuery bool
}

// supportedMysqlVersions enumerates the MySQL engine versions this
// platform recognises. The list mirrors the major MySQL LTS releases
// (5.7, 8.0, 8.4) plus their most current minor revisions. AWS RDS for
// MySQL supports many additional revisions; the list here intentionally
// covers the parameter-group families we serve so that
// DescribeDBParameters / DescribeDBClusterParameters can pick the right
// defaults for whichever EngineVersion the caller selected.
var supportedMysqlVersions = []MysqlEngineVersion{
	{"5.7.44", "mysql5.7", "5.7", "MySQL 5.7.44"},
	{"5.7.45", "mysql5.7", "5.7", "MySQL 5.7.45"},
	{"5.7.46", "mysql5.7", "5.7", "MySQL 5.7.46"},
	{"5.7.47", "mysql5.7", "5.7", "MySQL 5.7.47"},
	{"5.7.48", "mysql5.7", "5.7", "MySQL 5.7.48"},
	{"8.0.35", "mysql8.0", "8.0", "MySQL 8.0.35"},
	{"8.0.36", "mysql8.0", "8.0", "MySQL 8.0.36"},
	{"8.0.37", "mysql8.0", "8.0", "MySQL 8.0.37"},
	{"8.0.38", "mysql8.0", "8.0", "MySQL 8.0.38"},
	{"8.0.39", "mysql8.0", "8.0", "MySQL 8.0.39"},
	{"8.0.40", "mysql8.0", "8.0", "MySQL 8.0.40"},
	{"8.0.41", "mysql8.0", "8.0", "MySQL 8.0.41"},
	{"8.0.42", "mysql8.0", "8.0", "MySQL 8.0.42"},
	{"8.0.43", "mysql8.0", "8.0", "MySQL 8.0.43"},
	{"8.0.44", "mysql8.0", "8.0", "MySQL 8.0.44"},
	{"8.0.45", "mysql8.0", "8.0", "MySQL 8.0.45"},
	{"8.0.46", "mysql8.0", "8.0", "MySQL 8.0.46"},
	{"8.4.3", "mysql8.4", "8.4", "MySQL 8.4.3"},
	{"8.4.4", "mysql8.4", "8.4", "MySQL 8.4.4"},
	{"8.4.5", "mysql8.4", "8.4", "MySQL 8.4.5"},
	{"8.4.6", "mysql8.4", "8.4", "MySQL 8.4.6"},
	{"8.4.7", "mysql8.4", "8.4", "MySQL 8.4.7"},
	{"8.4.8", "mysql8.4", "8.4", "MySQL 8.4.8"},
	{"8.4.9", "mysql8.4", "8.4", "MySQL 8.4.9"},
	{"8.4.10", "mysql8.4", "8.4", "MySQL 8.4.10"},
}

// supportedNeptuneVersions enumerates the Neptune engine versions this
// platform recognises. ParallelQuery is true for 1.3.x and above,
// matching the Neptune engine's Streaming Changes / ParallelQuery
// feature availability.
var supportedNeptuneVersions = []NeptuneEngineVersion{
	{"1.2.1.0", "neptune1", false},
	{"1.2.2.0", "neptune1", false},
	{"1.2.3.0", "neptune1", false},
	{"1.3.0.0", "neptune1", true},
	{"1.3.1.0", "neptune1", true},
	{"1.3.2.0", "neptune1", true},
	{"1.4.0.0", "neptune1", true},
	{"1.4.0.1", "neptune1", true},
	{"1.4.1.0", "neptune1", true},
	{"1.4.2.0", "neptune1", true},
	{"1.4.3.0", "neptune1", true},
	{"1.4.4.0", "neptune1", true},
	{"1.4.5.0", "neptune1", true},
	{"1.4.5.1", "neptune1", true},
	{"1.4.6.0", "neptune1", true},
	{"1.4.6.2", "neptune1", true},
	{"1.4.7.0", "neptune1", true},
}

// SupportedMysqlVersions returns the list of supported MySQL engine
// versions. This is the single source of truth used by both the
// validator and the DescribeDBEngineVersions handler.
func SupportedMysqlVersions() []MysqlEngineVersion {
	return supportedMysqlVersions
}

// SupportedNeptuneVersions returns the list of supported Neptune engine
// versions.
func SupportedNeptuneVersions() []NeptuneEngineVersion {
	return supportedNeptuneVersions
}

// IsSupportedEngine returns true if the engine type is recognised by
// this platform.
func IsSupportedEngine(engine string) bool {
	switch engine {
	case "neptune", "mysql":
		return true
	}
	return false
}

// DefaultEngineVersion returns the default version for the given engine,
// matching AWS behaviour where CreateDBInstance / CreateDBCluster
// without an explicit EngineVersion selects the latest supported
// version. Returns empty string for unsupported engines.
func DefaultEngineVersion(engine string) string {
	switch engine {
	case "mysql":
		if len(supportedMysqlVersions) > 0 {
			return supportedMysqlVersions[len(supportedMysqlVersions)-1].Version
		}
	case "neptune":
		if len(supportedNeptuneVersions) > 0 {
			return supportedNeptuneVersions[len(supportedNeptuneVersions)-1].Version
		}
	}
	return ""
}

// ValidateEngine rejects unsupported engine types. Call this at every
// Create entry point to fail fast before the engine is dispatched to
// engineFor.
func ValidateEngine(engine string) error {
	if !IsSupportedEngine(engine) {
		return fmt.Errorf("Engine %q is not supported on this platform (supported: mysql, neptune)", engine)
	}
	return nil
}

// ValidateDBSnapshotIdentifier enforces AWS RDS DBSnapshotIdentifier
// constraints: 1-255 characters, alphanumeric and hyphens, must start
// with a letter, cannot end with a hyphen, cannot contain consecutive
// hyphens.
func ValidateDBSnapshotIdentifier(id string) error {
	if len(id) < 1 || len(id) > 255 {
		return fmt.Errorf("DBSnapshotIdentifier must be 1-255 characters, got %d", len(id))
	}
	if !isASCIILetter(id[0]) {
		return fmt.Errorf("DBSnapshotIdentifier %q must start with a letter", id)
	}
	for i := 0; i < len(id); i++ {
		c := id[i]
		if !(isASCIILetter(c) || isASCIIDigit(c) || c == '-') {
			return fmt.Errorf("DBSnapshotIdentifier %q contains invalid character %q (allowed: a-z, 0-9, -)", id, string(rune(c)))
		}
	}
	if id[len(id)-1] == '-' {
		return fmt.Errorf("DBSnapshotIdentifier %q cannot end with a hyphen", id)
	}
	if strings.Contains(id, "--") {
		return fmt.Errorf("DBSnapshotIdentifier %q cannot contain consecutive hyphens", id)
	}
	return nil
}

// allEngineVersions returns the complete list of engine versions the
// platform recognises. Used by DescribeDBEngineVersions.
func allEngineVersions() []*pb.DBEngineVersion {
	out := make([]*pb.DBEngineVersion, 0, len(supportedMysqlVersions)+len(supportedNeptuneVersions))
	for _, v := range supportedNeptuneVersions {
		out = append(out, &pb.DBEngineVersion{
			Engine:                 "neptune",
			Engineversion:          v.Version,
			Dbparametergroupfamily: v.Family,
		})
	}
	for _, v := range supportedMysqlVersions {
		out = append(out, &pb.DBEngineVersion{
			Engine:                     "mysql",
			Engineversion:              v.Version,
			Dbparametergroupfamily:     v.Family,
			Dbenginedescription:        v.DescShort,
			Dbengineversiondescription: v.DescShort,
		})
	}
	return out
}

// ValidateEngineVersion rejects an unsupported engine/version
// combination. An empty version string is allowed so that callers can
// apply DefaultEngineVersion after validation passes.
func ValidateEngineVersion(engine, version string) error {
	switch engine {
	case "mysql":
		if version == "" {
			return nil
		}
		for _, v := range supportedMysqlVersions {
			if v.Version == version {
				return nil
			}
		}
		return fmt.Errorf("MySQL engine version %q is not supported on this platform", version)
	case "neptune":
		if version == "" {
			return nil
		}
		for _, v := range supportedNeptuneVersions {
			if v.Version == version {
				return nil
			}
		}
		return fmt.Errorf("Neptune engine version %q is not supported on this platform", version)
	default:
		return fmt.Errorf("Engine %q is not supported on this platform (supported: mysql, neptune)", engine)
	}
}

// ValidatePort enforces the AWS RDS Port constraint: 1150-65535.
// A value of 0 means "unset" (the caller will choose a default) and
// is therefore accepted.
func ValidatePort(port int32) error {
	if port == 0 {
		return nil
	}
	if port < 1150 || port > 65535 {
		return fmt.Errorf("Port must be between 1150 and 65535, got %d", port)
	}
	return nil
}

// ValidateMonitoringInterval enforces the AWS RDS MonitoringInterval
// constraint: one of {0, 1, 5, 10, 15, 30, 60} seconds. A value of 0
// disables Enhanced Monitoring.
func ValidateMonitoringInterval(interval int32) error {
	switch interval {
	case 0, 1, 5, 10, 15, 30, 60:
		return nil
	}
	return fmt.Errorf("MonitoringInterval must be one of 0, 1, 5, 10, 15, 30, 60, got %d", interval)
}

// ValidateStorageType enforces the AWS RDS StorageType constraint.
// MySQL accepts gp2, gp3, io2, and standard. Neptune (Aurora) storage
// is managed by the service and does not accept a StorageType parameter.
// An empty string means "unset" and is accepted so callers that do not
// supply the parameter are not penalised.
func ValidateStorageType(storageType, engine string) error {
	if storageType == "" {
		return nil
	}
	switch engine {
	case "mysql":
		switch storageType {
		case "gp2", "gp3", "io2", "standard":
			return nil
		}
		return fmt.Errorf("StorageType %q is not valid for MySQL engine (allowed: gp2, gp3, io2, standard)", storageType)
	case "neptune":
		return fmt.Errorf("StorageType is not applicable for Neptune engine")
	default:
		switch storageType {
		case "gp2", "gp3", "io2", "standard":
			return nil
		}
		return fmt.Errorf("StorageType %q is not recognised (allowed: gp2, gp3, io2, standard)", storageType)
	}
}

// ValidateBackupRetentionPeriod enforces the AWS RDS
// BackupRetentionPeriod constraint: 0-35 days. A value of 0 disables
// automated backups.
func ValidateBackupRetentionPeriod(period int32) error {
	if period < 0 || period > 35 {
		return fmt.Errorf("BackupRetentionPeriod must be between 0 and 35, got %d", period)
	}
	return nil
}

// ValidateAllocatedStorage enforces the AWS RDS AllocatedStorage
// constraint. For MySQL the valid range is 20-65536 GiB. Neptune
// (Aurora) storage is automatically managed and the field is ignored,
// so an empty/zero value is accepted for Neptune.
func ValidateAllocatedStorage(storage int32, engine string) error {
	if engine == "neptune" {
		return nil
	}
	if storage == 0 {
		return nil
	}
	if storage < 20 || storage > 65536 {
		return fmt.Errorf("AllocatedStorage must be between 20 and 65536 GiB, got %d", storage)
	}
	return nil
}

func isASCIILetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isASCIIDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
