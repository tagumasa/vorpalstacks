package broker

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

	"vorpalstacks/internal/services/aws/iot/ca"
	"vorpalstacks/internal/services/aws/iot/policy"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// AuthProvider supplies the store and CA needed by the authentication hook.
type AuthProvider interface {
	GetStore() iotstore.IotStoreInterface
	GetCA() *ca.CertificateAuthority
}

// AuthHook validates client certificates and enforces IoT policy-based ACLs
// on publish and subscribe operations.
type AuthHook struct {
	mqtt.HookBase
	provider AuthProvider
}

// NewAuthHook creates an authentication hook backed by the given provider.
func NewAuthHook(provider AuthProvider) *AuthHook {
	return &AuthHook{provider: provider}
}

// ID returns the hook identifier.
func (h *AuthHook) ID() string { return "iot-auth" }

// Provides declares which hook events this handler supports.
func (h *AuthHook) Provides(b byte) bool {
	return b == mqtt.OnConnectAuthenticate || b == mqtt.OnACLCheck
}

// OnConnectAuthenticate validates the client certificate against the IoT
// certificate store and evaluates iot:Connect policy before accepting the
// connection.
func (h *AuthHook) OnConnectAuthenticate(cl *mqtt.Client, _ packets.Packet) bool {
	if h.provider == nil {
		return false
	}

	certID := extractCertificateID(cl)
	if certID == "" {
		return false
	}

	store := h.provider.GetStore()
	if store == nil {
		return false
	}

	cert, err := store.GetCertificate(certID)
	if err != nil {
		return false
	}

	if cert.Status != "ACTIVE" {
		return false
	}

	policies := getPoliciesForPrincipal(store, certID)
	if len(policies) == 0 {
		return false
	}

	allowed, err := policy.Evaluate(&policy.EvaluateRequest{
		Policies: policies,
		Action:   "iot:Connect",
		Resource: cl.ID,
		ClientID: cl.ID,
	})
	if err != nil || !allowed {
		return false
	}

	cl.Properties.Username = []byte(certID)
	return true
}

// OnACLCheck evaluates the client's attached IoT policies to determine
// whether the requested topic operation is permitted.
func (h *AuthHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	if h.provider == nil {
		return false
	}

	certID := string(cl.Properties.Username)
	if certID == "" {
		return false
	}

	action := "iot:Subscribe"
	if write {
		action = "iot:Publish"
	}

	store := h.provider.GetStore()
	if store == nil {
		return false
	}

	policies := getPoliciesForPrincipal(store, certID)
	if len(policies) == 0 {
		return false
	}

	allowed, err := policy.Evaluate(&policy.EvaluateRequest{
		Policies: policies,
		Action:   action,
		Resource: topic,
		ClientID: cl.ID,
		Topic:    topic,
	})
	if err != nil {
		return false
	}

	return allowed
}

// extractCertificateID extracts the SHA-256 fingerprint of the peer
// certificate from the TLS connection state.
func extractCertificateID(cl *mqtt.Client) string {
	conn := cl.Net.Conn
	if conn == nil {
		return ""
	}

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		return ""
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return ""
	}

	derBytes := state.PeerCertificates[0].Raw
	fingerprint := sha256.Sum256(derBytes)
	return hex.EncodeToString(fingerprint[:])
}

// getPoliciesForPrincipal retrieves all policies and parses them for
// evaluation. In a full implementation this would filter by attachment.
func getPoliciesForPrincipal(store iotstore.IotStoreInterface, certID string) []*policy.PolicyVersion {
	policyNames, err := store.ListPoliciesForPrincipal(certID)
	if err != nil || len(policyNames) == 0 {
		return nil
	}

	var result []*policy.PolicyVersion
	for _, name := range policyNames {
		p, err := store.GetPolicy(name)
		if err != nil {
			continue
		}
		pv, err := policy.ParsePolicyVersion([]byte(p.PolicyDocument))
		if err != nil {
			continue
		}
		result = append(result, pv)
	}
	return result
}

// GetCertificateIDFromPEM computes the SHA-256 fingerprint of a PEM-encoded
// certificate, matching the format used by AWS IoT for certificate IDs.
func GetCertificateIDFromPEM(certPEM string) string {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return ""
	}
	fingerprint := sha256.Sum256(block.Bytes)
	return hex.EncodeToString(fingerprint[:])
}

// GetCertificateIDFromX509 computes the SHA-256 fingerprint of an X.509
// certificate in DER form.
func GetCertificateIDFromX509(cert *x509.Certificate) string {
	fingerprint := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(fingerprint[:])
}
