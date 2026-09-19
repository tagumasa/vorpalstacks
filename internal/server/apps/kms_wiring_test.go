package apps

import (
	"reflect"
	"testing"

	"vorpalstacks/internal/core/storage"
	svckinesis "vorpalstacks/internal/services/aws/kinesis"
	svckms "vorpalstacks/internal/services/aws/kms"
	svckmshsm "vorpalstacks/internal/services/aws/kms/hsm"
	svcsqs "vorpalstacks/internal/services/aws/sqs"
	storesqs "vorpalstacks/internal/store/aws/sqs"
)

// newKMSWiringTestService builds a KMS service over a temporary HSM
// backend. The checker factory only wraps the service pointer, so the
// wiring test needs no keys or stores.
func newKMSWiringTestService(t *testing.T) *svckms.KMSService {
	t.Helper()
	backend, err := svckmshsm.NewPersistentBackend(t.TempDir())
	if err != nil {
		t.Fatalf("hsm backend: %v", err)
	}
	return svckms.NewKMSService("000000000000", "us-east-1", backend)
}

// checkerWired reports whether a service carries a KMS key checker. Each
// service keeps the checker unexported, so the boot wiring test reads the
// field through reflection rather than widening the services' APIs for
// one assertion.
func checkerWired(service any) bool {
	v := reflect.ValueOf(service).Elem().FieldByName("kmsChecker")
	return v.IsValid() && !v.IsNil()
}

// TestWireKMSCheckersCoversEveryEnableCombination pins the boot wiring
// contract: checker injection targets exactly the services that were
// constructed. A Kinesis-disabled deployment must boot (the Kinesis leg
// must not touch a nil service), a SQS-disabled deployment must still arm
// Kinesis key validation on Kinesis (the wiring must not depend on
// another service's initialiser having run), and a KMS-disabled
// deployment injects nothing.
func TestWireKMSCheckersCoversEveryEnableCombination(t *testing.T) {
	kms := newKMSWiringTestService(t)

	// KMS and SQS enabled, Kinesis disabled.
	regional, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	sqsStore := storesqs.NewSQSStore(regional, "000000000000", "us-east-1", "http://127.0.0.1:50080")
	st := &serviceState{
		kmsService: kms,
		sqsService: svcsqs.NewSQSServiceWithStore(sqsStore),
	}
	wireKMSCheckers(st) // must not panic on the nil Kinesis service
	if !checkerWired(st.sqsService) {
		t.Fatal("SQS service constructed alongside KMS carries no key checker")
	}

	// KMS and Kinesis enabled, SQS disabled.
	st = &serviceState{
		kmsService:     kms,
		kinesisService: svckinesis.NewKinesisService("000000000000"),
	}
	wireKMSCheckers(st)
	if !checkerWired(st.kinesisService) {
		t.Fatal("Kinesis service constructed alongside KMS carries no key checker")
	}

	// KMS disabled — no target is touched.
	st = &serviceState{
		sqsService:     svcsqs.NewSQSServiceWithStore(sqsStore),
		kinesisService: svckinesis.NewKinesisService("000000000000"),
	}
	wireKMSCheckers(st)
	if checkerWired(st.kinesisService) || checkerWired(st.sqsService) {
		t.Fatal("checker injected without a constructed KMS service")
	}
}
