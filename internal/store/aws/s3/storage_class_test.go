package s3

import "testing"

// The acceptance set must cover exactly the persisted constants: every
// constant is accepted, and the AWS enum's hardware- and backup-bound
// classes are not.
func TestIsValidStorageClassMatchesConstants(t *testing.T) {
	for _, class := range []ObjectStorageClass{
		StorageClassStandard,
		StorageClassReducedRedundancy,
		StorageClassGlacier,
		StorageClassStandardIA,
		StorageClassOneZoneIA,
		StorageClassIntelligentTiering,
		StorageClassGlacierIR,
		StorageClassDeepArchive,
	} {
		if !IsValidStorageClass(class) {
			t.Fatalf("persisted constant %s must be accepted", class)
		}
	}

	for _, class := range []ObjectStorageClass{
		"OUTPOSTS", "SNOW", "EXPRESS_ONEZONE",
		"FSX_OPENZFS", "FSX_ONTAP",
		"AWS_BACKUP_WARM", "AWS_BACKUP_LOW_COST_WARM",
		"", "not-a-class",
	} {
		if IsValidStorageClass(class) {
			t.Fatalf("hardware/backup-bound or unknown class %q must be rejected", class)
		}
	}
}
