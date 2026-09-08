package lambda

import (
	"net/http"
	"vorpalstacks/internal/common/defaults"

	"google.golang.org/protobuf/proto"

	pb "vorpalstacks/internal/pb/aws/lambda"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	"vorpalstacks/internal/utils/timeutils"
)

// getStore extracts the region from request headers and returns the
// full lambdaStore for that region. This is the sole entry point for
// store access in the admin handler layer; a storage-layer failure
// propagates so the handler answers an internal error instead of the
// cores dereferencing a nil store.
func (h *AdminHandler) getStore(header http.Header) (*lambdaStore, error) {
	region := defaults.GetRegionFromHeader(header)
	return h.service.getOrCreateLambdaStoreE(region)
}

// safeRuntime converts a store Runtime to a proto Runtime, falling back to
// nodejs22x when the value has no proto enum mapping (deprecated runtimes).
// The reverse lookup runs against protoToStoreRuntimeMap so the response
// direction can never drift from the request direction; pb.Runtime_value is
// keyed by enum name, not by the runtime string, so it cannot serve here.
func safeRuntime(v lambdastore.Runtime) pb.Runtime {
	for pbRuntime, storeRuntime := range protoToStoreRuntimeMap {
		if storeRuntime == string(v) {
			return pbRuntime
		}
	}
	return pb.Runtime_RUNTIME_NODEJS22X
}

// protoToStoreRuntime converts a proto Runtime enum to the canonical Lambda
// runtime string (e.g. pb.Runtime_RUNTIME_NODEJS22X → "nodejs22.x").
// Returns "" for unsupported or unmapped enum values so that the create and
// update cores reject them with a clear error. Values come from the store's
// CurrentRuntimes single source and the agreement is pinned by unit test;
// EOL runtimes such as nodejs20.x, nodejs18.x, nodejs16.x, python3.9,
// python3.8, dotnet6, ruby3.2, go1.x are intentionally excluded.
var protoToStoreRuntimeMap = map[pb.Runtime]string{
	pb.Runtime_RUNTIME_NODEJS24X:      string(lambdastore.RuntimeNodejs24X),
	pb.Runtime_RUNTIME_NODEJS22X:      string(lambdastore.RuntimeNodejs22X),
	pb.Runtime_RUNTIME_PYTHON314:      string(lambdastore.RuntimePython314),
	pb.Runtime_RUNTIME_PYTHON313:      string(lambdastore.RuntimePython313),
	pb.Runtime_RUNTIME_PYTHON312:      string(lambdastore.RuntimePython312),
	pb.Runtime_RUNTIME_PYTHON311:      string(lambdastore.RuntimePython311),
	pb.Runtime_RUNTIME_PYTHON310:      string(lambdastore.RuntimePython310),
	pb.Runtime_RUNTIME_JAVA25:         string(lambdastore.RuntimeJava25),
	pb.Runtime_RUNTIME_JAVA21:         string(lambdastore.RuntimeJava21),
	pb.Runtime_RUNTIME_JAVA17:         string(lambdastore.RuntimeJava17),
	pb.Runtime_RUNTIME_JAVA11:         string(lambdastore.RuntimeJava11),
	pb.Runtime_RUNTIME_JAVA17AL2023:   string(lambdastore.RuntimeJava17Al2023),
	pb.Runtime_RUNTIME_JAVA11AL2023:   string(lambdastore.RuntimeJava11Al2023),
	pb.Runtime_RUNTIME_JAVA8AL2023:    string(lambdastore.RuntimeJava8Al2023),
	pb.Runtime_RUNTIME_JAVA8AL2:       string(lambdastore.RuntimeJava8Al2),
	pb.Runtime_RUNTIME_DOTNET10:       string(lambdastore.RuntimeDotnet10),
	pb.Runtime_RUNTIME_DOTNET8:        string(lambdastore.RuntimeDotnet8),
	pb.Runtime_RUNTIME_RUBY40:         string(lambdastore.RuntimeRuby40),
	pb.Runtime_RUNTIME_RUBY34:         string(lambdastore.RuntimeRuby34),
	pb.Runtime_RUNTIME_RUBY33:         string(lambdastore.RuntimeRuby33),
	pb.Runtime_RUNTIME_PROVIDEDAL2023: string(lambdastore.RuntimeProvidedAl2023),
	pb.Runtime_RUNTIME_PROVIDEDAL2:    string(lambdastore.RuntimeProvidedAl2),
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
		Functionname:    proto.String(f.FunctionName),
		Functionarn:     proto.String(f.FunctionArn),
		Runtime:         safeRuntime(f.Runtime),
		Role:            proto.String(f.Role),
		Handler:         proto.String(f.Handler),
		Codesize:        proto.Int64(f.CodeSize),
		Codesha256:      proto.String(f.CodeSha256),
		Description:     proto.String(f.Description),
		Timeout:         proto.Int32(f.Timeout),
		Memorysize:      proto.Int32(f.MemorySize),
		Lastmodified:    proto.String(f.LastModified.Format(timeutils.ISO8601UTCFormat)),
		Revisionid:      proto.String(f.RevisionId),
		State:           safeState(f.State),
		Statereason:     proto.String(f.StateReason),
		Statereasoncode: safeStateReasonCode(f.StateReasonCode),
		Packagetype:     safePackageType(f.PackageType),
	}
	if f.EphemeralStorage != nil {
		pbFn.Ephemeralstorage = &pb.EphemeralStorage{Size: f.EphemeralStorage.Size}
	}
	return pbFn
}
