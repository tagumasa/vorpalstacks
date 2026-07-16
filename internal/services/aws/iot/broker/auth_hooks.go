package broker

import (
	"crypto/tls"
	"log/slog"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

	"vorpalstacks/internal/services/aws/iot/ca"
	"vorpalstacks/internal/services/aws/iot/policy"
	iotstore "vorpalstacks/internal/store/aws/iot"
	vcrypto "vorpalstacks/internal/utils/crypto"
)

// AuthProvider supplies the store and CA needed by the authentication hook.
type AuthProvider interface {
	GetStore() iotstore.IotStoreInterface
	GetCA() *ca.CertificateAuthority
	// CertificatePrincipal returns the principal ARN for a certificate
	// ID. IoT policy attachments are keyed by ARN, not by raw cert ID.
	CertificatePrincipal(certID string) string
}

// AuthHook validates client certificates and enforces IoT policy-based ACLs
// on publish and subscribe operations.
type AuthHook struct {
	mqtt.HookBase
	broker      *Broker
	connTracker *connectionTracker
}

// NewAuthHook creates an authentication hook that reads the AuthProvider
// from the Broker at connection time. Until SetAuthProvider is called,
// the provider is nil and all connections are denied (fail-closed).
func NewAuthHook(broker *Broker, connTracker *connectionTracker) *AuthHook {
	return &AuthHook{broker: broker, connTracker: connTracker}
}

// ID returns the hook identifier.
func (h *AuthHook) ID() string { return "iot-auth" }

// Provides declares which hook events this handler supports.
func (h *AuthHook) Provides(b byte) bool {
	return b == mqtt.OnConnectAuthenticate || b == mqtt.OnACLCheck ||
		b == mqtt.OnSessionEstablished || b == mqtt.OnDisconnect
}

// OnConnectAuthenticate validates the client certificate against the IoT
// certificate store and evaluates iot:Connect policy before accepting the
// connection.
func (h *AuthHook) OnConnectAuthenticate(cl *mqtt.Client, _ packets.Packet) bool {
	if h.broker == nil {
		return false
	}
	provider := h.broker.getAuthProvider()
	if provider == nil {
		return false
	}

	// extractCertificateID reads the client certificate fingerprint from
	// the TLS connection state. On an mTLS listener this yields the
	// SHA-256 fingerprint of the peer certificate; without TLS it
	// returns "" and authentication fails (fail-closed).
	certID := extractCertificateID(cl)
	if certID == "" {
		slog.Warn("mqtt auth denied", "client_id", cl.ID, "reason", "no certificate")
		return false
	}

	store := provider.GetStore()
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

	principal := provider.CertificatePrincipal(certID)
	policies := getPoliciesForPrincipal(store, principal)
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

	return true
}

// OnACLCheck evaluates the client's attached IoT policies to determine
// whether the requested topic operation is permitted.
func (h *AuthHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	if h.broker == nil {
		return false
	}
	provider := h.broker.getAuthProvider()
	if provider == nil {
		return false
	}

	certID := extractCertificateID(cl)
	if certID == "" {
		slog.Warn("mqtt acl denied", "client_id", cl.ID, "topic", topic, "reason", "no certificate")
		return false
	}

	action := "iot:Subscribe"
	if write {
		action = "iot:Publish"
	}

	store := provider.GetStore()
	if store == nil {
		slog.Warn("mqtt acl denied", "cert_id", certID, "topic", topic, "action", action, "reason", "no store")
		return false
	}

	principal := provider.CertificatePrincipal(certID)
	policies := getPoliciesForPrincipal(store, principal)
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

	return vcrypto.FingerprintX509(state.PeerCertificates[0])
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
			slog.Warn("mqtt auth: policy lookup failed, skipping (fail-closed)", "policy", name, "cert_id", certID, "error", err)
			continue
		}
		pv, err := policy.ParsePolicyVersion([]byte(p.PolicyDocument))
		if err != nil {
			slog.Warn("mqtt auth: policy parse failed, skipping (fail-closed)", "policy", name, "cert_id", certID, "error", err)
			continue
		}
		result = append(result, pv)
	}
	return result
}

// OnSessionEstablished records a successful MQTT connection in the
// connection tracker so that GetThingConnectivityData can report
// accurate connection state.
func (h *AuthHook) OnSessionEstablished(cl *mqtt.Client, _ packets.Packet) {
	certID := extractCertificateID(cl)
	if certID != "" && h.connTracker != nil {
		h.connTracker.onConnect(certID)
	}
}

// OnDisconnect removes the client from the connection tracker.
func (h *AuthHook) OnDisconnect(cl *mqtt.Client, _ error, _ bool) {
	certID := extractCertificateID(cl)
	if certID != "" && h.connTracker != nil {
		h.connTracker.onDisconnect(certID)
	}
}
