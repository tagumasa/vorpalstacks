package rds

import (
	"strings"
	"testing"
)

func TestValidateDBInstanceIdentifier(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"valid simple", "mydb", false},
		{"valid with hyphens", "my-db-instance", false},
		{"valid alphanumeric", "db123", false},
		{"valid max length", strings.Repeat("a", 63), false},
		{"empty", "", true},
		{"too long", strings.Repeat("a", 64), true},
		{"starts with digit", "1db", true},
		{"starts with hyphen", "-db", true},
		{"ends with hyphen", "db-", true},
		{"consecutive hyphens", "db--1", true},
		{"contains slash", "db/1", true},
		{"contains colon", "db:1", true},
		{"contains underscore", "db_1", true},
		{"contains space", "db 1", true},
		{"single letter", "a", false},
		{"uppercase", "MyDB", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDBInstanceIdentifier(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDBInstanceIdentifier(%q) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDBClusterIdentifier(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"valid simple", "mycluster", false},
		{"valid with hyphens", "my-cluster-1", false},
		{"empty", "", true},
		{"too long", strings.Repeat("a", 64), true},
		{"starts with digit", "1cluster", true},
		{"ends with hyphen", "cluster-", true},
		{"consecutive hyphens", "c--1", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDBClusterIdentifier(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDBClusterIdentifier(%q) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDatabaseName(t *testing.T) {
	tests := []struct {
		name    string
		dbName  string
		wantErr bool
	}{
		{"valid simple", "mydb", false},
		{"valid with digits", "db123", false},
		{"valid with underscore", "my_db", false},
		{"valid max length", strings.Repeat("a", 64), false},
		{"empty", "", true},
		{"too long", strings.Repeat("a", 65), true},
		{"starts with digit", "1db", true},
		{"starts with underscore", "_db", true},
		{"contains hyphen", "my-db", true},
		{"contains space", "my db", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDatabaseName(tt.dbName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDatabaseName(%q) error = %v, wantErr %v", tt.dbName, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDBInstanceClass(t *testing.T) {
	if err := ValidateDBInstanceClass(""); err == nil {
		t.Error("ValidateDBInstanceClass('') should error")
	}
	if err := ValidateDBInstanceClass("db.r5.large"); err != nil {
		t.Errorf("ValidateDBInstanceClass('db.r5.large') unexpected error: %v", err)
	}
}

func TestValidateEngineVersion(t *testing.T) {
	// MySQL versions — pick first and last from supportedMysqlVersions
	mysqlFirst := supportedMysqlVersions[0].Version
	mysqlLast := supportedMysqlVersions[len(supportedMysqlVersions)-1].Version

	// Neptune versions — pick first and last
	neptuneFirst := supportedNeptuneVersions[0].Version
	neptuneLast := supportedNeptuneVersions[len(supportedNeptuneVersions)-1].Version

	tests := []struct {
		name    string
		engine  string
		version string
		wantErr bool
	}{
		{"mysql first version", "mysql", mysqlFirst, false},
		{"mysql last version", "mysql", mysqlLast, false},
		{"mysql empty version allowed", "mysql", "", false},
		{"mysql unsupported version", "mysql", "3.0.0", true},
		{"mysql fictional version", "mysql", "99.99.99", true},
		{"neptune first version", "neptune", neptuneFirst, false},
		{"neptune last version", "neptune", neptuneLast, false},
		{"neptune empty version allowed", "neptune", "", false},
		{"neptune unsupported version", "neptune", "0.0.0.0", true},
		{"unsupported engine", "postgres", "15", true},
		{"unsupported aurora engine", "aurora-mysql", "8.0", true},
		{"empty engine", "", "1.0", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEngineVersion(tt.engine, tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEngineVersion(%q, %q) error = %v, wantErr %v",
					tt.engine, tt.version, err, tt.wantErr)
			}
		})
	}
}

func TestValidateEngine(t *testing.T) {
	tests := []struct {
		engine  string
		wantErr bool
	}{
		{"mysql", false},
		{"neptune", false},
		{"", true},
		{"aurora-mysql", true},
		{"aurora-postgresql", true},
		{"postgres", true},
		{"mariadb", true},
	}
	for _, tt := range tests {
		t.Run(tt.engine, func(t *testing.T) {
			err := ValidateEngine(tt.engine)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEngine(%q) error = %v, wantErr %v", tt.engine, err, tt.wantErr)
			}
		})
	}
}

func TestIsSupportedEngine(t *testing.T) {
	if !IsSupportedEngine("mysql") {
		t.Error("mysql should be supported")
	}
	if !IsSupportedEngine("neptune") {
		t.Error("neptune should be supported")
	}
	if IsSupportedEngine("aurora-mysql") {
		t.Error("aurora-mysql should not be supported")
	}
	if IsSupportedEngine("") {
		t.Error("empty engine should not be supported")
	}
}

func TestValidateDBSnapshotIdentifier(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"valid simple", "mysnap", false},
		{"valid with hyphens", "my-snapshot-1", false},
		{"valid max length 255", strings.Repeat("a", 255), false},
		{"empty", "", true},
		{"too long 256", strings.Repeat("a", 256), true},
		{"starts with digit", "1snap", true},
		{"starts with hyphen", "-snap", true},
		{"ends with hyphen", "snap-", true},
		{"consecutive hyphens", "s--1", true},
		{"contains slash", "snap/1", true},
		{"contains colon", "snap:1", true},
		{"contains underscore", "snap_1", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDBSnapshotIdentifier(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDBSnapshotIdentifier(%q) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}

func TestSupportedMysqlVersionsConsistency(t *testing.T) {
	versions := SupportedMysqlVersions()
	if len(versions) == 0 {
		t.Error("SupportedMysqlVersions should not be empty")
	}
	// Every version in the list should be accepted by ValidateEngineVersion.
	for _, v := range versions {
		if err := ValidateEngineVersion("mysql", v.Version); err != nil {
			t.Errorf("version %q in SupportedMysqlVersions but rejected by ValidateEngineVersion: %v",
				v.Version, err)
		}
	}
	// No duplicate versions.
	seen := make(map[string]bool)
	for _, v := range versions {
		if seen[v.Version] {
			t.Errorf("duplicate MySQL version %q", v.Version)
		}
		seen[v.Version] = true
	}
}

func TestSupportedNeptuneVersionsConsistency(t *testing.T) {
	versions := SupportedNeptuneVersions()
	if len(versions) == 0 {
		t.Error("SupportedNeptuneVersions should not be empty")
	}
	// Every version in the list should be accepted by ValidateEngineVersion.
	for _, v := range versions {
		if err := ValidateEngineVersion("neptune", v.Version); err != nil {
			t.Errorf("version %q in SupportedNeptuneVersions but rejected by ValidateEngineVersion: %v",
				v.Version, err)
		}
	}
	// No duplicate versions.
	seen := make(map[string]bool)
	for _, v := range versions {
		if seen[v.Version] {
			t.Errorf("duplicate Neptune version %q", v.Version)
		}
		seen[v.Version] = true
	}
}

func TestAllEngineVersionsCompleteness(t *testing.T) {
	versions := allEngineVersions()
	if len(versions) != len(supportedMysqlVersions)+len(supportedNeptuneVersions) {
		t.Errorf("allEngineVersions returned %d entries, expected %d (mysql + neptune)",
			len(versions), len(supportedMysqlVersions)+len(supportedNeptuneVersions))
	}
}

func TestDefaultEngineVersion(t *testing.T) {
	mysqlDefault := DefaultEngineVersion("mysql")
	if mysqlDefault == "" {
		t.Error("DefaultEngineVersion('mysql') should not be empty")
	}
	// Should be the last entry in supportedMysqlVersions
	expected := supportedMysqlVersions[len(supportedMysqlVersions)-1].Version
	if mysqlDefault != expected {
		t.Errorf("DefaultEngineVersion('mysql') = %q, want %q (latest)", mysqlDefault, expected)
	}
	// Should pass ValidateEngineVersion
	if err := ValidateEngineVersion("mysql", mysqlDefault); err != nil {
		t.Errorf("default mysql version %q failed validation: %v", mysqlDefault, err)
	}

	neptuneDefault := DefaultEngineVersion("neptune")
	if neptuneDefault == "" {
		t.Error("DefaultEngineVersion('neptune') should not be empty")
	}
	expectedNeptune := supportedNeptuneVersions[len(supportedNeptuneVersions)-1].Version
	if neptuneDefault != expectedNeptune {
		t.Errorf("DefaultEngineVersion('neptune') = %q, want %q", neptuneDefault, expectedNeptune)
	}
	if err := ValidateEngineVersion("neptune", neptuneDefault); err != nil {
		t.Errorf("default neptune version %q failed validation: %v", neptuneDefault, err)
	}

	if DefaultEngineVersion("aurora-mysql") != "" {
		t.Error("DefaultEngineVersion('aurora-mysql') should return empty for unsupported engine")
	}
}
