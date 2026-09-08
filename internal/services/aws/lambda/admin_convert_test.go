package lambda

import (
	"testing"

	pb "vorpalstacks/internal/pb/aws/lambda"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// TestProtoToStoreRuntimeMapAgreesWithCurrentRuntimes pins the admin
// console's runtime surface to the store single source: every currently
// supported runtime is creatable from the console and the map maps
// nothing the validators would reject. A runtime added to
// CurrentRuntimes without a proto enum mapping fails here, closing the
// drift this test was introduced for.
func TestProtoToStoreRuntimeMapAgreesWithCurrentRuntimes(t *testing.T) {
	current := make(map[string]bool, len(lambdastore.CurrentRuntimes))
	for _, r := range lambdastore.CurrentRuntimes {
		current[string(r)] = true
	}

	mapped := make(map[string]bool, len(protoToStoreRuntimeMap))
	for _, v := range protoToStoreRuntimeMap {
		mapped[v] = true
	}

	for r := range current {
		if !mapped[r] {
			t.Errorf("runtime %q is accepted by the API but missing from protoToStoreRuntimeMap — the console cannot create functions with it", r)
		}
	}
	for r := range mapped {
		if !current[r] {
			t.Errorf("runtime %q is mapped by protoToStoreRuntimeMap but is not a currently supported runtime", r)
		}
	}
}

// TestSafeRuntimeCoversCurrentRuntimes pins the response direction of the
// console plane: a function stored with any current runtime echoes back
// that runtime instead of the safeRuntime fallback.
func TestSafeRuntimeCoversCurrentRuntimes(t *testing.T) {
	for _, r := range lambdastore.CurrentRuntimes {
		if got := safeRuntime(r); got == pb.Runtime_RUNTIME_NODEJS22X && string(r) != "nodejs22.x" {
			t.Errorf("runtime %q falls back to the nodejs22.x default in safeRuntime", r)
		}
	}
}
