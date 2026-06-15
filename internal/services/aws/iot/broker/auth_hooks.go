package broker

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"log/slog"

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
		slog.Warn("mqtt auth denied", "reason", "no provider")
		return false
	}

	certID := extractCertificateID(cl)
	if certID == "" {
		slog.Warn("mqtt auth denied", "client_id", cl.ID, "reason", "no certificate")
		return false
	}

	store := h.provider.GetStore()
	if store == nil {
		slog.Warn("mqtt auth denied", "cert_id", certID, "reason", "no store")
		return false
	}

	cert, err := store.GetCertificate(certID)
	if err != nil {
		slog.Warn("mqtt auth denied", "cert_id", certID, "reason", "certificate lookup failed", "error", err)
		return false
	}

	if cert.Status != "ACTIVE" {
		slog.Warn("mqtt auth denied", "cert_id", certID, "reason", "certificate not active", "status", cert.Status)
		return false
	}

	policies := getPoliciesForPrincipal(store, certID)
	if len(policies) == 0 {
		slog.Warn("mqtt auth denied", "cert_id", certID, "reason", "no policies attached")
		return false
	}

	allowed, err := policy.Evaluate(&policy.EvaluateRequest{
		Policies: policies,
		Action:   "iot:Connect",
		Resource: cl.ID,
		ClientID: cl.ID,
	})
	if err != nil {
		slog.Warn("mqtt auth denied", "cert_id", certID, "reason", "policy evaluation error", "error", err)
		return false
	}
	if !allowed {
		slog.Warn("mqtt auth denied", "cert_id", certID, "reason", "iot:Connect denied by policy")
		return false
	}

	cl.Properties.Username = []byte(certID)
	return true
}

// OnACLCheck evaluates the client's attached IoT policies to determine
// whether the requested topic operation is permitted.
func (h *AuthHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	if h.provider == nil {
		slog.Warn("mqtt acl denied", "client_id", cl.ID, "topic", topic, "action", "unknown", "reason", "no provider")
		return false
	}

	certID := string(cl.Properties.Username)
	if certID == "" {
		slog.Warn("mqtt acl denied", "client_id", cl.ID, "topic", topic, "reason", "no certificate id")
		return false
	}

	action := "iot:Subscribe"
	if write {
		action = "iot:Publish"
	}

	store := h.provider.GetStore()
	if store == nil {
		slog.Warn("mqtt acl denied", "cert_id", certID, "topic", topic, "action", action, "reason", "no store")
		return false
	}

	policies := getPoliciesForPrincipal(store, certID)
	if len(policies) == 0 {
		slog.Warn("mqtt acl denied", "cert_id", certID, "topic", topic, "action", action, "reason", "no policies attached")
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
		slog.Warn("mqtt acl denied", "cert_id", certID, "topic", topic, "action", action, "reason", "policy evaluation error", "error", err)
		return false
	}

	if !allowed {
		slog.Warn("mqtt acl denied", "cert_id", certID, "topic", topic, "action", action, "reason", "denied by policy")
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

type policyLookup interface {
	iotstore.PolicyOps
	iotstore.PolicyAttachmentOps
}

// getPoliciesForPrincipal retrieves all policies and parses them for
// evaluation. In a full implementation this would filter by attachment.
func getPoliciesForPrincipal(store policyLookup, certID string) []*policy.PolicyVersion {
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

// GetCertificateIDFromX509 computes the SHA-256 fingerprint of an X.509
// certificate in DER form.
func GetCertificateIDFromX509(cert *x509.Certificate) string {
	fingerprint := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(fingerprint[:])
}
