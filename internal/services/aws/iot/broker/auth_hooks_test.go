package broker

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"

	mqttsrv "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
	vcrypto "vorpalstacks/internal/utils/crypto"
)

type mockPolicyLookup struct {
	policies      map[string]*iotstore.Policy
	principalPols map[string][]string
	getErr        error
	listErr       error
}

func (m *mockPolicyLookup) CreatePolicy(p *iotstore.Policy) (*iotstore.Policy, error) {
	if m.policies == nil {
		m.policies = make(map[string]*iotstore.Policy)
	}
	m.policies[p.PolicyName] = p
	return p, nil
}

func (m *mockPolicyLookup) GetPolicy(name string) (*iotstore.Policy, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if p, ok := m.policies[name]; ok {
		return p, nil
	}
	return nil, iotstore.ErrPolicyNotFound
}

func (m *mockPolicyLookup) DeletePolicy(name string) error {
	delete(m.policies, name)
	return nil
}

func (m *mockPolicyLookup) UpdatePolicy(p *iotstore.Policy) error {
	if m.policies == nil {
		m.policies = make(map[string]*iotstore.Policy)
	}
	m.policies[p.PolicyName] = p
	return nil
}

func (m *mockPolicyLookup) ListPolicies(opts storecommon.ListOptions) (*storecommon.ListResult[iotstore.Policy], error) {
	return &storecommon.ListResult[iotstore.Policy]{}, nil
}

func (m *mockPolicyLookup) AttachPolicyToPrincipal(policyName, principal string) error {
	if m.principalPols == nil {
		m.principalPols = make(map[string][]string)
	}
	m.principalPols[principal] = append(m.principalPols[principal], policyName)
	return nil
}

func (m *mockPolicyLookup) DetachPolicyFromPrincipal(policyName, principal string) error {
	pols := m.principalPols[principal]
	for i, p := range pols {
		if p == policyName {
			m.principalPols[principal] = append(pols[:i], pols[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockPolicyLookup) ListPrincipalsForPolicy(policyName string) ([]string, error) {
	return nil, nil
}

func (m *mockPolicyLookup) ListPoliciesForPrincipal(principal string) ([]string, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.principalPols[principal], nil
}

func generateTestCert(t *testing.T) *x509.Certificate {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("ecdsa.GenerateKey: %v", err)
	}
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test-device"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("x509.CreateCertificate: %v", err)
	}
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		t.Fatalf("x509.ParseCertificate: %v", err)
	}
	return cert
}

func TestFingerprintX509(t *testing.T) {
	cert := generateTestCert(t)
	id := vcrypto.FingerprintX509(cert)
	if len(id) != 64 {
		t.Errorf("ID length = %d, want 64", len(id))
	}
}

func TestFingerprintX509Deterministic(t *testing.T) {
	cert := generateTestCert(t)
	id1 := vcrypto.FingerprintX509(cert)
	id2 := vcrypto.FingerprintX509(cert)
	if id1 != id2 {
		t.Fatal("expected same cert to produce same ID")
	}
}

func TestFingerprintX509Distinct(t *testing.T) {
	cert1 := generateTestCert(t)
	cert2 := generateTestCert(t)
	id1 := vcrypto.FingerprintX509(cert1)
	id2 := vcrypto.FingerprintX509(cert2)
	if id1 == id2 {
		t.Fatal("expected distinct certs to produce distinct IDs")
	}
}

func TestGetPoliciesForPrincipalEmpty(t *testing.T) {
	store := &mockPolicyLookup{}
	result := getPoliciesForPrincipal(store, "nonexistent")
	if len(result) != 0 {
		t.Errorf("expected 0 policies, got %d", len(result))
	}
}

func TestGetPoliciesForPrincipalValid(t *testing.T) {
	store := &mockPolicyLookup{
		policies: map[string]*iotstore.Policy{
			"p1": {PolicyName: "p1", PolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"iot:Connect","Resource":"#"}]}`},
		},
		principalPols: map[string][]string{
			"cert-1": {"p1"},
		},
	}
	result := getPoliciesForPrincipal(store, "cert-1")
	if len(result) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(result))
	}
	if result[0].Version != "2012-10-17" {
		t.Errorf("version = %s", result[0].Version)
	}
}

func TestGetPoliciesForPrincipalSkipsBadJSON(t *testing.T) {
	store := &mockPolicyLookup{
		policies: map[string]*iotstore.Policy{
			"good": {PolicyName: "good", PolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"iot:*"}]}`},
			"bad":  {PolicyName: "bad", PolicyDocument: `{invalid`},
		},
		principalPols: map[string][]string{
			"cert-1": {"bad", "good"},
		},
	}
	result := getPoliciesForPrincipal(store, "cert-1")
	if len(result) != 1 {
		t.Fatalf("expected 1 (bad skipped), got %d", len(result))
	}
}

func TestGetPoliciesForPrincipalGetError(t *testing.T) {
	store := &mockPolicyLookup{
		policies: map[string]*iotstore.Policy{
			"p1": {PolicyName: "p1", PolicyDocument: `{"Version":"2012-10-17","Statement":[]}`},
		},
		principalPols: map[string][]string{
			"cert-1": {"p1"},
		},
		getErr: iotstore.ErrPolicyNotFound,
	}
	result := getPoliciesForPrincipal(store, "cert-1")
	if len(result) != 0 {
		t.Errorf("expected 0 on GetPolicy error, got %d", len(result))
	}
}

func TestGetPoliciesForPrincipalListError(t *testing.T) {
	store := &mockPolicyLookup{listErr: iotstore.ErrPolicyNotFound}
	result := getPoliciesForPrincipal(store, "cert-1")
	if len(result) != 0 {
		t.Errorf("expected 0 on list error, got %d", len(result))
	}
}

func TestNewAuthHook(t *testing.T) {
	b := &Broker{}
	hook := NewAuthHook(b, nil)
	if hook == nil {
		t.Fatal("expected non-nil hook")
	}
}

func TestAuthHookID(t *testing.T) {
	b := &Broker{}
	hook := NewAuthHook(b, nil)
	if hook.ID() != "iot-auth" {
		t.Errorf("ID = %s, want iot-auth", hook.ID())
	}
}

func TestAuthHookProvides(t *testing.T) {
	b := &Broker{}
	hook := NewAuthHook(b, nil)
	if !hook.Provides(mqttsrv.OnConnectAuthenticate) {
		t.Error("expected Provides true for OnConnectAuthenticate")
	}
	if !hook.Provides(mqttsrv.OnACLCheck) {
		t.Error("expected Provides true for OnACLCheck")
	}
	if !hook.Provides(mqttsrv.OnSessionEstablished) {
		t.Error("expected Provides true for OnSessionEstablished")
	}
	if !hook.Provides(mqttsrv.OnDisconnect) {
		t.Error("expected Provides true for OnDisconnect")
	}
}

func TestOnConnectAuthenticateNilProvider(t *testing.T) {
	b := &Broker{}
	hook := NewAuthHook(b, nil)
	cl := &mqttsrv.Client{}
	if hook.OnConnectAuthenticate(cl, packets.Packet{}) {
		t.Error("expected false for nil provider")
	}
}

func TestOnACLCheckNilProvider(t *testing.T) {
	b := &Broker{}
	hook := NewAuthHook(b, nil)
	cl := &mqttsrv.Client{}
	if hook.OnACLCheck(cl, "topic/test", true) {
		t.Error("expected false for nil provider")
	}
}
