package lambda

import (
	"reflect"
	"testing"
	"time"

	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
)

// atomicFunctionStore builds a FunctionStore over a fresh empty storage.
func atomicFunctionStore(t *testing.T) *FunctionStore {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("new region storage manager: %v", err)
	}
	st, err := mgr.GetStorage("us-east-1")
	if err != nil {
		t.Fatalf("get storage: %v", err)
	}
	return NewFunctionStore(st, "000000000000", "us-east-1")
}

// fullyPopulatedFunction returns a Function with every exported field
// non-zero, so the copy round-trip tests can detect a silently dropped
// field (a zero fixture member is indistinguishable from a dropped one,
// which is exactly how the historic services-side copier lost
// DeadLetterConfig, TracingConfig, SnapStart and the container IDs).
func fullyPopulatedFunction() *Function {
	reserved := int64(64)
	f := &Function{
		FunctionName:               "full",
		FunctionArn:                "arn:aws:lambda:us-east-1:000000000000:function:full",
		Runtime:                    RuntimeNodejs22X,
		Role:                       "arn:aws:iam::000000000000:role/lambda",
		Handler:                    "index.handler",
		CodeSize:                   128,
		CodeSha256:                 "sha256",
		CodeLocation:               "/data/us-east-1/code/full/$LATEST/code.zip",
		ImageUri:                   "012345678901.dkr.ecr.us-east-1.amazonaws.com/full:1",
		SourceCodeHash:             "source-hash",
		Description:                "description",
		Timeout:                    30,
		MemorySize:                 256,
		Publish:                    true,
		KMSKeyArn:                  "arn:aws:kms:us-east-1:000000000000:key/1",
		RevisionId:                 "revision-1",
		State:                      StateActive,
		StateReason:                "reason",
		StateReasonCode:            "code",
		LastUpdateStatus:           LastUpdateStatusSuccessful,
		LastUpdateReason:           "reason",
		LastUpdateStatusReason:     "reason",
		LastUpdateStatusReasonCode: "code",
		LastModified:               time.Date(2026, 9, 8, 1, 2, 3, 0, time.UTC),
		LastModifiedUser:           "user",
		PackageType:                "Zip",
		SigningProfileVersionArn:   "arn:aws:signer:signing-profile",
		SigningJobArn:              "arn:aws:signer:signing-job",
		CodeSigningConfigArn:       "arn:aws:lambda:code-signing-config",
		ContainerID:                "container-id",
		ContainerImageID:           "container-image-id",

		EphemeralStorage: &EphemeralStorage{Size: 1024},
		Architectures:    []string{"arm64"},
		VpcConfig: &VpcConfig{
			SubnetIds:               []string{"subnet-1"},
			SecurityGroupIds:        []string{"sg-1"},
			VpcId:                   "vpc-1",
			Ipv6AllowedForDualStack: true,
		},
		Environment: &Environment{
			Variables: map[string]string{"KEY": "value"},
			Error:     &EnvironmentError{ErrorCode: "code", Message: "message"},
		},
		DeadLetterConfig: &DeadLetterConfig{TargetArn: "arn:aws:sqs:us-east-1:000000000000:dlq"},
		TracingConfig:    &TracingConfig{Mode: "Active"},
		Layers:           []LayerReference{{Arn: "arn:aws:lambda:layer:l:1", CodeSize: 8}},
		Tags:             []types.Tag{{Key: "k", Value: "v"}},
		SnapStart:        &SnapStart{ApplyOn: "PublishedVersions", OptimizationStatus: "Successful"},
		UrlConfig: &FunctionUrlConfig{
			FunctionUrl:      "https://full.lambda-url.us-east-1.on.aws/",
			FunctionArn:      "arn:aws:lambda:us-east-1:000000000000:function:full",
			AuthType:         "AWS_IAM",
			InvokeMode:       "BUFFERED",
			CreationTime:     time.Date(2026, 9, 8, 1, 2, 3, 0, time.UTC),
			LastModifiedTime: time.Date(2026, 9, 8, 1, 2, 3, 0, time.UTC),
			Qualifier:        "prod",
			Cors: &CorsConfig{
				AllowCredentials: true,
				AllowHeaders:     []string{"x-header"},
				AllowMethods:     []string{"GET"},
				AllowOrigins:     []string{"https://origin"},
				ExposeHeaders:    []string{"x-expose"},
				MaxAge:           60,
			},
		},
		LoggingConfig: &LoggingConfig{
			LogFormat:           "JSON",
			ApplicationLogLevel: "INFO",
			SystemLogLevel:      "INFO",
			LogGroup:            "/aws/lambda/full",
		},
		ImageConfig: &ImageConfig{
			EntryPoint:       []string{"/entry"},
			Command:          []string{"cmd"},
			WorkingDirectory: "/work",
		},
		FileSystemConfigs: []FileSystemConfig{{Arn: "arn:aws:elasticfilesystem:fs-1", LocalMountPath: "/mnt"}},
		TenancyConfig:     &TenancyConfig{TenantIsolationMode: "tenant"},
		CapacityProviderConfig: map[string]interface{}{
			"capacity": float64(2),
		},
		RuntimeVersionConfig: &RuntimeVersionConfig{
			RuntimeVersionArn: "arn:aws:lambda:runtime-version",
			Error:             map[string]interface{}{"errorCode": "code"},
		},
		DurableConfig:       map[string]interface{}{"durable": true},
		ReservedConcurrency: &reserved,
		PolicyRevisionId:    "11111111-2222-3333-4444-555555555555",
		Versions:            []Version{*fullyPopulatedVersion()},
		Aliases:             []Alias{*fullyPopulatedAlias()},
		Policies: []FunctionPolicy{
			{
				Id:        "stmt",
				Effect:    "Allow",
				Principal: "*",
				Action:    "lambda:InvokeFunction",
				Resource:  "arn",
				Condition: map[string]interface{}{"k": "v"},
				Raw:       `{"Sid":"stmt","Effect":"Allow","Principal":"*","Action":"lambda:InvokeFunction","Resource":"arn"}`,
			},
		},
		ProvisionedConcurrency: []ProvisionedConcurrencyConfig{{
			FunctionName:                             "full",
			FunctionArn:                              "arn",
			Qualifier:                                "1",
			AllocatedProvisionedConcurrentExecutions: 1,
			AvailableProvisionedConcurrentExecutions: 1,
			RequestedProvisionedConcurrentExecutions: 1,
			Status:                                   "READY",
			StatusReason:                             "reason",
			LastModified:                             time.Date(2026, 9, 8, 1, 2, 3, 0, time.UTC),
		}},
		EventInvokeConfigs: []EventInvokeConfig{{
			FunctionName: "full",
			Qualifier:    "1",
			LastModified: time.Date(2026, 9, 8, 1, 2, 3, 0, time.UTC),
			DestinationConfig: &DestinationConfig{
				OnSuccess: &OnSuccess{Destination: "arn:aws:sqs:us-east-1:000000000000:ok"},
				OnFailure: &OnFailure{Destination: "arn:aws:sqs:us-east-1:000000000000:bad"},
			},
			MaximumEventAgeInSeconds: 3600,
			MaximumRetryAttempts:     2,
		}},
	}
	f.rebuildIndexes()
	return f
}

func fullyPopulatedVersion() *Version {
	return &Version{
		Version:                  "7",
		FunctionArn:              "arn:aws:lambda:us-east-1:000000000000:function:full:7",
		Runtime:                  RuntimeNodejs22X,
		Role:                     "arn:aws:iam::000000000000:role/lambda",
		Handler:                  "index.handler",
		CodeSize:                 128,
		CodeSha256:               "sha256",
		CodeLocation:             "/data/us-east-1/code/full/7/code.zip",
		ImageUri:                 "012345678901.dkr.ecr.us-east-1.amazonaws.com/full:1",
		Description:              "description",
		Timeout:                  30,
		MemorySize:               256,
		EphemeralStorage:         &EphemeralStorage{Size: 1024},
		Architectures:            []string{"arm64"},
		KMSKeyArn:                "arn:aws:kms:us-east-1:000000000000:key/1",
		RevisionId:               "revision-7",
		State:                    StateActive,
		StateReason:              "reason",
		StateReasonCode:          "code",
		LastUpdateStatus:         LastUpdateStatusSuccessful,
		LastModified:             time.Date(2026, 9, 8, 1, 2, 3, 0, time.UTC),
		VpcConfig:                &VpcConfig{SubnetIds: []string{"subnet-1"}, SecurityGroupIds: []string{"sg-1"}, VpcId: "vpc-1", Ipv6AllowedForDualStack: true},
		Environment:              &Environment{Variables: map[string]string{"KEY": "value"}, Error: &EnvironmentError{ErrorCode: "code", Message: "message"}},
		DeadLetterConfig:         &DeadLetterConfig{TargetArn: "arn:aws:sqs:us-east-1:000000000000:dlq"},
		TracingConfig:            &TracingConfig{Mode: "Active"},
		Layers:                   []LayerReference{{Arn: "arn:aws:lambda:layer:l:1", CodeSize: 8}},
		SnapStart:                &SnapStart{ApplyOn: "PublishedVersions", OptimizationStatus: "Successful"},
		PackageType:              "Zip",
		SigningProfileVersionArn: "arn:aws:signer:signing-profile",
		SigningJobArn:            "arn:aws:signer:signing-job",
		LoggingConfig:            &LoggingConfig{LogFormat: "JSON", ApplicationLogLevel: "INFO", SystemLogLevel: "INFO", LogGroup: "/aws/lambda/full"},
		ImageConfig:              &ImageConfig{EntryPoint: []string{"/entry"}, Command: []string{"cmd"}, WorkingDirectory: "/work"},
		FileSystemConfigs:        []FileSystemConfig{{Arn: "arn:aws:elasticfilesystem:fs-1", LocalMountPath: "/mnt"}},
		ContainerID:              "container-id",
		ContainerImageID:         "container-image-id",
	}
}

func fullyPopulatedAlias() *Alias {
	return &Alias{
		Name:            "prod",
		AliasArn:        "arn:aws:lambda:us-east-1:000000000000:function:full:prod",
		FunctionVersion: "7",
		Description:     "alias",
		FunctionName:    "full",
		RevisionId:      "revision-alias",
		RoutingConfig:   &RoutingConfig{AdditionalVersionWeights: map[string]float64{"6": 0.5}},
	}
}

// requireFullyPopulated walks every exported field and fails on a zero
// member, keeping the fixture strong enough that a dropped field in a
// copy necessarily shows up as an inequality.
func requireFullyPopulated(t *testing.T, v reflect.Value, path string) {
	t.Helper()
	switch v.Kind() {
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			if v.Interface().(time.Time).IsZero() {
				t.Fatalf("%s is zero — fixture incomplete", path)
			}
			return
		}
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).PkgPath != "" {
				continue
			}
			requireFullyPopulated(t, v.Field(i), path+"."+v.Type().Field(i).Name)
		}
	default:
		if v.IsZero() {
			t.Fatalf("%s is zero — fixture incomplete", path)
		}
	}
}

// TestFunctionDeepCopyRoundTripsEveryField pins the single Function
// deep-copy: a fully populated function copies equal in every field and
// stays detached — mutating the original's nested pointers, slices and
// maps never reaches the copy.
func TestFunctionDeepCopyRoundTripsEveryField(t *testing.T) {
	f := fullyPopulatedFunction()
	requireFullyPopulated(t, reflect.ValueOf(*f).FieldByName("FunctionName"), "FunctionName") // smoke: walk below is the real guard
	orig := reflect.ValueOf(*f)
	for i := 0; i < orig.NumField(); i++ {
		if orig.Type().Field(i).PkgPath != "" {
			continue
		}
		requireFullyPopulated(t, orig.Field(i), "Function."+orig.Type().Field(i).Name)
	}

	copyFn := f.DeepCopy()

	// The copy must equal the fully populated original in every exported
	// field — a field DeepCopy omits stays zero here and fails the walk,
	// which is what the package comment promises. Unexported index maps
	// are rebuilt lazily and are skipped.
	for i := 0; i < orig.NumField(); i++ {
		if orig.Type().Field(i).PkgPath != "" {
			continue
		}
		if !reflect.DeepEqual(orig.Field(i).Interface(), reflect.ValueOf(*copyFn).Field(i).Interface()) {
			t.Fatalf("DeepCopy lost or altered field %s", orig.Type().Field(i).Name)
		}
	}

	// Detachment probes: every mutation of the original must leave the
	// copy holding the pre-mutation values.
	f.VpcConfig.SubnetIds[0] = "mutated"
	f.Environment.Variables["KEY"] = "mutated"
	f.Environment.Error.Message = "mutated"
	f.Tags[0].Key = "mutated"
	f.UrlConfig.Cors.AllowOrigins[0] = "mutated"
	f.UrlConfig.Cors.AllowHeaders[0] = "mutated"
	f.Versions[0].DeadLetterConfig.TargetArn = "mutated"
	f.Versions[0].TracingConfig.Mode = "mutated"
	f.Versions[0].SnapStart.ApplyOn = "mutated"
	f.Versions[0].ContainerID = "mutated"
	f.Aliases[0].RoutingConfig.AdditionalVersionWeights["6"] = 9
	f.Policies[0].Condition["k"] = "mutated"
	f.EventInvokeConfigs[0].DestinationConfig.OnFailure.Destination = "mutated"
	f.CapacityProviderConfig["capacity"] = "mutated"
	f.RuntimeVersionConfig.Error["errorCode"] = "mutated"
	*f.ReservedConcurrency = 1
	f.ImageConfig.EntryPoint[0] = "mutated"
	f.Layers[0] = LayerReference{}
	f.FileSystemConfigs[0].LocalMountPath = "mutated"
	f.Architectures[0] = "mutated"

	if got := copyFn.VpcConfig.SubnetIds[0]; got != "subnet-1" {
		t.Fatalf("copy shares VpcConfig.SubnetIds: %q", got)
	}
	if got := copyFn.Environment.Variables["KEY"]; got != "value" {
		t.Fatalf("copy shares Environment.Variables: %q", got)
	}
	if got := copyFn.Environment.Error.Message; got != "message" {
		t.Fatalf("copy shares Environment.Error: %q", got)
	}
	if got := copyFn.Tags[0].Key; got != "k" {
		t.Fatalf("copy shares Tags: %q", got)
	}
	if got := copyFn.UrlConfig.Cors.AllowOrigins[0]; got != "https://origin" {
		t.Fatalf("copy shares UrlConfig.Cors.AllowOrigins: %q", got)
	}
	if got := copyFn.UrlConfig.Cors.AllowHeaders[0]; got != "x-header" {
		t.Fatalf("copy shares UrlConfig.Cors.AllowHeaders: %q", got)
	}
	if got := copyFn.Versions[0].DeadLetterConfig.TargetArn; got != "arn:aws:sqs:us-east-1:000000000000:dlq" {
		t.Fatalf("copy shares nested version DeadLetterConfig: %q", got)
	}
	if got := copyFn.Versions[0].TracingConfig.Mode; got != "Active" {
		t.Fatalf("copy shares nested version TracingConfig: %q", got)
	}
	if got := copyFn.Versions[0].SnapStart.ApplyOn; got != "PublishedVersions" {
		t.Fatalf("copy shares nested version SnapStart: %q", got)
	}
	if got := copyFn.Versions[0].ContainerID; got != "container-id" {
		t.Fatalf("copy shares nested version ContainerID: %q", got)
	}
	if got := copyFn.Aliases[0].RoutingConfig.AdditionalVersionWeights["6"]; got != 0.5 {
		t.Fatalf("copy shares alias routing weights: %v", got)
	}
	if got := copyFn.Policies[0].Condition["k"]; got != "v" {
		t.Fatalf("copy shares policy condition: %v", got)
	}
	if got := copyFn.EventInvokeConfigs[0].DestinationConfig.OnFailure.Destination; got != "arn:aws:sqs:us-east-1:000000000000:bad" {
		t.Fatalf("copy shares event invoke destinations: %q", got)
	}
	if got := copyFn.CapacityProviderConfig["capacity"]; got != float64(2) {
		t.Fatalf("copy shares CapacityProviderConfig map: %v", got)
	}
	if got := copyFn.RuntimeVersionConfig.Error["errorCode"]; got != "code" {
		t.Fatalf("copy shares RuntimeVersionConfig.Error map: %v", got)
	}
	if got := *copyFn.ReservedConcurrency; got != 64 {
		t.Fatalf("copy shares ReservedConcurrency: %d", got)
	}
	if got := copyFn.ImageConfig.EntryPoint[0]; got != "/entry" {
		t.Fatalf("copy shares ImageConfig.EntryPoint: %q", got)
	}
	if got := copyFn.Layers[0].Arn; got != "arn:aws:lambda:layer:l:1" {
		t.Fatalf("copy shares Layers: %q", got)
	}
	if got := copyFn.FileSystemConfigs[0].LocalMountPath; got != "/mnt" {
		t.Fatalf("copy shares FileSystemConfigs: %q", got)
	}
	if got := copyFn.Architectures[0]; got != "arm64" {
		t.Fatalf("copy shares Architectures: %q", got)
	}
	if copyFn.versionsByNum["7"].ContainerID != "container-id" {
		t.Fatalf("copy index points at shared version storage")
	}
}

// TestDeepCopyFileSystemConfigsS3Files pins the pointer isolation of the
// S3FilesConfig member: the copy must not alias the source's struct.
func TestDeepCopyFileSystemConfigsS3Files(t *testing.T) {
	src := []FileSystemConfig{{
		Arn:            "arn:aws:s3files:us-east-1:123456789012:file-system/fs-0000000000000000000000000000000/access-point/fsap-0123456789abcdef0",
		LocalMountPath: "/mnt/s3",
		S3FilesConfig:  &S3FilesConfig{DirectS3Read: DirectS3ReadEnabled},
	}}
	dst := deepCopyFileSystemConfigs(src)
	if dst[0].S3FilesConfig == src[0].S3FilesConfig {
		t.Fatalf("copy shares the S3FilesConfig pointer")
	}
	src[0].S3FilesConfig.DirectS3Read = DirectS3ReadDisabled
	if got := dst[0].S3FilesConfig.DirectS3Read; got != DirectS3ReadEnabled {
		t.Fatalf("copy shares S3FilesConfig state: %q", got)
	}
}

// TestVersionDeepCopyRoundTripsEveryField pins the single Version
// deep-copy the same way, with the historically dropped members
// (DeadLetterConfig, TracingConfig, SnapStart, ContainerID,
// ContainerImageID) explicitly probed.
func TestVersionDeepCopyRoundTripsEveryField(t *testing.T) {
	v := fullyPopulatedVersion()
	orig := reflect.ValueOf(*v)
	for i := 0; i < orig.NumField(); i++ {
		if orig.Type().Field(i).PkgPath != "" {
			continue
		}
		requireFullyPopulated(t, orig.Field(i), "Version."+orig.Type().Field(i).Name)
	}

	copyVer := v.DeepCopy()

	v.DeadLetterConfig.TargetArn = "mutated"
	v.TracingConfig.Mode = "mutated"
	v.SnapStart.ApplyOn = "mutated"
	v.ContainerID = "mutated"
	v.ContainerImageID = "mutated"
	v.Environment.Error.Message = "mutated"
	v.VpcConfig.SubnetIds[0] = "mutated"
	v.ImageConfig.Command[0] = "mutated"

	if copyVer.DeadLetterConfig.TargetArn != "arn:aws:sqs:us-east-1:000000000000:dlq" {
		t.Fatalf("copy shares DeadLetterConfig: %q", copyVer.DeadLetterConfig.TargetArn)
	}
	if copyVer.TracingConfig.Mode != "Active" {
		t.Fatalf("copy shares TracingConfig: %q", copyVer.TracingConfig.Mode)
	}
	if copyVer.SnapStart.ApplyOn != "PublishedVersions" {
		t.Fatalf("copy shares SnapStart: %q", copyVer.SnapStart.ApplyOn)
	}
	if copyVer.ContainerID != "container-id" {
		t.Fatalf("copy shares ContainerID: %q", copyVer.ContainerID)
	}
	if copyVer.ContainerImageID != "container-image-id" {
		t.Fatalf("copy shares ContainerImageID: %q", copyVer.ContainerImageID)
	}
	if copyVer.Environment.Error.Message != "message" {
		t.Fatalf("copy shares Environment.Error: %q", copyVer.Environment.Error.Message)
	}
	if copyVer.VpcConfig.SubnetIds[0] != "subnet-1" {
		t.Fatalf("copy shares VpcConfig: %q", copyVer.VpcConfig.SubnetIds[0])
	}
	if copyVer.ImageConfig.Command[0] != "cmd" {
		t.Fatalf("copy shares ImageConfig.Command: %q", copyVer.ImageConfig.Command[0])
	}
}

// snapshotMetadataFields are assigned by the publish operation itself
// (identity, precondition and status members) rather than copied from
// the function's current configuration.
var snapshotMetadataFields = map[string]bool{
	"Version":          true,
	"FunctionArn":      true,
	"RevisionId":       true,
	"Description":      true,
	"State":            true,
	"StateReason":      true,
	"StateReasonCode":  true,
	"LastUpdateStatus": true,
	"LastModified":     true,
}

// TestPublishSnapshotCopiesEveryModelledMember pins the field-by-field
// version construction inside PublishVersionAtomically: every Version
// member that shares its name with a Function member must carry the
// function's value, so a newly modelled member cannot silently miss the
// snapshot.
func TestPublishSnapshotCopiesEveryModelledMember(t *testing.T) {
	s := atomicFunctionStore(t)
	f := fullyPopulatedFunction()
	f.FunctionName = "snapshot"
	f.Versions = nil
	f.rebuildIndexes()
	if _, err := s.Create(f); err != nil {
		t.Fatalf("create fixture function: %v", err)
	}

	version, err := s.PublishVersionAtomically("snapshot", "", "")
	if err != nil {
		t.Fatalf("publish: %v", err)
	}

	versionVal := reflect.ValueOf(*version)
	functionVal := reflect.ValueOf(*f)
	versionType := versionVal.Type()
	functionType := functionVal.Type()
	copied := 0
	for i := 0; i < versionType.NumField(); i++ {
		name := versionType.Field(i).Name
		if versionType.Field(i).PkgPath != "" || snapshotMetadataFields[name] {
			continue
		}
		fj, ok := functionType.FieldByName(name)
		if !ok || fj.PkgPath != "" {
			t.Fatalf("Version member %s has no exported Function counterpart — extend this test with its mapping", name)
		}
		got := versionVal.Field(i).Interface()
		want := functionVal.FieldByName(name).Interface()
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("snapshot member %s = %v, want the function's %v", name, got, want)
		}
		copied++
	}
	if copied < 20 {
		t.Fatalf("only %d snapshot members verified — the shared-name mapping is incomplete", copied)
	}
}
