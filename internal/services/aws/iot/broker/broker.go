package broker

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"

	"vorpalstacks/internal/services/aws/iot/ca"
)

type MessageHandler func(clientID string, topic string, payload []byte)

type Broker struct {
	server         *mqtt.Server
	handler        MessageHandler
	ca             *ca.CertificateAuthority
	authProviderMu sync.RWMutex
	authProvider   AuthProvider
	port           int
	running        atomic.Bool
	ctx            context.Context
	cancel         context.CancelFunc
	closeOnce      sync.Once
	connTracker    *connectionTracker
}

func NewBroker(certificateAuthority *ca.CertificateAuthority, handler MessageHandler) *Broker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Broker{
		ca:          certificateAuthority,
		handler:     handler,
		ctx:         ctx,
		cancel:      cancel,
		connTracker: newConnectionTracker(),
	}
}

// SetAuthProvider configures the broker to use certificate and policy
// authentication. Can be called at any time; AuthHook reads the provider
// at connection time (fail-closed until set).
func (b *Broker) SetAuthProvider(provider AuthProvider) {
	b.authProviderMu.Lock()
	defer b.authProviderMu.Unlock()
	b.authProvider = provider
}

// getAuthProvider returns the current auth provider in a thread-safe
// manner. Returns nil before SetAuthProvider is called.
func (b *Broker) getAuthProvider() AuthProvider {
	b.authProviderMu.RLock()
	defer b.authProviderMu.RUnlock()
	return b.authProvider
}

// SetMessageHandler configures the handler invoked for every published
// MQTT message. Must be called before Start.
func (b *Broker) SetMessageHandler(handler MessageHandler) {
	b.handler = handler
}

func (b *Broker) Start(port int) error {
	if !b.running.CompareAndSwap(false, true) {
		return fmt.Errorf("broker already running")
	}

	b.port = port

	b.server = mqtt.New(&mqtt.Options{
		InlineClient: true,
	})

	// Always install AuthHook. It reads b.authProvider at connection time
	// via getAuthProvider(). Before SetAuthProvider is called, the provider
	// is nil and all connections are denied (fail-closed).
	_ = b.server.AddHook(NewAuthHook(b, b.connTracker), nil)
	_ = b.server.AddHook(&iotHook{
		handler: b.handler,
	}, nil)

	tlsConfig, err := b.buildTLSConfig()
	if err != nil {
		b.running.Store(false)
		return fmt.Errorf("failed to build mTLS config: %w", err)
	}

	tcp := listeners.NewTCP(listeners.Config{
		ID:        "iot-mqtt",
		Address:   fmt.Sprintf(":%d", port),
		TLSConfig: tlsConfig,
	})

	if err := b.server.AddListener(tcp); err != nil {
		b.running.Store(false)
		return fmt.Errorf("failed to add MQTT listener on port %d: %w", port, err)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("IoT MQTT broker panic recovered", "panic", r)
			}
		}()
		if err := b.server.Serve(); err != nil && b.running.Load() {
			slog.Error("IoT MQTT broker stopped", "error", err)
		}
	}()

	return nil
}

// buildTLSConfig constructs a mutual TLS configuration from the broker's
// CA. The server certificate is freshly issued on each Start; clients
// trust the CA root and need not pin the server certificate. If the CA
// is nil (e.g. unit tests), returns nil for plain TCP.
func (b *Broker) buildTLSConfig() (*tls.Config, error) {
	if b.ca == nil {
		return nil, nil
	}

	serverCert, err := b.ca.IssueServerCertificate("vorpalstacks-iot-broker")
	if err != nil {
		return nil, fmt.Errorf("failed to issue MQTT server certificate: %w", err)
	}

	clientCAs := x509.NewCertPool()
	clientCAs.AddCert(b.ca.RootCA())

	return &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAs,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

func (b *Broker) Stop() error {
	if !b.running.CompareAndSwap(true, false) {
		return nil
	}

	b.cancel()

	b.closeOnce.Do(func() {
		if b.server != nil {
			_ = b.server.Close()
		}
	})

	return nil
}

func (b *Broker) IsRunning() bool {
	return b.running.Load()
}

// Port returns the TCP port the broker is listening on, or 0 before Start.
func (b *Broker) Port() int {
	return b.port
}

func (b *Broker) Publish(topic string, payload []byte) error {
	if !b.running.Load() || b.server == nil {
		return fmt.Errorf("broker not running")
	}

	return b.server.Publish(topic, payload, false, 0)
}

func (b *Broker) Server() *mqtt.Server {
	return b.server
}

type iotHook struct {
	mqtt.HookBase
	handler MessageHandler
}

func (h *iotHook) ID() string { return "iot-hook" }

func (h *iotHook) Provides(b byte) bool {
	return b == mqtt.OnPublished
}

func (h *iotHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	if h.handler != nil {
		h.handler(cl.ID, pk.TopicName, pk.Payload)
	}
}

// IsCertConnected returns whether a certificate is currently connected to
// this broker, and the Unix timestamp of the connection establishment.
func (b *Broker) IsCertConnected(certID string) (bool, int64) {
	if b.connTracker == nil {
		return false, 0
	}
	return b.connTracker.isConnected(certID)
}

// connectionTracker records MQTT session establishment and disconnection
// keyed by certificate ID. It is safe for concurrent use.
type connectionTracker struct {
	clients sync.Map // certID → int64 (Unix timestamp)
}

func newConnectionTracker() *connectionTracker {
	return &connectionTracker{}
}

func (ct *connectionTracker) onConnect(certID string) {
	ct.clients.Store(certID, time.Now().Unix())
}

func (ct *connectionTracker) onDisconnect(certID string) {
	ct.clients.Delete(certID)
}

func (ct *connectionTracker) isConnected(certID string) (bool, int64) {
	v, ok := ct.clients.Load(certID)
	if !ok {
		return false, 0
	}
	ts, _ := v.(int64)
	return true, ts
}

// ConnectedCertIDs returns all certificate IDs currently connected.
func (ct *connectionTracker) ConnectedCertIDs() []string {
	var ids []string
	ct.clients.Range(func(key, _ any) bool {
		if id, ok := key.(string); ok {
			ids = append(ids, id)
		}
		return true
	})
	return ids
}

// ConnectedCertIDs returns all certificate IDs currently connected to
// this broker.
func (b *Broker) ConnectedCertIDs() []string {
	if b.connTracker == nil {
		return nil
	}
	return b.connTracker.ConnectedCertIDs()
}
