package lambda

import (
	"net/http"

	"google.golang.org/protobuf/proto"
	svccommon "vorpalstacks/internal/common"

	pb "vorpalstacks/internal/pb/aws/lambda"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	"vorpalstacks/internal/utils/timeutils"
)

// getStore extracts the region from request headers and returns the
// full lambdaStore for that region. This is the sole entry point for
// store access in the admin handler layer.
func (h *AdminHandler) getStore(header http.Header) (*lambdaStore, error) {
	region := svccommon.GetRegionFromHeader(header)
	return h.service.getOrCreateLambdaStore(region), nil
}

// safeRuntime converts a store Runtime to a proto Runtime, falling back to
// nodejs22x when the value is not present in the proto enum map.
func safeRuntime(v lambdastore.Runtime) pb.Runtime {
	if val, ok := pb.Runtime_value[string(v)]; ok {
		return pb.Runtime(val)
	}
	return pb.Runtime_RUNTIME_NODEJS22X
}

// protoToStoreRuntime converts a proto Runtime enum to the canonical Lambda
// runtime string (e.g. pb.Runtime_RUNTIME_NODEJS22X → "nodejs22.x").
// Returns "" for unsupported or unmapped enum values so that ValidateRuntime
// rejects them with a clear error. Only currently-supported runtimes (per
// validRuntimes in validators.go) are mapped; EOL runtimes such as
// nodejs20.x, nodejs18.x, nodejs16.x, python3.9, python3.8, dotnet6,
// ruby3.2, go1.x are intentionally excluded.
var protoToStoreRuntimeMap = map[pb.Runtime]string{
	pb.Runtime_RUNTIME_NODEJS24X:      "nodejs24.x",
	pb.Runtime_RUNTIME_NODEJS22X:      "nodejs22.x",
	pb.Runtime_RUNTIME_PYTHON314:      "python3.14",
	pb.Runtime_RUNTIME_PYTHON313:      "python3.13",
	pb.Runtime_RUNTIME_PYTHON312:      "python3.12",
	pb.Runtime_RUNTIME_PYTHON311:      "python3.11",
	pb.Runtime_RUNTIME_PYTHON310:      "python3.10",
	pb.Runtime_RUNTIME_JAVA25:         "java25",
	pb.Runtime_RUNTIME_JAVA21:         "java21",
	pb.Runtime_RUNTIME_JAVA17:         "java17",
	pb.Runtime_RUNTIME_JAVA11:         "java11",
	pb.Runtime_RUNTIME_JAVA8AL2:       "java8.al2",
	pb.Runtime_RUNTIME_DOTNET10:       "dotnet10",
	pb.Runtime_RUNTIME_DOTNET8:        "dotnet8",
	pb.Runtime_RUNTIME_RUBY40:         "ruby4.0",
	pb.Runtime_RUNTIME_RUBY34:         "ruby3.4",
	pb.Runtime_RUNTIME_RUBY33:         "ruby3.3",
	pb.Runtime_RUNTIME_PROVIDEDAL2023: "provided.al2023",
	pb.Runtime_RUNTIME_PROVIDEDAL2:    "provided.al2",
}

func protoToStoreRuntime(r pb.Runtime) string {
	if s, ok := protoToStoreRuntimeMap[r]; ok {
		return s
	}
	return ""
}

// protoToPackageType converts a proto PackageType enum to the canonical
// Lambda package type string ("Zip" or "Image").
func protoToPackageType(pt pb.PackageType) string {
	switch pt {
	case pb.PackageType_PACKAGE_TYPE_IMAGE:
		return "Image"
	case pb.PackageType_PACKAGE_TYPE_ZIP:
		return "Zip"
	default:
		return "Zip"
	}
}

// safeState converts a store State to a proto State, defaulting to Active.
func safeState(v lambdastore.State) pb.State {
	if val, ok := pb.State_value[string(v)]; ok {
		return pb.State(val)
	}
	return pb.State_STATE_ACTIVE
}

// safeStateReasonCode maps a state reason code string to the proto enum,
// falling back to Idle when the value is unrecognised.
func safeStateReasonCode(v string) pb.StateReasonCode {
	if val, ok := pb.StateReasonCode_value[v]; ok {
		return pb.StateReasonCode(val)
	}
	return pb.StateReasonCode_STATE_REASON_CODE_IDLE
}

// safePackageType maps a package type string to the proto enum, falling back
// to Zip when the value is unrecognised.
func safePackageType(v string) pb.PackageType {
	if val, ok := pb.PackageType_value[v]; ok {
		return pb.PackageType(val)
	}
	return pb.PackageType_PACKAGE_TYPE_ZIP
}

// functionToProto converts a store Function to its proto representation,
// using safe enum mappers to avoid silent zero-value fallbacks.
func functionToProto(f *lambdastore.Function) *pb.FunctionConfiguration {
	pbFn := &pb.FunctionConfiguration{
		Functionname:    f.FunctionName,
		Functionarn:     f.FunctionArn,
		Runtime:         safeRuntime(f.Runtime),
		Role:            f.Role,
		Handler:         f.Handler,
		Codesize:        proto.Int64(f.CodeSize),
		Codesha256:      f.CodeSha256,
		Description:     f.Description,
		Timeout:         proto.Int32(f.Timeout),
		Memorysize:      proto.Int32(f.MemorySize),
		Lastmodified:    f.LastModified.Format(timeutils.ISO8601UTCFormat),
		Revisionid:      f.RevisionId,
		State:           safeState(f.State),
		Statereason:     f.StateReason,
		Statereasoncode: safeStateReasonCode(f.StateReasonCode),
		Packagetype:     safePackageType(f.PackageType),
	}
	if f.EphemeralStorage != nil {
		pbFn.Ephemeralstorage = &pb.EphemeralStorage{Size: f.EphemeralStorage.Size}
	}
	return pbFn
}
